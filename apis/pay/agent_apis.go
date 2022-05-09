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

import (
	"fmt"
	"github.com/wengoldx/wing/comm"
	"github.com/wengoldx/wing/invar"
	"github.com/wengoldx/wing/logger"
)

// Generate a new virtual card ticket, and return card number
//	@param ticket The first ticket node of trade
//	@return - string Trade transaction number
//			- error Exception message
func (a *PayAgent) GenCard(ticket *TradeNode) (string, error) {
	return a.postReqString("card", ticket)
}

// Generate a new trade by given payment datas, and return trade number
//	@param ticket The first ticket node of trade
//	@return - string Trade transaction number
//			- error Exception message
func (a *PayAgent) GenTrade(ticket *TradeNode) (string, error) {
	return a.postReqString("trade", ticket)
}

// Generate a new refund by given payment datas, and return trade number
//	@param ticket The first ticket node of refund
//	@return - string Refund transaction number
//			- error Exception message
func (a *PayAgent) GenRefund(ticket *RefundNode) (string, error) {
	return a.postReqString("refund", ticket)
}

// Change trade ticket amount which on TSUnpaid or TSPayError stauts
//	@param tno Trade transaction number
//	@param amount Dest trade amount to change
//	@return - error Exception message
func (a *PayAgent) ChangeTAmount(tno string, amount int64) error {
	return a.reqChangeAmount("ta", tno, amount)
}

// Change refund ticket amount which on TSInProgress or TSRefundError stauts
//	@param rno Refund transaction number
//	@param amount Dest refund amount to change
//	@return - error Exception message
func (a *PayAgent) ChangeRAmount(rno string, amount int64) error {
	return a.reqChangeAmount("ra", rno, amount)
}

// Revoke trade transaction by user
//	@param rno Refund transaction number
//	@return - error Exception message
func (a *PayAgent) RevokeTrade(tno string) error {
	return a.reqRevokeTicket("trade", tno)
}

// Revoke trade transaction by user
//	@param rno Refund transaction number
//	@return - error Exception message
func (a *PayAgent) RevokeRefund(rno string) error {
	return a.reqRevokeTicket("refund", rno)
}

// Update trade, it not modify the any exist tickt nodes but generate a new
// ticket and append to trade nodes list.
//	@param tno Trade transaction number
//	@param ticket The changed trade ticket node
//	@return - error Exception message
//
// `DEPRECATED`:
//
// This function is deprecate, use agent.ChangeTAmount() instead it to change
// ticket trade amount.
func (a *PayAgent) UpdateTrade(tno string, ticket *TradeNode) error {
	return a.postReqParams("trade", "tno", tno, ticket)
}

// Update refund, it not modify the any exist tickt nodes but generate a new
// ticket and append to refund nodes list.
//	@param rno Refund transaction number
//	@param ticket The changed refund ticket node
//	@return - error Exception message
//
// `DEPRECATED`:
//
// This function is deprecate, use agent.ChangeRAmount() instead it to change
// ticket refund amount.
func (a *PayAgent) UpdateRefund(rno string, ticket *RefundNode) error {
	return a.postReqParams("refund", "rno", rno, ticket)
}

// Get the latest trade ticket node
//	@param tno Trade transaction number
//	@return - TradeNode Trade ticket node
//			- error Exception message
func (a *PayAgent) TradeTicket(tno string) (*TradeNode, error) {
	resp := &TradeNode{}
	if err := a.getReqStruct("trade", "tno", tno, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// Get the latest dividing ticket node
//	@param tno Dividing transaction number
//	@return - DiviNode Dividing ticket node
//			- error Exception message
func (a *PayAgent) DiviTicket(tno string) (*DiviNode, error) {
	resp := &DiviNode{}
	if err := a.getReqStruct("dividing", "tno", tno, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// Get the latest refund ticket node
//	@param rno Refund transaction number
//	@return - RefundNode Refund ticket node
//			- error Exception message
func (a *PayAgent) RefundTicket(rno string) (*RefundNode, error) {
	resp := &RefundNode{}
	if err := a.getReqStruct("refund", "rno", rno, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// ------------------------
//    Internal Functions
// ------------------------

// Post http request with input params, then return transaction number as string
func (a *PayAgent) postReqString(api string, params interface{}) (string, error) {
	if a.Domain == "" {
		logger.E("Not set domain, please set first!")
		return "", invar.ErrInvalidClient
	}

	payapi := fmt.Sprintf("%s/wgpay/v2/chain/%s", a.Domain, api)
	ret, err := comm.HttpPostString(payapi, params)
	if err != nil {
		return "", err
	}
	return ret, nil
}

// Post http request with input params, it will append key and value into request url
func (a *PayAgent) postReqParams(api, key, val string, params interface{}) error {
	if a.Domain == "" {
		logger.E("Not set domain, please set first!")
		return invar.ErrInvalidClient
	}

	if key == "" || val == "" {
		logger.E("Invalid key:", key, "or value:", val)
		return invar.ErrInvalidParams
	}

	payapi := fmt.Sprintf("%s/wgpay/v2/chain/update/%s?%s=%s", a.Domain, api, key, val)
	if _, err := comm.HttpPost(payapi, params); err != nil {
		return err
	}
	return nil
}

// Post http get request after append key and value into request url,
// then return  struct data from response
func (a *PayAgent) getReqStruct(api, key, val string, resp interface{}) error {
	if a.Domain == "" {
		logger.E("Not set domain, please set first!")
		return invar.ErrInvalidClient
	}

	if key == "" || val == "" {
		logger.E("Invalid key:", key, "or value:", val)
		return invar.ErrInvalidParams
	}

	payapi := fmt.Sprintf("%s/wgpay/v2/chain/ticket/%s?%s=%s", a.Domain, api, key, val)
	if err := comm.HttpGetStruct(payapi, resp); err != nil {
		return err
	}
	return nil
}

// Post http get request to change ticket amount
func (a *PayAgent) reqChangeAmount(api, tid string, amount int64) error {
	if a.Domain == "" {
		logger.E("Not set domain, please set first!")
		return invar.ErrInvalidClient
	}

	if tid == "" || amount <= 0 {
		logger.E("Invalid ticket id:", tid, "or amount:", amount)
		return invar.ErrInvalidParams
	}

	payapi := fmt.Sprintf("%s/wgpay/v2/chain/update/%s?tid=%s&amount=%v", a.Domain, api, tid, amount)
	if _, err := comm.HttpGet(payapi); err != nil {
		return err
	}
	return nil
}

// Post http get request to revoke ticket
func (a *PayAgent) reqRevokeTicket(api, tid string) error {
	if a.Domain == "" {
		logger.E("Not set domain, please set first!")
		return invar.ErrInvalidClient
	}

	if tid == "" {
		logger.E("Invalid ticket id:", tid, "to revoke")
		return invar.ErrInvalidParams
	}

	payapi := fmt.Sprintf("%s/wgpay/v2/chain/revoke/%s?tid=%s", a.Domain, api, tid)
	if _, err := comm.HttpGet(payapi); err != nil {
		return err
	}
	return nil
}
