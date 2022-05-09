// Copyright (c) 2018-2028 Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2019/05/22   yangping       New version
// -------------------------------------------------------------------

package mvc

type InToken struct {
	Token string `json:"token" validate:"required" description:"jwt auth token"`
}

type InID struct {
	ID int64 `json:"id" validate:"required" description:"target id, must over 0"`
}

type InUID struct {
	UID string `json:"uuid" validate:"required" description:"user unique id"`
}

type InCode struct {
	Code string `json:"code" validate:"required" description:"secure code string"`
}

type InTNO struct {
	TradeNo string `json:"trade_no" validate:"required" description:"trade no string"`
}

type InOID struct {
	Token  string `json:"token"  validate:"required" description:"jwt auth token"`
	OpenID string `json:"openid" validate:"required" description:"wechat openid"`
}

type InTokenID struct {
	Token string `json:"token" validate:"required" description:"jwt auth token"`
	ID    int64  `json:"id"    validate:"required" description:"target id, must over 0"`
}

type InTokenUID struct {
	Token string `json:"token" validate:"required" description:"jwt auth token"`
	UID   string `json:"uuid"  validate:"required" description:"user unique id"`
}

type InTCode struct {
	Token string `json:"token" validate:"required" description:"jwt auth token"`
	Code  string `json:"code"  validate:"required" description:"secure code string"`
}

type InTrade struct {
	Token   string `json:"token"    validate:"required" description:"jwt auth token"`
	TradeNo string `json:"trade_no" validate:"required" description:"trade no string"`
}
