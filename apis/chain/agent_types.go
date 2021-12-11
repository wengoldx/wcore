// Copyright (c) 2018-2022 WING All Rights Reserved.
//
// Author : jidi
// Email  : j18041361158@163.com
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2021/11/17   jidi           New version
// -------------------------------------------------------------------

package chain

const (
	PLACE_ORDER int64 = iota + 1 // place an order, status 1
	CANCELLED                    // order was cancelled, status 2
	UNPAID                       // order not paid or paied fail, status 3
	PAID                         // goods to be delivered, status 4
)

// PaychainAgent agent of paychain
//
// `USAGE` :
//
// You can generate a PaychainAgent instance to quick access paychain service APIs.
//
// * `Step 1` - register a agent account on paychain service by API /v2/register
//
// * `Step 2` - check the agent if activity, you may use API /v2/state to change it
//
// * `Step 3` - config paychian server domain in project config files as follow:
//
// #### conf/app.conf
//
//	[dev]
//	; Domain url
//	domain=https://www.sampledomein.com:xxx4
//
//	[prod]
//	; Domain url
//	domain=https://www.sampledomein.com:xxx3
//
// * `Step 4` - generate a PaychainAgent and use them public methods
//
// #### generate instance
//
//	agentIns := &order.PaychainAgent {
//		Aid    : "agent-id",
//		Devmac : "xx:xx:xx:xx:xx:xx",
//		Pubkey : "xxxxxxxxxxxxxxxxx",
//		Domain : beego.AppConfig.String("domain"),
//	}
//
// #### public methods
//
//	// generate a new trade ticket
//	tno, err := agentIns.GenTicket(ps)
//
//	// get the latest trade ticket node
//	ticket, err := agentIns.TradeTicket(tno)
//
//	// get the latest dividing ticket node
//	ticket, err := agentIns.DiviTicket(tno)
//
//	// get the latest refund ticket node
//	ticket, err := agentIns.RefundTicket(tno)
//
//	// update indicated trade ticket
//	err := agentIns.UpdateTicket(tno, ps)
type PaychainAgent struct {
	Aid    string // agent id
	Devmac string // device mac bind with current agent
	Pubkey string // public key of current agent

	// Paychain service access url domain, it allways get from app.conf
	// by beego.AppConfig.String("domain") code.
	Domain string
}

// EncryptNode encrypt node data struct
type EncryptNode struct {
	SecureKey string `json:"securekey"`
	Timestamp int64  `json:"timestamp"`
	PayBody   string `json:"paybody"`
}

// TradeNode Trade ticket node
type TradeNode struct {
	Service    string `json:"service"               description:"service name"`
	CUUID      string `json:"cuuid"                 description:"payer uuid"`
	SUUID      string `json:"suuid"                 description:"payee uuid, merchant id"`
	SubMchID   string `json:"sub_mchid"             description:"payee sub merchant id"`
	Amount     int64  `json:"amount"                description:"amount price, unit one cent CNY"`
	RedundFee  int64  `json:"refundfee"             description:"total refund price, unit one cent CNY"`
	Desc       string `json:"desc"                  description:"this order description"`
	NotifyURL  string `json:"notifyurl"             description:"the notify url to tell service that payment success"`
	PayWay     string `json:"payway"                description:"payment way, such as 'wehcat' and 'alipay'"`
	IsFrozen   bool   `json:"isfrozen"              description:"whether frozen amount when payment finished。it must be true, when you want to share money"`
	Status     int64  `json:"status"                description:"payment status, such as 'cancle', 'unpaid', 'paid'"`
	TimeExpire string `json:"time_expire,omitempty" description:"expire time"`
}

// Dividing ticket node
type DiviNode struct {
	Service       string `json:"service"        description:"service name"`
	SubMchID      string `json:"sub_mchid"      description:"payee sub merchant id"`
	TransactionID string `json:"transaction_id" description:"wechat transaction order id"`
	Commission    int64  `json:"commission"     description:"share out money, unit one cent CNY"`
	Desc          string `json:"desc"           description:"this share description"`
	IsFinsh       bool   `json:"isfinsh"        description:"finish trade order, and unfrozen order"`
}

// RefundNode refund ticket node
type RefundNode struct {
	Service   string `json:"service"       description:"service name"`
	TradeNo   string `json:"tradeno"       description:"total refund price, unit one cent CNY"`
	CUUID     string `json:"cuuid"         description:"payer uuid"`
	SUUID     string `json:"suuid"         description:"payee uuid"`
	SubMchID  string `json:"sub_mchid"     description:"payee sub merchant id"`
	RefundID  string `json:"refund_id"     description:"refund id"`
	Total     int64  `json:"total"         description:"Original order price, unit one cent CNY"`
	RedundFee int64  `json:"refundfee"     description:"total refund price, unit one cent CNY"`
	Desc      string `json:"desc"          description:"this order description"`
	NotifyURL string `json:"notifyurl"     description:"the notify url to tell service that payment success"`
}

// TicketNode ticket node detail
type TicketNode struct {
	PayBody string `json:"paybody"`
	UpTime  int64  `json:"uptime"`
	Action  int64  `json:"action"`
}

// InTicketNo order system get order detail request
type InTicketNo struct {
	AID   string `json:"aid"`
	PayNo string `json:"payno"`
}

// InTicketData ticket node datas for generate request
type InTicketData struct {
	AID       string `json:"aid"`
	Encode    bool   `json:"encode"`
	PayBody   string `json:"paybody"`
	SignKey   string `json:"signkey"`
	Timestamp int64  `json:"timestamp"`
}

// InTicketMod ticket node datas for update request
type InTicketMod struct {
	AID       string `json:"aid"`
	PayBody   string `json:"paybody"`
	SignKey   string `json:"signkey"`
	Timestamp int64  `json:"timestamp"`
	PayNo     string `json:"payno"`
}

// PayInfo payment information
type PayInfo struct {
	PayWay    string `json:"payway"    description:"payment way, such as 'wechat', 'wechatJSAPI', 'alipay'"`
	Status    int64  `json:"status"    description:"payment status, such as 'cancle', 'unpaid', 'paid'"`
	WxPayInfo string `json:"wxpayinfo" description:"wechat payment app information"`
	AlPayInfo string `json:"alpayinfo" description:"alipay payment information"`
}
