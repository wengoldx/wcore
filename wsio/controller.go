// Copyright (c) 2019-2029 DY All Rights Reserved.
//
// Author : yangping
// Email  : youhei_yp@163.com
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2022/02/09   yangping       New version
// -------------------------------------------------------------------

package wsio

import (
	"encoding/json"
	sio "github.com/googollee/go-socket.io"
	"github.com/wengoldx/wing/logger"
)

// Auth client outset, it will disconnect when return no-nil error
//	@param token client login jwt-token contain uuid or optional data in claims key string
//	@return - string client uuid
//			- interface{} client optional data parsed from token
//			- error Exception message
type AuthHandler func(token string) (string, string, error)

// Client connected callback, it will disconnect when return no-nil error
//	@param uuid client unique id
//	@param option client login optional data, maybe nil
//	@return - error Exception message
type ConnectHandler func(uuid, option string) error

// Client disconnected handler function
//	@param uuid client unique id
//	@param option client login optional data, maybe nil
//
// `NOTICE` :
//
// The client of uuid already released when call this event function.
type DisconnectHandler func(uuid, option string)

// Socket signlaing event function
type SignalingEvent func(sc sio.Socket, uuid, params string) string

// Socket signaling adaptor to register events
type SignalingAdaptor interface {

	// Retruen socket signaling events
	Signalings() []string

	// Dispath socket signaling callback by event
	Dispatch(evt string) SignalingEvent
}

// Socket event ack
type EventAck struct {
	State   int    `json:"state"`
	Message string `json:"message"`
}

const (
	// StSuccess success status
	StSuccess = iota + 1

	// StError error status
	StError
)

// Response normal ack to socket client
func AckResp(msg string) string {
	resp, _ := json.Marshal(&EventAck{
		State: StSuccess, Message: msg,
	})
	logger.D("SIO Response data >>", msg)
	return string(resp)
}

// Response error ack to socket client
func AckError(msg string) string {
	resp, _ := json.Marshal(&EventAck{
		State: StError, Message: msg,
	})
	logger.E("SIO Response err >>", msg)
	return string(resp)
}
