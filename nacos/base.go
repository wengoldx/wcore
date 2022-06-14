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
	"github.com/wengoldx/wing/comm"
	"github.com/wengoldx/wing/logger"
)

const (
	LogLevel = "warn"          // default level to print nacos logs on warn, it not same nacos-sdk-go:info
	logDir   = "./nacos/logs"  // default nacos logs dir
	cacheDir = "./nacos/cache" // default nacos caches dir

	NS_META = "dunyu-meta-configs" // META namespace id
	NS_PROD = "dunyu-server-prod"  // PROD namespace id
	NS_DEV  = "dunyu-server-dev"   // DEV  namespace id

	GP_BASIC = "group.basic" // BASIC group name
	GP_IFSC  = "group.ifsc"  // IFSC  group name
	GP_DTE   = "group.dte"   // DTE   group name
	GP_CWS   = "group.cws"   // CWS   group name
)

// Generate nacos client config, contain nacos remote server and
// current business servers configs, this client keep alive with
// 5s pingpong heartbeat and output logs on warn leven.
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
		Username:            "accessor", // secure account
		Password:            "accessor", // secure passowrd
	}

	return vo.NacosClientParam{
		ClientConfig: cc, ServerConfigs: sc,
	}
}

// Load register server's nacos informations
//	@return - string nacos remote server ip
//			- string group name of local server
func LoadNacosSvrConfigs() (string, string) {
	svr := beego.AppConfig.String("nacossvr")
	if svr == "" {
		panic("Not found nacos server host!")
	}

	gp := beego.AppConfig.String("nacosgp")
	if !(gp == GP_BASIC || gp == GP_IFSC || gp == GP_DTE || gp == GP_CWS) {
		panic("Invalid register cluster group!")
	}
	return svr, gp
}

// -------- Auto Register Define --------

// Server register informations
type ServerItem struct {
	Name     string         // Server name, same as beego app name
	Group    string         // Server group, range in [GP_BASIC, GP_IFSC, GP_DTE, GP_CWS]
	Callback ServerCallback // Server register datas changed callback
}

// Callback to listen server address and port changes
type ServerCallback func(svr string, addr string, port int)

// Register current server to nacos, you must set configs in app.conf
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
//	; Inner net address for dev servers access
//	nacosaddr = "10.239.20.99"
//
//	; Inner net port for grpc access
//	nacosport = 3000
//
//	[prod]
//	; Inner net address for prod servers access
//	nacosaddr = "10.239.40.64"
//
//	; Inner net port for grpc access
//	nacosport = 3000
func RegisterServer() *ServerStub {
	svr, group := LoadNacosSvrConfigs()
	return RegisterServer2(svr, group)
}

// Register current server to nacos by given nacos server host and group
func RegisterServer2(svr, group string) *ServerStub {
	// Local server listing ip
	addr := beego.AppConfig.String("nacosaddr")
	if addr == "" {
		panic("Not found local server ip to register!")
	}

	// Server access port for grpc, it maybe same as httpport config
	// when the local server not support grpc but for http
	port, err := beego.AppConfig.Int("nacosport")
	if err != nil || port < 3000 /* remain 0 ~ 3000 */ {
		panic("Not found port number or less 3000!")
	}

	// Namespace id of local server
	ns := comm.Condition(beego.BConfig.RunMode == "prod",
		NS_PROD, NS_DEV).(string)

	// Generate nacos server stub and setup it
	stub := NewServerStub(ns, svr)
	if err := stub.Setup(); err != nil {
		panic(err)
	}

	// Fixed app name as nacos server name to register,
	// and pick server port from config 'nacosport' not form 'httpport' value,
	// becase it maybe support either grpc or http hanlder to accesse.
	//
	// And here not use cluster name, please keep it empty!
	app, port := beego.BConfig.AppName, beego.BConfig.Listen.HTTPPort
	if err := stub.Register(app, addr, uint64(port), group); err != nil {
		panic(err)
	}

	logger.I("Registered server", app+"@"+addr)
	return stub
}

// Listing services address and port changes, it will call the callback
// immediately to return target service host when them allready registerd
// to service central of nacos.
//	@params servers []*ServerItem target server registry informations
func (ss *ServerStub) ListenServers(servers []*ServerItem) {
	for _, s := range servers {
		if err := ss.Subscribe(s.Name, s.OnChanged, s.Group); err != nil {
			panic("Subscribe server " + s.Name + " err:" + err.Error())
		}
	}
}

// Subscribe callback called when target service address and port changed
func (si *ServerItem) OnChanged(services []model.SubscribeService, err error) {
	if err != nil {
		logger.E("Received server", si.Name, "change, err:", err)
		return
	}

	if len(services) > 0 {
		addr, port := services[0].Ip, services[0].Port
		logger.I("Update server", si.Name, "to {", addr, "-", port, "}")

		si.Callback(si.Name, addr, int(port))
	}
}

// -------- Config Service Define --------

// Meta config informations
type MetaConfig struct {
	Group     string                        // Server group, range in [GP_BASIC, GP_IFSC, GP_DTE, GP_CWS]
	Stub      *ConfigStub                   // Nacos config client instance
	Callbacks map[string]MetaConfigCallback // Meta config changed callback maps, key is dataid
}

// Callback to listen server address and port changes
type MetaConfigCallback func(dataId, data string)

// Generate a meta config client to get or listen configs changes
//	@return - *MetaConfig nacos config client instance on NS_META namespace
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
func GenMetaConfig() *MetaConfig {
	svr, group := LoadNacosSvrConfigs()
	return GenMetaConfig2(svr, group)
}

// Generate a meta config client by given nacos server host and group
func GenMetaConfig2(svr, group string) *MetaConfig {
	stub := NewConfigStub(NS_META /* Fixed ns */, svr)
	if err := stub.Setup(); err != nil {
		panic("Gen config stub, err:" + err.Error())
	}

	cbs := make(map[string]MetaConfigCallback)
	return &MetaConfig{
		Group: group, Stub: stub, Callbacks: cbs,
	}
}

// Get and listing the config of indicated dataId
func (mc *MetaConfig) ListenConfig(dataId string, cb MetaConfigCallback) {
	mc.Callbacks[dataId] = cb // cache callback

	// get config first before listing
	data, err := mc.Stub.GetString(dataId, mc.Group)
	if err != nil {
		panic("Get config " + dataId + "err: " + err.Error())
	}
	cb(dataId, data)

	// listing config changes
	logger.I("Listen config { dataId:", dataId, "group:", mc.Group, "}")
	mc.Stub.Listen(dataId, mc.Group, mc.OnChanged)
}

// Listing callback called when target configs changed
func (mc *MetaConfig) OnChanged(namespace, group, dataId, data string) {
	if namespace != NS_META || group != mc.Group {
		logger.E("Invalid meta config ns:", namespace, "or group:", group)
		return
	}

	if callback, ok := mc.Callbacks[dataId]; ok {
		logger.I("Update config dataId", dataId, "to:", data)
		callback(dataId, data)
	}
}
