// Copyright (c) 2018-2028 Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2021/12/26   youhei         New version
// -------------------------------------------------------------------

package mea

import (
	"fmt"
	"github.com/wengoldx/wing/comm"
	"github.com/wengoldx/wing/invar"
	"github.com/wengoldx/wing/logger"
)

func (m *MeaAgent) Measure(param *Measure) (string, error) {
	if m.Domain == "" {
		logger.E("Not set domain, please set first!")
		return "", invar.ErrInvalidClient
	}

	apiurl := m.Domain + "/measure/ai/capture"
	reqid, err := comm.HttpPostString(apiurl, param)
	if err != nil {
		logger.I("Send measure body request err", err)
		return "", err
	}
	return reqid, nil
}

func (m *MeaAgent) ReMeasure(param *BodyUpReq) error {
	if m.Domain == "" {
		logger.E("Not set domain, please set first!")
		return invar.ErrInvalidClient
	}

	apiurl := m.Domain + "/measure/ai/remeasure/capture"
	if _, err := comm.HttpPost(apiurl, param); err != nil {
		logger.I("Send re-measure body request err", err)
		return err
	}
	return nil
}

func (m *MeaAgent) Capture(reqid string) error {
	if m.Domain == "" {
		logger.E("Not set domain, please set first!")
		return invar.ErrInvalidClient
	}

	apiurl := m.Domain + "/measure/ai/recapture"
	param := &ReqID{ReqID: reqid}

	if _, err := comm.HttpPost(apiurl, param); err != nil {
		logger.I("Post re-capture body show picture err", err)
		return err
	}
	return nil
}

func (m *MeaAgent) BodyList(reqids []string) ([]*BodyBasicResp, error) {
	if m.Domain == "" {
		logger.E("Not set domain, please set first!")
		return nil, invar.ErrInvalidClient
	}

	apiurl := m.Domain + "/measure/ai/list"
	param := &ReqIDs{ReqIDs: reqids}
	resp := &([]*BodyBasicResp{})

	if err := comm.HttpPostStruct(apiurl, param, resp); err != nil {
		logger.I("Get body list data err", err)
		return nil, err
	}
	return *resp, nil
}

func (m *MeaAgent) BodyDetail(reqid string) (*BodyDetailResp, error) {
	if m.Domain == "" {
		logger.E("Not set domain, please set first!")
		return nil, invar.ErrInvalidClient
	}

	apiurl := fmt.Sprintf("%s/measure/ai/query?reqid=%s", m.Domain, reqid)
	resp := &BodyDetailResp{}

	if err := comm.HttpGetStruct(apiurl, resp); err != nil {
		logger.I("Get body detail data err", err)
		return nil, err
	}
	return resp, nil
}

func (m *MeaAgent) Delete(reqid string) error {
	if m.Domain == "" {
		logger.E("Not set domain, please set first!")
		return invar.ErrInvalidClient
	}

	apiurl := m.Domain + "/measure/ai/del"
	param := &ReqID{ReqID: reqid}

	if _, err := comm.HttpPost(apiurl, param); err != nil {
		logger.I("Delete body request err", err)
		return err
	}
	return nil
}
