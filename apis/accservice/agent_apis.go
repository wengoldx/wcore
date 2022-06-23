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

import (
	"github.com/wengoldx/wing/comm"
	"github.com/wengoldx/wing/invar"
	"github.com/wengoldx/wing/logger"
)

// Verify token from accserveice service
//	@param token User logined jwt token
//	@return - string Account uuid
//			- error Exception message
func (a *AccAgent) AuthToken(token string) (string, error) {
	if a.Domain == "" {
		logger.E("Not set domain, please set first!")
		return "", invar.ErrInvalidClient
	}

	apiurl := a.Domain + "/accservice/acc/via/token"
	params := &Token{Token: token}

	resp, err := comm.HttpPostString(apiurl, params)
	if err != nil {
		logger.E("Verify account token, err:", err)
		return "", err
	}

	return resp, nil
}

// Get user base profile from accserveice service
//	@param uids Account uuids
//	@return - []ProfSumm Accounts simple profiles
//			- error Exception message
func (a *AccAgent) ProfBases(uids []string) ([]ProfSumm, error) {
	if a.Domain == "" {
		logger.E("Not set domain, please set first!")
		return nil, invar.ErrInvalidClient
	}

	resp := &([]ProfSumm{})
	if len(uids) > 0 {
		apiurl := a.Domain + "/accservice/prof/summs"
		params := &UUIDs{UID: uids}

		if err := comm.HttpPostStruct(apiurl, params, resp); err != nil {
			logger.E("Get users base profile, err:", err)
			return nil, err
		}
	}

	return *resp, nil
}
