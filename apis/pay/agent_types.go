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
	Cashier    string `json:"cashier" validate:"required" description:"cashier name who provide transaction by wgpay server"`
	Payer      string `json:"parer"   validate:"required" description:"payer unique id, such as user uuid"`
	Payee      string `json:"payee,omitempty"             description:"payee unique id, such as merchant id"`
	SMchID     string `json:"smid,omitempty"              description:"sub merchant id of payee"`
	Amount     int64  `json:"amount"  validate:"required" description:"total amount price, unit one cent CNY"`
	Desc       string `json:"desc"    validate:"required" description:"current trade ticket descriptions"`
	NotifyURL  string `json:"ntfurl,omitempty"            description:"ansync notifier url from wgpay to notify trade status changed, must returen OK if success to stop wgpay notify event looper"`
	PayWay     string `json:"payway,omitempty"            description:"payment way, such as 'wehcat', 'wechatJSAPI' and 'alipay'"`
	IsFrozen   bool   `json:"isfrozen,omitempty"          description:"whether frozen amount when payment completed, it must be true for dividing payment"`
	Status     string `json:"status"  validate:"required" description:"payment status, such as 'UNPAID', 'PAY_ERROR', 'REVOKED', 'PAID', 'COMPLETED', 'CLOSED'"`
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
	Amount    int64  `json:"amount"   validate:"required" description:"total amount price completed transaction, unit one cent CNY"`
	Refund    int64  `json:"refund"   validate:"required" description:"total refund price should return back to payer, unit one cent CNY"`
	Desc      string `json:"desc"     validate:"required" description:"current refund transacte description"`
	Status    string `json:"status"   validate:"required" description:"refund status, such as 'REFUND_IN_PROGRESS', 'REFUND_ERROR', 'REFUND', 'CLOSED'"`
	NotifyURL string `json:"ntfurl,omitempty"             description:"ansync notifier url from wgpay to notify refund status changed, must return OK if success to stop wgpay notify event looper"`
}

// PayInfo payment information
type PayInfo struct {
	PayWay  string `json:"payway"    description:"payment way, such as 'wechat', 'wechatJSAPI', 'alipay'"`
	Status  string `json:"status"    description:"payment status, such as 'PLACE_ORDER', 'UNPAID', 'PAY_ERROR', 'REVOKED', 'PAID', 'COMPLETED', 'REFUND_IN_PROGRESS', 'REFUND_ERROR', 'REFUND', 'CLOSED'"`
	WxInfo  string `json:"wxpayinfo" description:"wechat payment app information"`
	AliInfo string `json:"alpayinfo" description:"alipay payment information"`
}
