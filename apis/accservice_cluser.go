// Copyright (c) 2018-2028 Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2021/12/06   jidi           New version
// -------------------------------------------------------------------

package apis

import (
	"github.com/astaxie/beego"
	"github.com/wengoldx/wing/comm"
	"github.com/wengoldx/wing/logger"
)

const (
	agent_router = "/accservice"
)

type token struct {
	Token string
}

type uuids struct {
	UUID []string
}

// ProfSumm base personal information
type ProfSumm struct {
	UUID     string `json:"uuid"     description:"account uuid(account number)"`
	Sex      int    `json:"sex"      description:"account sex"`
	Nickname string `json:"nickname" description:"user nickname"`
	HeardURL string `json:"heardurl" description:"user heard image url"`
}

// ApiAuthAccToken verify token from accserveice service
func ApiAuthAccToken(param string) (string, error) {
	apiurl := beego.AppConfig.String("domain") + agent_router + "/acc/via/token"
	params := &token{Token: param}

	res, err := comm.HttpPostString(apiurl, params)
	if err != nil {
		logger.E("ApiAuthAccToken verify token to accservice service err:", err)
		return "", err
	}

	return res, nil
}

// ApiAccGetProBases batch get user base profile, and support to get the basic information of a user at the same time
func ApiAccGetProBases(uid []string) ([]ProfSumm, error) {
	apiurl := beego.AppConfig.String("domain") + agent_router + "/prof/summs"
	params := &uuids{UUID: uid}
	pro := &([]ProfSumm{})

	if len(uid) > 0 {
		if err := comm.HttpPostStruct(apiurl, params, pro); err != nil {
			logger.E("ApiAccGetProBases get users base profile information from accservice err:", err)
			return nil, err
		}
	}

	return *pro, nil
}
