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
	DID_ACC_CONFIGS  = "dunyu.acc.configs"  // Group by BASIC, data id of accservice cofnigs
	DID_API_ROUTERS  = "dunyu.api.routers"  // Group by BASIC, data id of swagger restful routers
	DID_DTALK_NTFERS = "dunyu.dtalk.ntfers" // Group by BASIC, data id of dingtalk notifiers
	DID_ES_AGENTS    = "dunyu.es.agents"    // Group by IFSC,  data id of elastic search agents
	DID_MIO_PATHS    = "dunyu.mio.paths"    // Group by BASIC, data id of minio source paths
	DID_MQTT_AGENTS  = "dunyu.mqtt.agents"  // Group by IFSC,  data id of mqtt agents
	DID_OTA_BUILDS   = "dunyu.ota.builds"   // Group by BASIC, data id of all projects OTA informations
	DID_WX_AGENTS    = "dunyu.wx.agents"    // Group by BASIC, data id of wechat agents

	// For GRPC certs content data ids
	DID_CERTK_ACC  = "dunyu.cert.acc.key"  // Group by BASIC, data id of cert key of accservice for grpc
	DID_CERTP_ACC  = "dunyu.cert.acc.pem"  // Group by BASIC, data id of cert pem of accservice for grpc
	DID_CERTK_CHAT = "dunyu.cert.chat.key" // Group by BASIC, data id of cert key of wgchat for grpc
	DID_CERTP_CHAT = "dunyu.cert.chat.pem" // Group by BASIC, data id of cert pem of wgchat for grpc
	DID_CERTK_MEA  = "dunyu.cert.mea.key"  // Group by BASIC, data id of cert key of measure for grpc
	DID_CERTP_MEA  = "dunyu.cert.mea.pem"  // Group by BASIC, data id of cert pem of measure for grpc
	DID_CERTK_PAY  = "dunyu.cert.pay.key"  // Group by BASIC, data id of cert key of wgpay for grpc
	DID_CERTP_PAY  = "dunyu.cert.pay.pem"  // Group by BASIC, data id of cert pem of wgpay for grpc
	DID_CERTK_WSS  = "dunyu.cert.wss.key"  // Group by BASIC, data id of cert key of webss for grpc
	DID_CERTP_WSS  = "dunyu.cert.wss.pem"  // Group by BASIC, data id of cert pem of webss for grpc
)

/* -------------------------- */
/* Export Configs defines     */
/* -------------------------- */

// Nacos config for data id DID_ACC_CONFIGS
type AccConfs struct {

	// Email sender service
	Email struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		User     string `json:"user"`
		Pwd      string `json:"pwd"`
		Identity string `json:"identity"`
	} `json:"email"`

	// SMS sender service
	Sms struct {
		Secret    string `json:"secret"`
		KeyID     string `json:"keyid"`
		URLFormat string `json:"urlformat"`
	} `json:"sms"`

	// Account secure settings
	Secures struct {
		SecureSalt   string `json:"secureSalt"`   // Secure salt key to decode account login token
		ApiTaxCode   string `json:"apiTaxCode"`   // Auth code to access API of check company tax code
		ApiIDViaCode string `json:"apiIDViaCode"` // Auth code to access API of identification check
		PageLimits   int    `json:"pageLimits"`   // One times to get list item counts on a page
	} `json:"secure"`

	// Administrators to allow login SysMgr
	Admins []string `json:"admin"`
}

// Nacos config for OTA upgrade by using DID_OTA_BUILDS data id
type OTAInfo struct {
	BuildVersion string `json:"BuildVersion" description:"Build version string"`
	BuildNumber  int    `json:"BuildNumber"  description:"Build number, pase form BuildVersion string as version = major*10000 + middle*100 + minor"`
	DownloadUrl  string `json:"DownloadUrl"  description:"Bin file download url"`
	UpdateDate   string `json:"UpdateDate"   description:"Bin file update date"`
	HashSums     string `json:"HashSums"     description:"Bin file hash sums"`
}

// Nacos config for DingTalk notify sender
type DTalkSender struct {
	WebHook   string   `json:"webhook"`   // DingTalk group chat session webhook
	Secure    string   `json:"secure"`    // DingTalk group chat senssion secure key
	Receivers []string `json:"receivers"` // The target @ users
}
