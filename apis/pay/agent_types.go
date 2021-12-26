// Copyright (c) 2019-2029 DY All Rights Reserved.
//
// Author : yangping
// Email  : youhei_yp@163.com
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2021/12/26   youhei         New version
// -------------------------------------------------------------------

package pay

type PayAgent struct {
	// Wgpay service access url domain, it allways get from app.conf
	// by beego.AppConfig.String("domain") code.
	Domain string
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
