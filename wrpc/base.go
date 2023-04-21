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

package wrpc

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/xml"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/wengoldx/wing/logger"
	"github.com/wengoldx/wing/nacos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"net"
	"strings"
)

const (
	SvrAcc  = "accservice" // server name of AccService backend
	SvrMea  = "measure"    // server name of Measure    backend
	SvrWss  = "webss"      // server name of WebSS      backend
	SvrChat = "wgchat"     // server name of WgChat     backend
	SvrPay  = "wgpay"      // server name of WgPay      backend
)

// GrpcHandlerFunc grpc server handler for register
type GrpcHandlerFunc func(svr *grpc.Server)

// ConnHandlerFunc grpc client handler for connect
type ConnHandlerFunc func(conn *grpc.ClientConn) interface{}

type GrpcStub struct {
	Certs   map[string]*nacos.GrpcCert // Grpc handler certs
	Clients map[string]interface{}     // Grpc client handlers

	// Current grpc server if registried
	isRegistried bool

	// Global handler function to return grpc server handler
	SvrHandlerFunc GrpcHandlerFunc
}

// Singleton grpc stub instance
var grpcStub *GrpcStub

// Return Grpc global singleton
func Singleton() *GrpcStub {
	if grpcStub == nil {
		grpcStub = &GrpcStub{
			isRegistried: false,
			Certs:        make(map[string]*nacos.GrpcCert),
			Clients:      make(map[string]interface{}),
		}
	}
	return grpcStub
}

// RegistServer start and excute grpc server, you numst setup global grpc
// register handler first as follow.
//
// `USAGE`
//
//	// set grpc server register handler
//	stub := wrpc.Singleton()
//	stub.SvrHandlerFunc = func(svr *grpc.Server) {
//		proto.RegisterAccServer(svr, &(handler.Acc{}))
//	}
//
//	// parse grps certs before register
//	stub.ParseCerts(data)
//
//	// register local server as grpc server
//	go stub.RegistServer()
func (stub *GrpcStub) RegistServer() {
	if stub.SvrHandlerFunc == nil {
		logger.E("Not setup global grpc handler!")
		return
	} else if stub.isRegistried {
		return // drop the duplicate registry
	}

	secure, ok := stub.Certs[beego.BConfig.AppName]
	if !ok || secure.Key == "" || secure.Pem == "" {
		logger.E("Not found grpc cert, abort register", beego.BConfig.AppName)
		return
	}

	// load grpc grpc server local port to listen
	port := beego.AppConfig.String("nacosport")
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		logger.E("Listen grpc server, err:", err)
		return
	}

	// generate TLS cert from pem datas
	cert, err := tls.X509KeyPair([]byte(secure.Pem), []byte(secure.Key))
	if err != nil {
		logger.E("Gen grpc cert, err:", err)
		return
	}

	// generate grpc server handler with TLS secure
	cred := credentials.NewServerTLSFromCert(&cert)
	svr := grpc.NewServer(grpc.Creds(cred))
	stub.SvrHandlerFunc(svr)
	logger.I("Grpc server runing on", port)

	stub.isRegistried = true
	defer func(stub *GrpcStub) { stub.isRegistried = false }(stub)
	if err := svr.Serve(lis); err != nil {
		logger.E("Start grpc server, err:", err)
	}
}

// Generate grpc client handler
func (stub *GrpcStub) GenClient(svrkey, addr string, port int, cb ConnHandlerFunc) {
	secure, ok := stub.Certs[svrkey]
	if !ok || secure.Key == "" || secure.Pem == "" {
		return
	}

	// generate TLS cert from pem datas
	cp := x509.NewCertPool()
	if !cp.AppendCertsFromPEM([]byte(secure.Pem)) {
		logger.E("Failed generate grpc cert!")
		return
	}

	// generate grpc client handler with TLS secure
	domain := fmt.Sprintf("%s:%d", addr, port)
	cred := credentials.NewClientTLSFromCert(cp, svrkey)
	conn, err := grpc.Dial(domain, grpc.WithTransportCredentials(cred))
	if err != nil {
		logger.E("dial grpc address", domain, " fialed", err)
		return
	}

	logger.I("Connectd grpc server", domain)
	stub.Clients[svrkey] = cb(conn)
}

// Parse all grpc certs from nacos config data, and cache to certs map
func (stub *GrpcStub) ParseCerts(data string) {
	certs := nacos.GrpcCerts{}
	if err := xml.Unmarshal([]byte(data), &certs); err != nil {
		logger.E("Parse grpc certs, err:", err)
		return
	}

	svrs := []string{}
	for _, cert := range certs.Certs {
		svrs = append(svrs, cert.Svr)
		stub.Certs[cert.Svr] = &cert
	}
	logger.D("Update grpc certs for:", strings.Join(svrs, ","))
}
