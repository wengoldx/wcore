// Copyright (c) 2018-2028 Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2022/05/11   yangping     New version
// -------------------------------------------------------------------

package nacos

import (
	"github.com/astaxie/beego"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/wengoldx/wing/logger"
)

const (
	LogLevel = "warn"          // default level to print nacos logs on warn, it not same nacos-sdk-go:info
	logDir   = "./nacos/logs"  // default nacos logout dir
	cacheDir = "./nacos/cache" // default nacos caches dir

	NS_META = "dunyu-meta-configs" // META namespace id
	NS_PROD = "dunyu-server-prod"  // PROD namespace id
	NS_DEV  = "dunyu-server-dev"   // DEV  namespace id
	NS_TEST = "dunyu-test-ns"      // TEST namespace id

	GP_BASIC = "group.basic" // BASIC group name
	GP_IFSC  = "group.ifsc"  // IFSC  group name
	GP_DTE   = "group.dte"   // DTE   group name
	GP_CWS   = "group.cws"   // CWS   group name
)

// Generate nacos client config, contain nacos remote server and
// current business servers configs, this client keep alive with
// 5s pingpong heartbeat and logout logs on info leven.
//
//	`NOTICE`
//	The remote server must access on http://{svr}:8848/nacos
func genClientParam(ns, svr string) vo.NacosClientParam {
	sc := []constant.ServerConfig{
		constant.ServerConfig{
			Scheme: "http", ContextPath: "/nacos", IpAddr: svr, Port: 8848,
		},
	}

	cc := &constant.ClientConfig{
		NamespaceId:         ns,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              logDir,
		CacheDir:            cacheDir,
		LogRollingConfig:    &constant.ClientLogRollingConfig{MaxSize: 10},
		LogLevel:            LogLevel,
	}

	return vo.NacosClientParam{
		ClientConfig: cc, ServerConfigs: sc,
	}
}

// Check namespace and group values, it will case panic when this
// tow feilds values invalid
func accessValidParams(ns, gp string) {
	validns := (ns == NS_PROD || ns == NS_DEV)
	validgp := (gp == GP_BASIC || gp == GP_IFSC || gp == GP_DTE || gp == GP_CWS)
	if !validns || !validgp {
		panic("Invalid namespace and group!")
	}
}

// Load register server's nacos informations
//	@return - string nacos remote server ip
//			- string local server listing ip
//			- string namespace id of local server
//			- string group name of local server
func loadConfigs() (string, string, string, string) {
	svr := beego.AppConfig.String("nacossvr")
	addr := beego.AppConfig.String("nacosaddr")
	if svr == "" || addr == "" {
		panic("Not found nacos configs!")
	}

	ns := beego.AppConfig.String("nacosns")
	gp := beego.AppConfig.String("nacosgp")
	accessValidParams(ns, gp)

	return svr, addr, ns, gp
}

// -------- Auto Register Define --------

// Server register informations
type ServerItem struct {
	Name     string           // Server name, same as beego app name
	Group    string           // Server group, range in [GP_BASIC, GP_IFSC, GP_DTE, GP_CWS]
	Callback RegisterCallback // Server register datas changed callback
}

// Callback to listen server address and port changed
type RegisterCallback func(svr string, addr string, port int)

// Register the given server, you must set configs in app.conf
//	@return - *ServerStub nacos server stub instance
//
//	`NOTICE` : nacos config as follows.
//
// ----
//
//	; Nacos remote server host
//	nacossvr = "10.239.40.24"
//
//	; Server nacos group name
//	nacosgp = "group.ifsc"
//
//	[dev]
//	; Nacos namespace id
//	nacosns = "dunyu-server-dev"
//
//	; Inner net address for dev servers access
//	nacosaddr = "10.239.20.99"
//
//	[prod]
//	; Nacos namespace id
//	nacosns = "dunyu-server-prod"
//
//	; Inner net address for prod servers access
//	nacosaddr = "10.239.40.64"
func RegisterServer() *ServerStub {
	svr, addr, ns, group := loadConfigs()

	// Generate nacos server stub and setup it
	stub := NewServerStub(ns, svr)
	if err := stub.Setup(); err != nil {
		panic(err)
	}

	// Fixed app name as nacos server name to register, and pick server port
	// from config 'httpport' value, but not pick server host ip from config
	// 'httpaddr' value, becase it empty and the nginx proxy server need it
	// keep empty, so get server host as input param by 'addr'.
	//
	// And here not use cluster name, please keep it empty!
	app := beego.BConfig.AppName
	port := beego.BConfig.Listen.HTTPPort
	if err := stub.Register(app, addr, uint64(port), group); err != nil {
		panic(err)
	}

	logger.I("Registered server:", app, "addr:", addr, "to nacos")
	return stub
}

// Get target server informations and listen status change
//	@params stub    *ServerStub   nacos server stub instance
//	@params targets []*ServerItem target server registry informations
func GetAndListen(stub *ServerStub, targets []*ServerItem) {
	if stub == nil || len(targets) == 0 {
		panic("Invalid server stub or empty targets servers!")
	}

	for _, tag := range targets {
		if svr, err := stub.GetServer(tag.Name, tag.Group); err != nil {
			panic("Get server, err:" + err.Error())
		} else if len(svr.Hosts) > 0 {
			addr, port := svr.Hosts[0].Ip, svr.Hosts[0].Port
			logger.I("Get server:", tag.Name, "on", addr, port)

			tag.Callback(tag.Name, addr, int(port))
		}

		// subscribe target server registry changed event
		if err := stub.Subscribe(tag.Name, tag.OnChanged, tag.Group); err != nil {
			panic("Subscribe " + tag.Name + " err:" + err.Error())
		}
	}
}

// Subscribe callback called when target server registry changed
func (s *ServerItem) OnChanged(services []model.SubscribeService, err error) {
	if err == nil && len(services) > 0 {
		addr, port := services[0].Ip, services[0].Port
		logger.I("Update server", s.Name, "changed to", addr, port)

		s.Callback(s.Name, addr, int(port))
		return
	}
}
