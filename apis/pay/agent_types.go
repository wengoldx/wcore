// Copyright (c) 2018-2028 Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2021/12/26   youhei         New version
// -------------------------------------------------------------------

package pay

type PayAgent struct {
	Domain string // Wgpay server access url, such as http://192.168.1.100:3000
}

// Trade ticket node
type TradeNode struct {
	Cashier    string `json:"cashier" validate:"required" description:"cashier name who provide transaction by wgpay server"`
	Payer      string `json:"parer"   validate:"required" description:"payer unique id, such as user uuid"`
	Payee      string `json:"payee,omitempty"             description:"payee unique id, such as merchant id"`
	SMchID     string `json:"smid,omitempty"              description:"sub merchant id of payee"`
	Amount     int64  `json:"amount"  validate:"required" description:"total amount price, unit one cent CNY"`
	Desc       string `json:"desc"    validate:"required" description:"current trade ticket descriptions"`
	NotifyURL  string `json:"ntfurl,omitempty"            description:"ansync notifier url from wgpay to notify trade status changed, must returen OK if success to stop wgpay notify event looper"`
	PayWay     string `json:"payway,omitempty"            description:"payment way, such as 'wehcat', 'wechatJSAPI' and 'alipay'"`
	IsFrozen   bool   `json:"isfrozen,omitempty"          description:"whether frozen amount when payment completed, it must be true for dividing payment"`
	Status     string `json:"status"                      description:"payment status, such as 'UNPAID', 'PAY_ERROR', 'REVOKED', 'PAID', 'COMPLETED', 'CLOSED'"`
	TimeExpire string `json:"expire,omitempty"            description:"expire time for virture products such as coupon, courtesy card, and so on"`
}

// Dividing ticket node
type DiviNode struct {
	Cashier    string `json:"cashier" validate:"required" description:"cashier name who provide transaction by wgpay server"`
	SMchID     string `json:"smid"                        description:"sub merchant id of payee"`
	TranID     string `json:"transaction_id"              description:"transaction id of wechat pay platform"`
	Commission int64  `json:"commission"                  description:"commission of dividing transaction, unit one cent CNY"`
	Desc       string `json:"desc"    validate:"required" description:"dividing transacte description"`
	IsFinsh    bool   `json:"isfinsh"                     description:"finish transation, and unfrozen transation"`
}

// Refund ticket node
type RefundNode struct {
	Cashier   string `json:"cashier"  validate:"required" description:"cashier name who provider refund by wgpay server"`
	TranNo    string `json:"trade_no" validate:"required" description:"the original trade transaction number"`
	Payer     string `json:"payer"    validate:"required" description:"payer unique id, such as user uuid"`
	Payee     string `json:"payee,omitempty"              description:"payee unique id, such as merchant id"`
	SMchID    string `json:"smid,omitempty"               description:"sub merchant id of payee"`
	RefundID  string `json:"refund_id,omitempty"          description:"refund transaction id of wechat pay"`
	Amount    int64  `json:"amount"   validate:"required" description:"total amount price to refund, unit one cent CNY"`
	Desc      string `json:"desc,omitempty"               description:"current refund transacte description"`
	Status    string `json:"status,omitempty"             description:"refund status, such as 'REFUND_IN_PROGRESS', 'REFUND_ERROR', 'REFUND', 'CLOSED'"`
	NotifyURL string `json:"ntfurl,omitempty"             description:"ansync notifier url from wgpay to notify refund status changed, must return OK if success to stop wgpay notify event looper"`
}

// PayInfo payment information
type PayInfo struct {
	PayWay  string `json:"payway"    description:"payment way, such as 'wechat', 'wechatJSAPI', 'alipay'"`
	Status  string `json:"status"    description:"payment status, such as 'PLACE_ORDER', 'UNPAID', 'PAY_ERROR', 'REVOKED', 'PAID', 'COMPLETED', 'REFUND_IN_PROGRESS', 'REFUND_ERROR', 'REFUND', 'CLOSED'"`
	WxInfo  string `json:"wxpayinfo" description:"wechat payment app information"`
	AliInfo string `json:"alpayinfo" description:"alipay payment information"`
}

// Combine Trade ticket node
type CombineNode struct {
	Cashier    string       `json:"cashier"   validate:"required"        description:"cashier name who provide transaction by wgpay server"`
	TimeExpire string       `json:"expire,omitempty"                     description:"expire time for virture products such as coupon, courtesy card, and so on"`
	SubOrders  []*TradeNode `json:"sub_order" validate:"required,max=10" description:"the array of sub trade number"`
}
