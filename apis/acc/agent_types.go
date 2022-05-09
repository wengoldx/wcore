// Copyright (c) 2018-2028 Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2021/12/26   youhei         New version
// -------------------------------------------------------------------

package acc

type AccAgent struct {
	// Wgpay service access url domain, it allways get from app.conf
	// by beego.AppConfig.String("domain") code.
	Domain string
}

type Token struct {
	Token string `json:"token" description:"jwt auth token"`
}

type UUIDs struct {
	UID []string `json:"uuid" description:"account uuid(account number)"`
}

// Account simple profile datas
type ProfSumm struct {
	UUID     string `json:"uuid"     description:"account uuid(account number)"`
	Sex      int    `json:"sex"      description:"account sex, 0:none, 1:male, 2:female, 3:neutral"`
	Nickname string `json:"nickname" description:"user nickname"`
	HeardURL string `json:"heardurl" description:"user heard image url"`
}
