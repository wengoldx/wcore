// Copyright (c) 2018-2028 Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2019/05/22   yangping       New version
// 00002       2019/06/30   zhaixing       Add function from godfs
// -------------------------------------------------------------------

package comm

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/plugins/cors"
	"github.com/wengoldx/wing/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"os"
	"os/signal"
	"syscall"
)

// HttpServer start and excute http server base on beego
func HttpServer(allowCredentials ...bool) {
	ignoreSysSignalPIPE()
	if len(allowCredentials) > 0 {
		accessAllowOriginBy(beego.BeforeRouter, "*", allowCredentials[0])
		accessAllowOriginBy(beego.BeforeStatic, "*", allowCredentials[0])
	} else {
		accessAllowOriginBy(beego.BeforeRouter, "*", false)
		accessAllowOriginBy(beego.BeforeStatic, "*", false)
	}

	// just output log to file on prod mode
	if beego.BConfig.RunMode != "dev" &&
		logger.GetLevel() != logger.LevelDebug {
		beego.BeeLogger.DelLogger(logs.AdapterConsole)
	}
	beego.Run()
}

// Allow cross domain access for localhost,
// the port number must config in /conf/app.conf file like :
//
// ---
//
//	; Server port of HTTP
//	httpport=3200
func AccessAllowOriginByLocal(category int, allowCredentials bool) {
	if beego.BConfig.Listen.HTTPPort > 0 {
		localhosturl := fmt.Sprintf("http://127.0.0.1:%v/", beego.BConfig.Listen.HTTPPort)
		accessAllowOriginBy(category, localhosturl, allowCredentials)
	}
}

// Ignore system PIPE signal
func ignoreSysSignalPIPE() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGPIPE)
	go func() {
		for {
			select {
			case sig := <-sc:
				if sig == syscall.SIGPIPE {
					logger.E("!! IGNORE BROKEN PIPE SIGNAL !!")
				}
			}
		}
	}()
}

// Allow cross domain access for the given origins
func accessAllowOriginBy(category int, origins string, allowCredentials bool) {
	beego.InsertFilter("*", category, cors.Allow(&cors.Options{
		AllowAllOrigins:  !allowCredentials,
		AllowCredentials: allowCredentials,
		AllowOrigins:     []string{origins}, // use to set allow Origins
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
	}))
}

// ----------------

// Callback to register Grpc client handler to server
type RegisterGrpcHandler func(svr *grpc.Server)

// GrpcServer start and excute grpc server, you can register grpc
// client handler by regfunc callback as follow:
//
// `CODE`
//
//	go comm.GrpcServer(func(svr *grpc.Server) {
//		proto.RegisterAccServer(svr, &(handler.Acc{}))
//	})
func GrpcServer(regfunc RegisterGrpcHandler, options ...string) {
	portkey := "nacosport"
	if len(options) != 0 && options[0] != "" {
		portkey = options[0]
	}

	port := beego.AppConfig.String(portkey)
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		panic("Listen grpc server, err:" + err.Error())
	}

	svr := grpc.NewServer()
	regfunc(svr)

	logger.I("Grpc server runing on", port)
	if err := svr.Serve(lis); err != nil {
		logger.E("Start grpc server, err:", err)
		panic(err)
	}
}

func DialGrpcServer(addr string, port int) *grpc.ClientConn {
	domain := fmt.Sprintf("%s:%d", addr, port)
	conn, err := grpc.Dial(domain, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.E("dial grpc address", domain, " fialed", err)
		return nil
	}
	logger.I("Connect grpc server", domain, " successed")
	return conn
}
