// Copyright (c) 2019-2029 Dunyu All Rights Reserved.
//
// Author      : jidi
// Email       : j18041361158@163.com
// Version     : 1.0
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2019/12/16   jidi           New version
// 00002       2020/08/24   jidi           modify request status update
// 00003       2021/07/07   tangxiaoyu     add method
// -------------------------------------------------------------------

package wsio

import (
	"encoding/json"
	"github.com/wengoldx/wing/logger"
)

const (
	// StSuccess success status
	SioSuccess = 1 + iota

	// StError error status
	SioError
)

// Socket signlaing event function
type SocketEvent func(c *SocketController, uuid, params string) string

// SocketController signaling controller
type SocketController struct {
	evt     string
	handler SocketEvent
}

// EventAck socket event ack
type EventAck struct {
	State   int    `json:"state"`
	Message string `json:"message"`
}

// SocketHandler socket handler callbacks
type SocketHandler struct {
	Authenticate func(token string) (string, error)
	Connect      func(uuid string) error
	Disconnect   func()
}

// AckResp response normal ack to socket client
func (c *SocketController) AckResp(msg string) string {
	resp, _ := json.Marshal(&EventAck{
		State: SioSuccess, Message: msg,
	})
	logger.I("Ack evt[", c.evt, "] resp: ", msg)
	return string(resp)
}

// AckError response error ack to socket client
func (c *SocketController) AckError(msg string) string {
	resp, _ := json.Marshal(&EventAck{
		State: SioError, Message: msg,
	})
	logger.E("Ack  evt[", c.evt, "] err: ", msg)
	return string(resp)
}
