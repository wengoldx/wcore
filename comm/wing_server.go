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
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/plugins/cors"
	"github.com/wengoldx/wing/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"net"
	"os"
	"os/signal"
	"syscall"
)

// ===========================
// For setup HTTP Server
// ===========================

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

// ===========================
// For GRPC Client and Server
// ===========================

// GrpcHandlerFunc register grpc client handler to server
type GrpcHandlerFunc func(svr *grpc.Server)

// Global handler function to register caller server as grpc server
var GGrpcHandlerFunc GrpcHandlerFunc

// GrpcServer start and excute grpc server, you can setup global grpc
// register handler first as follow, it maybe throw panic when case error.
//
// `USAGE`
//
//	GGrpcHandlerFunc = func(svr *grpc.Server) {
//		proto.RegisterAccServer(svr, &(handler.Acc{}))
//	}
func GrpcServer(certkey, certpem string, options ...string) {
	if GGrpcHandlerFunc == nil {
		panic("Not setup global grpc handler!")
	}

	portkey := "nacosport"
	if len(options) != 0 && options[0] != "" {
		portkey = options[0]
	}

	port := beego.AppConfig.String(portkey)
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		panic("Listen grpc server, err:" + err.Error())
	}

	// generate TLS cert from pem datas
	cert, err := tls.X509KeyPair([]byte(certpem), []byte(certkey))
	if err != nil {
		panic("Gen grpc server cert, err:" + err.Error())
	}

	// generate grpc server handler with TLS secure
	cred := credentials.NewServerTLSFromCert(&cert)
	svr := grpc.NewServer(grpc.Creds(cred))
	GGrpcHandlerFunc(svr)

	logger.I("Grpc server runing on", port)
	if err := svr.Serve(lis); err != nil {
		logger.E("Start grpc server, err:", err)
		panic(err)
	}
}

// Generate grpc client handler
func GrpcClient(addr string, port int, server, certpem string) *grpc.ClientConn {
	// generate TLS cert from pem datas
	cp := x509.NewCertPool()
	if !cp.AppendCertsFromPEM([]byte(certpem)) {
		logger.E("Failed generate grpc cert!")
		return nil
	}

	// generate grpc client handler with TLS secure
	domain := fmt.Sprintf("%s:%d", addr, port)
	cred := credentials.NewClientTLSFromCert(cp, server)
	conn, err := grpc.Dial(domain, grpc.WithTransportCredentials(cred))
	if err != nil {
		logger.E("dial grpc address", domain, " fialed", err)
		return nil
	}
	logger.I("Connect grpc server", domain, " successed")
	return conn
}
