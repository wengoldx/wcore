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

/* -------------------------- */
/* Internal constants defines */
/* -------------------------- */

const (
	nacosSysSecure = "accessor"      // CICD system secure authentications
	nacosLogLevel  = "warn"          // default level to print nacos logs on warn, it not same nacos-sdk-go:info
	nacosDirLogs   = "./nacos/logs"  // default nacos logs dir
	nacosDirCache  = "./nacos/cache" // default nacos caches dir

	configKeySvr   = "nacossvr"  // Nacos remote server IP address
	configKeyGroup = "nacosgp"   // Local server group
	configKeyAddr  = "nacosaddr" // Local server access IP address
	configKeyPort  = "nacosport" // Local server access port for grpc connect
	configKeyHPort = "httpport"  // Local server http port
)

/* -------------------------- */
/* Export constants defines   */
/* -------------------------- */

// Nacos namespace string for wing/nacos
const (
	NS_PROD = "dunyu-server-prod" // PROD namespace id
	NS_DEV  = "dunyu-server-dev"  // DEV  namespace id
)

// Nacos group string for wing/nacos
const (
	GP_BASIC = "group.basic" // BASIC group name
	GP_IFSC  = "group.ifsc"  // IFSC  group name
	GP_DTE   = "group.dte"   // DTE   group name
	GP_CWS   = "group.cws"   // CWS   group name
)

// Nacos data id for wing/nacos
const (
	DID_ACC_ADMINS   = "dunyu.acc.admins"   // Group by BASIC, data id of sysmgr admins
	DID_ACC_SALT     = "dunyu.acc.secure"   // Group by BASIC, data id of account secure salt
	DID_ACC_EMAIL    = "dunyu.acc.email"    // Group by BASIC, data id of email sender settings
	DID_ACC_SMS      = "dunyu.acc.sms"      // Group by BASIC, data id of sms sender settings
	DID_WX_AGENTS    = "dunyu.wx.agents"    // Group by BASIC, data id of wechat agents
	DID_MIO_PATHS    = "dunyu.mio.paths"    // Group by BASIC, data id of minio source paths
	DID_NTF_SENDER   = "dunyu.ntf.sender"   // Group by BASIC, data id of dingtalk gym watchdog nofitier
	DID_NTF_PROPOSER = "dunyu.ntf.proposer" // Group by BASIC, data id of dingtalk suggestion nofitier
	DID_NTF_WGPAY    = "dunyu.ntf.wgpay"    // Group by BASIC, data id of dingtalk wgpay nofitier
	DID_NTF_ORDER    = "dunyu.ntf.order"    // Group by BASIC, data id of dingtalk trade order nofitier
	DID_ES_AGENTS    = "dunyu.es.agents"    // Group by IFSC, data id of elastic search agents
	DID_MQTT_AGENTS  = "dunyu.mqtt.agents"  // Group by IFSC, data id of mqtt agents
	DID_OTA_STOR     = "dunyu.ota.store"    // Group by BASIC, data id of store OTA
	DID_OTA_MAKER    = "dunyu.ota.maker"    // Group by BASIC, data id of maker OTA
	DID_OTA_SHOW     = "dunyu.ota.show"     // Group by BASIC, data id of show OTA
	DID_OTA_QKS      = "dunyu.ota.qks"      // Group by BASIC, data id of QKS OTA
)
