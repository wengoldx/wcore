// Copyright (c) 2018-2022 WING All Rights Reserved.
//
// Author : jidi
// Email  : j18041361158@163.com
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2021/11/17   jidi           New version
// -------------------------------------------------------------------

package order

import (
	"encoding/json"
	"github.com/wengoldx/wing/comm"
	"github.com/wengoldx/wing/invar"
	"github.com/wengoldx/wing/logger"
	"github.com/wengoldx/wing/secure"
	"time"
)

// Generate a new trade by given payment datas, and return trade number
//	@param ps  "The first ticket node of trade"
//	@return - string "Trade number"
//			- error "Exception messages"
func (a *PaychainAgent) GenTicket(ps interface{}) (string, error) {
	if a.Domain == "" {
		logger.E("Not set domain, please set first!")
		return "", invar.ErrInvalidClient
	}

	body, err := json.Marshal(ps)
	if err != nil {
		logger.E("Marshal ticket node err:", err)
		return "", err
	}

	logger.D("Generate trade the first ticket:", string(body))
	tradeno, err := a.genTicketNode(a.Domain+"/gen", string(body))
	if err != nil {
		return "", err
	}
	return tradeno, nil
}

// Get the latest trade ticket node
//	@param tno "Trade number"
//	@return - TradeNode "Trade ticket node"
//			- error "Exception messages"
func (a *PaychainAgent) TradeTicket(tno string) (*TradeNode, error) {
	if a.Domain == "" {
		logger.E("Not set domain, please set first!")
		return nil, invar.ErrInvalidClient
	}

	ticket := &TradeNode{}
	logger.D("Get trade ticket by no:", tno)
	if err := a.lastTicketNode(tno, ticket); err != nil {
		return nil, err
	}
	return ticket, nil
}

// Get the latest dividing ticket node
//	@param tno "Trade number"
//	@return - ProfitShareInfo "Dividing ticket node"
//			- error "Exception messages"
func (a *PaychainAgent) DiviTicket(tno string) (*DiviNode, error) {
	if a.Domain == "" {
		logger.E("Not set domain, please set first!")
		return nil, invar.ErrInvalidClient
	}

	ticket := &DiviNode{}
	logger.D("Get dividing ticket by no:", tno)
	if err := a.lastTicketNode(tno, ticket); err != nil {
		return nil, err
	}
	return ticket, nil
}

// Get the latest refund ticket node
//	@param tno "Trade number"
//	@return - RefundNode "Refund ticket node"
//			- error "Exception messages"
func (a *PaychainAgent) RefundTicket(tno string) (*RefundNode, error) {
	if a.Domain == "" {
		logger.E("Not set domain, please set first!")
		return nil, invar.ErrInvalidClient
	}

	ticket := &RefundNode{}
	logger.D("Get refund ticket by no:", tno)
	if err := a.lastTicketNode(tno, ticket); err != nil {
		return nil, err
	}
	return ticket, nil
}

// Update trade, it not modify the any exist tickt nodes but generate a new
// ticket and append to trade nodes list.
//	@param tno "Trade number"
//	@param ps  "The new trade ticket node"
//	@return - error "Exception messages"
func (a *PaychainAgent) UpdateTicket(tno string, ps interface{}) error {
	if a.Domain == "" {
		logger.E("Not set domain, please set first!")
		return invar.ErrInvalidClient
	}

	node, err := json.Marshal(ps)
	if err != nil {
		logger.E("Marshal ticket ndoe err:", err)
		return err
	}

	logger.D("Update ticket node by trade no:", tno)
	if err = a.appendTradeTicket(tno, string(node)); err != nil {
		return err
	}
	return nil
}

// ------------------------
//    Internal Functions
// ------------------------

// Get the last ticket node from given trade nodes list
//	@param tno "Trade number"
//	@return - out "Out data of TradeNode, DiviNode or RefundNode"
//			- error "Exception messages"
func (a *PaychainAgent) lastTicketNode(tno string, out interface{}) error {
	tickets, err := a.getTradeTickets(tno)
	if err != nil {
		return err
	}

	length := len(tickets)
	if length == 0 || tickets[length-1] == nil {
		logger.E("Invalid responsed tickets nodes!")
		return invar.ErrInvalidData
	}

	ticket := tickets[length-1]
	if err := json.Unmarshal([]byte(ticket.PayBody), out); err != nil {
		logger.E("Unmarshal last ticket node err:", err)
		return err
	}
	return nil
}

// Get trade tickets list by trade number from paychain server
//	@param tno "Trade number"
//	@return - []TicketNode "Trade tickets nodes"
//			- error "Exception messages"
//
// `TODO`
//
// This method should use rpc instead of http post.
func (a *PaychainAgent) getTradeTickets(tno string) ([]*TicketNode, error) {
	params := &InTicketNo{AID: a.Aid, PayNo: tno}
	respByte, err := comm.HttpPost(a.Domain+"/detail", params)
	if err != nil {
		logger.E("Request trade tickets err:", err)
		return nil, err
	}

	resp := []*TicketNode{}
	if err = json.Unmarshal(respByte, &resp); err != nil {
		logger.E("Unmarshal trade tickets err:", err)
		return nil, err
	}

	logger.D("Got trade:", tno, "tickets")
	return resp, nil
}

// Append a new ticket node to paychain server
//	@param tno  "Trade number"
//	@param node "Json string of the new ticket node"
//	@return - error "Exception messages"
//
// `TODO`
//
// This method should use rpc instead of http post.
func (a *PaychainAgent) appendTradeTicket(tno, node string) error {
	signkey, eb, ts, err := a.Encrypt(node)
	if err != nil {
		logger.E("Encrypt ticket node err:", err)
		return err
	}

	params := &InTicketMod{
		AID:       a.Aid,
		PayBody:   eb,
		SignKey:   signkey,
		Timestamp: ts,
		PayNo:     tno,
	}
	if _, err = comm.HttpPostString(a.Domain+"/mod", params); err != nil {
		logger.E("Post update trade err:", err)
		return err
	}

	logger.D("Updated trade:", tno, "tickets")
	return nil
}

// Generate a new ticket node as the first node of trade nodes list,
// and return the trade number
//	@param url  "Paychain API router"
//	@param node "Json string of the new ticket node"
//	@return - string "Trade number"
//			- error "Exception messages"
//
// `TODO` :
//
// This method should use rpc instead of http post.
func (a *PaychainAgent) genTicketNode(url, node string) (string, error) {
	signkey, eb, ts, err := a.Encrypt(node)
	if err != nil {
		logger.E("Encrypt ticket node err:", err)
		return "", err
	}

	params := &InTicketData{
		AID:       a.Aid,
		Encode:    true,
		PayBody:   eb,
		SignKey:   signkey,
		Timestamp: ts,
	}

	tno, err := comm.HttpPostString(url, params)
	if err != nil {
		logger.E("Post generate ticket node err:", err)
		return "", err
	}

	logger.D("Generate ticket:", tno)
	return tno, nil
}

//	Encrypt encrypt the given body, and return sign code, timestamp and body ciphertext.
//
//	@param body payment content to be encrypted.
//	@return string sign code of hashed payment body.
//	@return string encrypted body ciphertext.
//	@return int64  encrypte timestamp.
//	@return error  invar.ErrInvalidClient or exception errors.
//
// `WARNING` :
//
// The body string max lenght DO NOT lagger than 400 chars.
func (a *PaychainAgent) Encrypt(body string) (string, string, int64, error) {
	hashcode := secure.EncodeMD5(body)
	timestamp := time.Now().UnixNano()
	bodybytes, err := json.Marshal(&EncryptNode{
		SecureKey: secure.EncodeMD5(a.Devmac), Timestamp: timestamp, PayBody: body,
	})
	if err != nil {
		return "", "", 0, err
	}

	ciphertext, err := secure.RSAEncrypt([]byte(a.Pubkey), bodybytes)
	if err != nil {
		return "", "", 0, err
	}

	ciphertextb64 := secure.EncodeBase64(string(ciphertext))
	return hashcode, ciphertextb64, timestamp, nil
}
