// Copyright (c) 2019-2029 DY All Rights Reserved.
//
// Author : yangping
// Email  : youhei_yp@163.com
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2022/02/09   yangping       New version
// -------------------------------------------------------------------

package core

import (
	sio "github.com/googollee/go-socket.io"
	"github.com/wengoldx/wing/comm"
	"github.com/wengoldx/wing/invar"
	"github.com/wengoldx/wing/logger"
	"sync"
)

// ClientPool client pool
type ClientPool struct {
	lock    sync.Mutex         // Mutex sync lock
	clients map[string]*client // Client map, seach key is client id
	s2c     map[string]string  // Socket id to Client id, seach key is socket id, value is client id

	waitOnCreate bool           // open or close waiting function, default is disable
	waitings     map[string]int // Idel clients map, seach key is client id, value is weights
}

// clientPool singleton instance
var clientPool *ClientPool

func init() {
	clientPool = &ClientPool{
		clients:  make(map[string]*client),
		s2c:      make(map[string]string),
		waitings: make(map[string]int),
	}
}

// Return ClientPool singleton
func Clients() *ClientPool {
	return clientPool
}

// Return client, it maybe nil if unexist.
func (cp *ClientPool) Client(cid string) *client {
	return cp.clients[cid]
}

// Register client and bind socket.
func (cp *ClientPool) Register(cid string, sc sio.Socket, option ...interface{}) error {
	cp.lock.Lock()
	defer cp.lock.Unlock()

	opt := comm.Condition(len(option) > 0, option[0], nil)
	if err := cp.registerLocked(cid, sc, opt); err != nil {
		logger.E("Regisger client err:", err.Error())
		return err
	}

	if cp.waitOnCreate {
		cp.waitingLocked(cid)
	}
	return nil
}

// Deregister client and unbind socket.
func (cp *ClientPool) Deregister(sc sio.Socket) (string, interface{}) {
	cp.lock.Lock()
	defer cp.lock.Unlock()

	return cp.deregisterLocked(sc)
}

// Check the client if exist.
func (cp *ClientPool) IsExist(cid string) bool {
	_, ok := cp.clients[cid]
	return ok
}

// Return client id from socket id.
func (cp *ClientPool) ClientID(sid string) string {
	cp.lock.Lock()
	defer cp.lock.Unlock()

	return cp.s2c[sid]
}

// Flag client to waiting state on create.
func (cp *ClientPool) WaitOnCreate(wait bool) {
	cp.waitOnCreate = wait
}

// Increate 1 of waiting weight for client.
func (cp *ClientPool) Waiting(cid string) {
	cp.lock.Lock()
	defer cp.lock.Unlock()

	cp.waitingLocked(cid)
}

// Reduce 1 of waiting weight for client, it will remove client from waiting
// map when weight value countdown to zero, and return true.
func (cp *ClientPool) Countdown(cid string) bool {
	cp.lock.Lock()
	defer cp.lock.Unlock()

	return cp.countdownLocked(cid)
}

// Remove client out of waiting map whatever weight value over zero or not.
func (cp *ClientPool) LeaveWaiting(cid string) {
	cp.lock.Lock()
	defer cp.lock.Unlock()

	cp.leaveWaitingLocked(cid)
}

// Return waiting clients ids
func (cp *ClientPool) IdleClients() []string {
	var idles []string
	for k, _ := range cp.waitings {
		idles = append(idles, k)
	}
	return idles
}

// -------- quick handle functions for indicate client

// Return client optinal data, it maybe nil.
func (cp *ClientPool) Option(cid string) interface{} {
	if c, ok := cp.clients[cid]; ok {
		return c.option
	}
	return nil
}

// Set the client optinal data, maybe return error if not exist client.
func (cp *ClientPool) SetOption(cid string, opt interface{}) error {
	if c, ok := cp.clients[cid]; ok {
		c.option = opt
		return nil
	}
	return invar.ErrNotFound
}

// Send signaling with message to indicate client.
func (cp *ClientPool) Signaling(cid string, evt invar.Event, data string) error {
	if c, ok := cp.clients[cid]; ok {
		return c.Send(evt, data)
	}
	return invar.ErrTagOffline
}

// --------

// Register the client without acquiring the lock.
func (cp *ClientPool) registerLocked(cid string, sc sio.Socket, opt interface{}) error {
	var newOne *client
	sid := sc.Id()

	if oldOne, ok := cp.clients[cid]; ok {
		oldOneID := oldOne.socket.Id()
		if oldOneID == sid {
			logger.W("Client", cid, "already bind socket", sid)
			return nil
		}

		logger.D("Drop socket", oldOneID, "of client", cid)
		delete(cp.s2c, oldOneID)
		oldOne.deregister() // reset and  disconnet the old socket
		newOne = oldOne
	} else {
		newOne = newClient(cid)
	}

	// bind client with socket
	if err := newOne.register(sc, opt); err != nil {
		return err
	}

	logger.D("Add client", cid, "bind socket", sid)
	cp.clients[cid] = newOne
	cp.s2c[sid] = cid // same as uuid
	return nil
}

// Deregister the client without acquiring the lock.
func (cp *ClientPool) deregisterLocked(sc sio.Socket) (string, interface{}) {
	sid := sc.Id()
	if cid := cp.s2c[sid]; cid != "" {
		delete(cp.s2c, sid)

		cp.leaveWaitingLocked(cid)
		if c := cp.clients[cid]; c != nil {
			delete(cp.clients, cid)
			c.deregister()
			return cid, c.option
		}
	}

	logger.I("Disconnect unkown client socket", sid)
	sc.Disconnect()
	return "", nil
}

// Increate waiting weight for client without acquiring the lock.
func (cp *ClientPool) waitingLocked(cid string) {
	weight := 1
	if w, ok := cp.waitings[cid]; ok {
		weight = w + 1
	}
	logger.D("Increate client", cid, "waiting weight:", weight)
	cp.waitings[cid] = weight
}

// Reduce waiting weight for client without acquiring the lock.
func (cp *ClientPool) countdownLocked(cid string) bool {
	if weight, ok := cp.waitings[cid]; ok {
		if weight > 1 {
			cp.waitings[cid] = weight - 1
			logger.D("Countdown client", cid, "waiting weight:", (weight - 1))
		} else if weight == 1 {
			logger.D("Countdown client", cid, "leave waiting")
			delete(cp.waitings, cid)
			return true
		}
	}
	return false
}

// Move client out of waiting state without acquiring the lock.
func (cp *ClientPool) leaveWaitingLocked(cid string) {
	if _, ok := cp.waitings[cid]; ok {
		logger.D("Client", cid, "leave waiting")
		delete(cp.waitings, cid)
	}
}
