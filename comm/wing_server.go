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
	"google.golang.org/grpc/credentials"
	// "google.golang.org/grpc/credentials/insecure"
	"net"
	"os"
	"os/signal"
	"syscall"
)

var (
	upperDir   = GetUpperFileDir()
	tlsKeyFile = upperDir + "/apis/%s/tls-keys/%s.key"
	tlsPemFile = upperDir + "/apis/%s/tls-keys/%s.pem"
)

// Start and excute http server base on beego, by default, it just
// support restful interface not socket.io connection, but you can
// set allowCredentials as true to on socket.io conect function.
//
// `USAGE` :
//
//	// use for restful interface server
//	func main() {}
//		// comm.HttpServer(false) or
//		comm.HttpServer()
//
//	// use for both restful and socket.io server
//	func main() {
//		comm.HttpServer(true)
//	}
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

// Start and excute both restful and socket.io server
func Rest4SioServer() {
	HttpServer(true)
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
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type", "Authoration", "Token"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
	}))
}

// ----------------

// Callback to register Grpc client handler to server
type RegisterGrpcHandler func(svr *grpc.Server)

// GrpcServer start and excute grpc server, you can register grpc
// client handler by regfunc callback as follow:
//
// `USAGE`
//
//	go comm.GrpcServer(func(svr *grpc.Server) {
//		proto.RegisterAccServer(svr, &(handler.Acc{}))
//	})
func GrpcServer(regfunc RegisterGrpcHandler, server string, options ...string) {
	portkey := "nacosport"
	if len(options) != 0 && options[0] != "" {
		portkey = options[0]
	}

	port := beego.AppConfig.String(portkey)
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		panic("Listen grpc server, err:" + err.Error())
	}

	// TLS auth
	pem := fmt.Sprintf(tlsPemFile, server, server)
	key := fmt.Sprintf(tlsKeyFile, server, server)
	cred, err := credentials.NewServerTLSFromFile(pem, key)
	if err != nil {
		panic("Generate new TLS, err:" + err.Error())
	}

	svr := grpc.NewServer(grpc.Creds(cred))
	regfunc(svr)

	logger.I("Grpc server runing on", port)
	if err := svr.Serve(lis); err != nil {
		logger.E("Start grpc server, err:", err)
		panic(err)
	}
}

func DialGrpcServer(addr string, port int, server string) *grpc.ClientConn {
	domain := fmt.Sprintf("%s:%d", addr, port)

	// TLS auth
	pem := fmt.Sprintf(tlsPemFile, server, server)
	cred, err := credentials.NewClientTLSFromFile(pem, server)
	if err != nil {
		logger.E("Generate new TLS, err:", err)
		return nil
	}

	conn, err := grpc.Dial(domain, grpc.WithTransportCredentials(cred) /*insecure.NewCredentials())*/)
	if err != nil {
		logger.E("dial grpc address", domain, " fialed", err)
		return nil
	}
	logger.I("Connect grpc server", domain, " successed")
	return conn
}
