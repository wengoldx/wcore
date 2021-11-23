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

// requestOrderDetail get the latest order detail by given trade number
func (a *PaychainAgent) requestOrderDetail(url, tradeno string, out interface{}) error {
	orderInfo, err := a.OrderDetail(url, tradeno)
	if err != nil {
		return err
	}

	length := len(orderInfo)
	if length == 0 || orderInfo[length-1] == nil {
		logger.E("Invalid responsed orders detail!")
		return invar.ErrInvalidData
	}

	lastOrder := orderInfo[length-1]
	if err := json.Unmarshal([]byte(lastOrder.PayBody), out); err != nil {
		logger.E("Unmarshal order body err:", err)
		return err
	}
	return nil
}

// TradeNoDetail get the latest trade order detail
func (a *PaychainAgent) TradeNoDetail(url, tradeno string) (*OrderBodyInfo, error) {
	logger.D("Get trade order detail by no:", tradeno)
	tradeInfo := &OrderBodyInfo{}
	if err := a.requestOrderDetail(url, tradeno, tradeInfo); err != nil {
		return nil, err
	}
	return tradeInfo, nil
}

// ShareNoDetail get the latest share order detail
func (a *PaychainAgent) ShareNoDetail(url, tradeno string) (*ProfitShareInfo, error) {
	logger.D("Get share order detail by no:", tradeno)
	shareInfo := &ProfitShareInfo{}
	if err := a.requestOrderDetail(url, tradeno, shareInfo); err != nil {
		return nil, err
	}
	return shareInfo, nil
}

// RefundNoDetail get the latest refund order detail
func (a *PaychainAgent) RefundNoDetail(url, tradeno string) (*RefundBodyInfo, error) {
	logger.D("Get refund order detail by no:", tradeno)
	refundInfo := &RefundBodyInfo{}
	if err := a.requestOrderDetail(url, tradeno, refundInfo); err != nil {
		return nil, err
	}
	return refundInfo, nil
}

// UpdateOrderBody update trade number information by struct
func (a *PaychainAgent) UpdateOrderBody(url, tradeno string, ps interface{}) error {
	body, err := json.Marshal(ps)
	if err != nil {
		logger.E("Marshal order body struct err:", err)
		return err
	}

	logger.D("Update order body no:", tradeno)
	if err = a.OrderUpdate(url, tradeno, string(body)); err != nil {
		return err
	}
	return nil
}

// GenerateOrderBody generate trade body by struct, and return trade number
func (a *PaychainAgent) GenerateOrderBody(url string, ps interface{}) (string, error) {
	body, err := json.Marshal(ps)
	if err != nil {
		logger.E("Marshal order body struct err:", err)
		return "", err
	}

	logger.D("Generate order body ", string(body))
	tradeno, err := a.OrderGen(url, string(body))
	if err != nil {
		return "", err
	}

	return tradeno, nil
}

// OrderDetail get the order detail by tardeno from paychain server
//	@TODO this method should use rpc instead of http post.
func (a *PaychainAgent) OrderDetail(url, tradeno string) ([]*OrderDetailResp, error) {
	params := &OrderDetailReq{
		AID:   a.Aid,
		PayNo: tradeno,
	}

	respByte, err := comm.HttpPost(url, params)
	if err != nil {
		logger.E("Post request order detail err:", err)
		return nil, err
	}

	resp := []*OrderDetailResp{}
	if err = json.Unmarshal(respByte, &resp); err != nil {
		logger.E("Unmarshal order detail err:", err)
		return nil, err
	}

	logger.D("Got order:", tradeno, "detail")
	return resp, nil
}

// OrderUpdate update order information to paychain server
//	@TODO this method should use rpc instead of http post.
func (a *PaychainAgent) OrderUpdate(url, tradeno, body string) error {
	signkey, eb, ts, err := a.Encrypt(body)
	if err != nil {
		logger.E("Encrypt order body err:", err)
		return err
	}

	params := &OrderModReq{
		AID:       a.Aid,
		PayBody:   eb,
		SignKey:   signkey,
		Timestamp: ts,
		PayNo:     tradeno,
	}
	if _, err = comm.HttpPostString(url, params); err != nil {
		logger.E("Post request update order err:", err)
		return err
	}

	logger.D("Updated order:", tradeno, "detail")
	return nil
}

// OrderGen generate the out request number
// 	@TODO this method should use rpc instead of http post.
func (a *PaychainAgent) OrderGen(url, body string) (string, error) {
	signkey, eb, ts, err := a.Encrypt(body)
	if err != nil {
		logger.E("Encrypt order body err:", err)
		return "", err
	}

	params := &OrderGenReq{
		AID:       a.Aid,
		Encode:    true,
		PayBody:   eb,
		SignKey:   signkey,
		Timestamp: ts,
	}

	tradeno, err := comm.HttpPostString(url, params)
	if err != nil {
		logger.E("Post request generate order err:", err)
		return "", err
	}

	logger.D("Generate order:", tradeno)
	return tradeno, nil
}

/**
 * Encrypt encrypt the given body, and return sign code, timestamp and body ciphertext.
 *
 * @param body payment content to be encrypted.
 * @return string sign code of hashed payment body.
 * @return string encrypted body ciphertext.
 * @return int64  encrypte timestamp.
 * @return error  invar.ErrInvalidClient or exception errors.
 *
 * WARNING : the body string max lenght DO NOT lagger than 400 chars.
 */
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
