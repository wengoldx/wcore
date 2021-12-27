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

import (
	"github.com/wengoldx/wing/comm"
	"github.com/wengoldx/wing/invar"
	"github.com/wengoldx/wing/logger"
)

// Generate a new trade by given payment datas, and return trade number
//	@param ticket The first ticket node of trade
//	@return - string Trade transaction number
//			- error Exception message
func (a *PayAgent) GenTrade(ticket *TradeNode) (string, error) {
	return a.postReqString("/wgpay/v2/chain/trade", ticket)
}

// Generate a new refund by given payment datas, and return trade number
//	@param ticket The first ticket node of refund
//	@return - string Refund transaction number
//			- error Exception message
func (a *PayAgent) GenRefund(ticket *RefundNode) (string, error) {
	return a.postReqString("/wgpay/v2/chain/refund", ticket)
}

// Update trade, it not modify the any exist tickt nodes but generate a new
// ticket and append to trade nodes list.
//	@param tno Trade transaction number
//	@param ticket The changed trade ticket node
//	@return - error Exception message
func (a *PayAgent) UpdateTrade(tno string, ticket *TradeNode) error {
	return a.postReqParams("/wgpay/v2/chain/update/trade", "tno", tno, ticket)
}

// Update refund, it not modify the any exist tickt nodes but generate a new
// ticket and append to refund nodes list.
//	@param rno Refund transaction number
//	@param ticket The changed refund ticket node
//	@return - error Exception message
func (a *PayAgent) UpdateRefund(rno string, ticket *RefundNode) error {
	return a.postReqParams("/wgpay/v2/chain/update/refund", "rno", rno, ticket)
}

// Get the latest trade ticket node
//	@param tno Trade transaction number
//	@return - TradeNode Trade ticket node
//			- error Exception message
func (a *PayAgent) TradeTicket(tno string) (*TradeNode, error) {
	resp, apiurl := &TradeNode{}, "/wgpay/v2/chain/ticket/trade"
	if err := a.getReqStruct(apiurl, "tno", tno, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// Get the latest dividing ticket node
//	@param tno Dividing transaction number
//	@return - DiviNode Dividing ticket node
//			- error Exception message
func (a *PayAgent) DiviTicket(tno string) (*DiviNode, error) {
	resp, apiurl := &DiviNode{}, "/wgpay/v2/chain/ticket/dividing"
	if err := a.getReqStruct(apiurl, "tno", tno, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// Get the latest refund ticket node
//	@param rno Refund transaction number
//	@return - RefundNode Refund ticket node
//			- error Exception message
func (a *PayAgent) RefundTicket(rno string) (*RefundNode, error) {
	resp, apiurl := &RefundNode{}, "/wgpay/v2/chain/ticket/refund"
	if err := a.getReqStruct(apiurl, "rno", rno, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// ------------------------
//    Internal Functions
// ------------------------

// Post http request with input params, then return transaction number as string
func (a *PayAgent) postReqString(apiurl string, params interface{}) (string, error) {
	if a.Domain == "" {
		logger.E("Not set domain, please set first!")
		return "", invar.ErrInvalidClient
	}

	payapi := a.Domain + apiurl
	ret, err := comm.HttpPostString(payapi, params)
	if err != nil {
		return "", err
	}
	return ret, nil
}

// Post http request with input params, it will append key and value into request url
func (a *PayAgent) postReqParams(apiurl, key, val string, params interface{}) error {
	if a.Domain == "" {
		logger.E("Not set domain, please set first!")
		return invar.ErrInvalidClient
	}

	if key == "" || val == "" {
		logger.E("Invalid key:", key, "or value:", val)
		return invar.ErrInvalidClient
	}

	payapi := a.Domain + apiurl + "?" + key + "=" + val
	if _, err := comm.HttpPost(payapi, params); err != nil {
		return err
	}
	return nil
}

// Post http get request after append key and value into request url,
// then return  struct data from response
func (a *PayAgent) getReqStruct(apiurl, key, val string, resp interface{}) error {
	if a.Domain == "" {
		logger.E("Not set domain, please set first!")
		return invar.ErrInvalidClient
	}

	if key == "" || val == "" {
		logger.E("Invalid key:", key, "or value:", val)
		return invar.ErrInvalidClient
	}

	payapi := a.Domain + apiurl + "?" + key + "=" + val
	if err := comm.HttpGetStruct(payapi, resp); err != nil {
		return err
	}
	return nil
}
