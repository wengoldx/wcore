// Copyright (c) 2019-2029 DY All Rights Reserved.
//
// Author : yangping
// Email  : youhei_yp@163.com
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2022/02/09   yangping       New version
// -------------------------------------------------------------------

package wsio

import (
	"github.com/astaxie/beego"
	sio "github.com/googollee/go-socket.io"
	"github.com/wengoldx/wing/invar"
	"github.com/wengoldx/wing/logger"
	"github.com/wengoldx/wing/wsio/clients"
	"net/http"
	"sync"
	"time"
	"unsafe"
)

// client datas for temp cache
type clientOpt struct {
	UID string
	Opt string
}

// Socket.io connecte information, it will generate a socket server
// and register as http Handler to listen socket signalings.
//
// ----
//
// `USAGE` :
//
//	// routers.go : register socket events
//	import "github.com/wengoldx/wing/wsio"
//
//	init() {
//		// set socket io hander and signalings
//		wsio.SetHandlers(ctrl.Authenticate, nil, nil)
//		adaptor := &ctrl.DefSioAdaptor{}
//		if err := wsio.SetAdapter(adaptor); err != nil {
//			panic(err)
//		}
//	}
type wingSIO struct {
	// Mutex sync lock, protect client connecting
	lock sync.Mutex

	// http request pointer to client, cache datas temporary
	// only for client authenticate-connect process.
	caches map[uintptr]*clientOpt

	// socket server
	server *sio.Server

	// socket golbal handler to execute clients authenticate action.
	authHandler AuthHandler

	// socket golbal handler to execute clients connect action.
	connHandler ConnectHandler

	// socket golbal handler to execute clients disconnect actions
	discHandler DisconnectHandler
}

// Socket connection server
var wsc *wingSIO

// Check client option if empty when connnection is established,
// if optinal data is empty the connect will not establish and disconnect.
var UsingOption = false

var (
	serverPingInterval = 30 * time.Second
	serverPingTimeout  = 60 * time.Second
	maxConnectCount    = 200000
)

func init() {
	setupWsioConfigs()
	wsc = &wingSIO{
		caches: make(map[uintptr]*clientOpt),
	}

	// set http handler for socke.io
	handler, err := wsc.createHandler()
	if err != nil {
		panic(err)
	}

	// set socket.io routers
	beego.Handler("/"+beego.BConfig.AppName+"/socket.io", handler)
	logger.I("Initialized socket.io routers...")
}

// Set handler to execute clients authenticate, connect and disconnect.
func SetHandlers(auth AuthHandler, conn ConnectHandler, disc DisconnectHandler) {
	wsc.authHandler, wsc.connHandler, wsc.discHandler = auth, conn, disc
	logger.D("Set wsio handlers")
}

// Set adapter to register socket signaling events.
func SetAdapter(adaptor SignalingAdaptor) error {
	if adaptor == nil {
		logger.W("Invalid socket event adaptor!")
		return nil
	}

	evts := adaptor.Signalings()
	if evts == nil || len(evts) == 0 {
		logger.W("No signaling event keys!")
		return nil
	}

	// register socket signaling events
	for _, evt := range evts {
		if evt != "" {
			callback := adaptor.Dispatch(evt)
			if callback != nil {
				if err := wsc.server.On(evt, callback); err != nil {
					return err
				}
				logger.D("Bind signaling event", evt)
			}
		}
	}
	return nil
}

// read wsio configs from file
func setupWsioConfigs() {
	if interval, err := beego.AppConfig.Int64("wsio::interval"); err != nil {
		logger.W("Read wsio::interval, err:", err)
	} else if interval > 0 {
		serverPingInterval = time.Duration(interval) * time.Second
	}

	if timeout, err := beego.AppConfig.Int64("wsio::timeout"); err != nil {
		logger.W("Read wsio::timeout, err:", err)
	} else if timeout > 0 {
		serverPingTimeout = time.Duration(timeout) * time.Second
	}

	if using, err := beego.AppConfig.Bool("wsio::optinal"); err != nil {
		logger.W("Load wsio::optinal, err:", err)
	} else if using {
		UsingOption = using
	}

	// logout the configs value
	logger.D("Configs interval:", serverPingInterval,
		"timeout:", serverPingTimeout, "optional:", UsingOption)
}

// createHandler create http handler for socket.io
func (cc *wingSIO) createHandler() (http.Handler, error) {
	server, err := sio.NewServer(nil)
	if err != nil {
		return nil, err
	}
	cc.server = server

	// set socket.io ping interval and timeout
	logger.I("Set socket ping-pong and timeout")
	server.SetPingInterval(serverPingInterval)
	server.SetPingTimeout(serverPingTimeout)

	// set max connection count
	server.SetMaxConnection(maxConnectCount)

	// set auth middleware for socket.io connection
	server.SetAllowRequest(func(req *http.Request) error {
		if err = cc.onAuthentication(req); err != nil {
			logger.E("Authenticate err:", err)
			return err
		}
		return nil
	})

	// set connection event
	server.On("connection", func(sc sio.Socket) {
		cc.onConnect(sc)
	})

	// set disconnection event
	server.On("disconnection", func(sc sio.Socket) {
		cc.onDisconnected(sc)
	})

	logger.I("Created socket.io handler")
	return server, nil
}

// onAuthentication event of authentication
func (cc *wingSIO) onAuthentication(req *http.Request) error {
	if err := req.ParseForm(); err != nil {
		return err
	}

	token := req.Form.Get("token")
	if token == "" {
		return invar.ErrInvalidClient
	}

	// auth client token by handler if set Authenticate function
	// handler, or just use token as uuid when not set.
	uuid, option := token, ""
	if cc.authHandler != nil {
		uid, opt, err := cc.authHandler(token)
		if err != nil || uid == "" {
			return invar.ErrAuthDenied
		} else if UsingOption && opt == "" {
			logger.E("Empty client", uid, "option data!")
			return invar.ErrAuthDenied
		}

		logger.D("Decode token, uuid:", uid, "opt:", opt)
		uuid, option = uid, opt
	}

	// bind http.Request -> uuid
	h := uintptr(unsafe.Pointer(req))
	logger.D("Bind request:", h, "with client:", uuid)
	cc.bindHTTP2UUIDLocked(h, uuid, option)
	return nil
}

// onConnect event of connect
func (cc *wingSIO) onConnect(sc sio.Socket) {
	// find client uuid and unbind -> http.Request
	h := uintptr(unsafe.Pointer(sc.Request()))
	co := cc.unbindUUIDFromHTTPLocked(h)
	if co == nil || co.UID == "" {
		logger.E("Invalid socket request bind!")
		sc.Disconnect()
		return
	}

	clientPool := clients.Clients()
	if err := clientPool.Register(co.UID, sc, co.Opt); err != nil {
		logger.E("Faild register socket client:", co.UID)
		sc.Disconnect()
		return
	}

	// handle connect callback for socket with uuid
	if cc.connHandler != nil {
		if err := cc.connHandler(co.UID, co.Opt); err != nil {
			logger.E("Client:", co.UID, "connect socket err:", err)
			sc.Disconnect()
		}
	}
	logger.I("Connected socket client:", co.UID)
}

// onDisconnected event of disconnect
func (cc *wingSIO) onDisconnected(sc sio.Socket) {
	uuid, opt := clients.Clients().Deregister(sc)
	if cc.discHandler != nil {
		cc.discHandler(uuid, opt)
	}
}

// bindHTTP2UUIDLocked bind http request pointer -> uuid on locked status
func (cc *wingSIO) bindHTTP2UUIDLocked(h uintptr, uuid, opt string) {
	cc.lock.Lock()
	defer cc.lock.Unlock()

	cc.caches[h] = &clientOpt{UID: uuid, Opt: opt}
}

// unbindUUIDFromHTTPLocked unbind uuid -> http request pointer on locked status
func (cc *wingSIO) unbindUUIDFromHTTPLocked(h uintptr) *clientOpt {
	cc.lock.Lock()
	defer cc.lock.Unlock()

	if data, ok := cc.caches[h]; ok {
		co := &clientOpt{UID: data.UID, Opt: data.Opt}
		delete(cc.caches, h)
		return co
	}
	return nil
}
