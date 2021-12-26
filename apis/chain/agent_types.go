// Copyright (c) 2019-2029 DY All Rights Reserved.
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
	CANCELLED                    // trade was cancelled, status 2
	UNPAID                       // trade not paid or paied fail, status 3
	PAID                         // goods to be delivered, status 4
)

// Agent of paychain to simple access server RESTful APIs
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
//	agentIns := &chain.PaychainAgent {
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
type ChainAgent struct {
	Aid    string // agent id
	Devmac string // device mac bind with current agent
	Pubkey string // public key of current agent

	// Paychain service access url domain, it allways get from app.conf
	// by beego.AppConfig.String("domain") code.
	Domain string
}

// Encrypt node data struct
type EncryptNode struct {
	SecureKey string `json:"securekey"`
	Timestamp int64  `json:"timestamp"`
	PayBody   string `json:"paybody"`
}

// Trade ticket node
type TradeNode struct {
	Cashier    string `json:"service"               description:"cashier name who provide transaction by wgpay server"`
	Payer      string `json:"cuuid"                 description:"payer uuid"`
	Payee      string `json:"suuid"                 description:"payee uuid, merchant id"`
	SMchID     string `json:"sub_mchid"             description:"sub merchant id of payee"`
	Amount     int64  `json:"amount"                description:"total amount price, unit one cent CNY"`
	Refund     int64  `json:"refundfee"             description:"total refund price, unit one cent CNY"`
	Desc       string `json:"desc"                  description:"this ticket description"`
	NotifyURL  string `json:"notifyurl"             description:"ansync notifier url from wgpay to notify pay status changed, must returen OK if success"`
	PayWay     string `json:"payway"                description:"payment way, such as 'wehcat', 'wechatJSAPI' and 'alipay'"`
	IsFrozen   bool   `json:"isfrozen"              description:"whether frozen amount when payment finished, it must be true for dividing payment"`
	Status     int64  `json:"status"                description:"payment status, such as 'cancle', 'unpaid', 'paid'"`
	TimeExpire string `json:"time_expire,omitempty" description:"expire time for virture products such as coupon, courtesy card, and so on"`
}

// Dividing ticket node
type DiviNode struct {
	Cashier    string `json:"service"        description:"cashier name who provide transaction by wgpay server"`
	SMchID     string `json:"sub_mchid"      description:"sub merchant id of payee"`
	TranID     string `json:"transaction_id" description:"transaction id of wechat pay platform"`
	Commission int64  `json:"commission"     description:"commission of dividing transaction, unit one cent CNY"`
	Desc       string `json:"desc"           description:"dividing transacte description"`
	IsFinsh    bool   `json:"isfinsh"        description:"finish transation, and unfrozen transation"`
}

// Refund ticket node
type RefundNode struct {
	Cashier   string `json:"service"   description:"cashier name who provide transaction by wgpay server"`
	TranNo    string `json:"tradeno"   description:"transaction number of mall pay platform"`
	Payer     string `json:"cuuid"     description:"payer uuid"`
	Payee     string `json:"suuid"     description:"payee uuid"`
	SMchID    string `json:"sub_mchid" description:"sub merchant id of payee"`
	RefundID  string `json:"refund_id" description:"refund id"`
	Amount    int64  `json:"total"     description:"total amount price, unit one cent CNY"`
	Refund    int64  `json:"refundfee" description:"total refund price, unit one cent CNY"`
	Desc      string `json:"desc"      description:"refund transacte description"`
	Status    int64  `json:"status"    description:"refund status, such as 'cancle', 'unpaid', 'paid'"`
	NotifyURL string `json:"notifyurl" description:"ansync notifier url from wgpay to notify refund status changed, must return OK if success"`
}

// Base ticket node to save into paychain database
type TicketNode struct {
	PayBody string `json:"paybody"`
	UpTime  int64  `json:"uptime"`
	Action  int64  `json:"action"`
}

// Request params contain agent id and transaction number
type ChainNo struct {
	AID string `json:"aid"`
	TNo string `json:"payno"`
}

// Request params to save trade ticket as packet into paychain database when start a trade
type ChainData struct {
	AID       string `json:"aid"`
	Encode    bool   `json:"encode"`
	PayBody   string `json:"paybody"`
	SignKey   string `json:"signkey"`
	Timestamp int64  `json:"timestamp"`
}

// Request params to save transaction tickt as packet into paychain database when update a trade both chande or refund
type ChainMod struct {
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
