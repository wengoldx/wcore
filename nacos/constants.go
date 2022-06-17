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
	NS_META = "dunyu-meta-configs" // META namespace id
	NS_PROD = "dunyu-server-prod"  // PROD namespace id
	NS_DEV  = "dunyu-server-dev"   // DEV  namespace id
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
	DID_ADMINS     = "dunyu.acc.admins" // Data id of sysmgr admins,      it must group by BASIC
	DID_WX_AGENTS  = "dunyu.wx.agents"  // Data id of wechat agents,      it must group by BASIC
	DID_NTF_SENDER = "dunyu.ntf.sender" // Data id of dingtalk nofitier,  it must group by BASIC
	DID_MIO_PATHS  = "dunyu.mio.paths"  // Data id of minio source paths, it must group by BASIC
)
