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
	sio "github.com/googollee/go-socket.io"
	"github.com/wengoldx/wing/invar"
	"github.com/wengoldx/wing/logger"
)

// client socket connected client
type client struct {
	id     string     // client id same as uuid
	socket sio.Socket // access socket.io connection.
}

// newClient create a new client
func newClient(cid string) *client {
	c := &client{
		id:     cid,
		socket: nil,
	}
	return c
}

// register binds the socket with client.
func (c *client) register(sc sio.Socket) error {
	if c.registered() {
		cid, sid := c.id, sc.Id()
		logger.E("Client", cid, "already bind socket", sid)
		return invar.ErrDupRegister
	}
	c.socket = sc
	return nil
}

// deregister closes the client's socket.
func (c *client) deregister() {
	if c.registered() {
		sid := c.socket.Id()
		logger.I("Unbind socket", sid, "of client", c.id)
		c.socket.Disconnect()
		c.socket = nil
	}
}

// registered check client register status.
func (c *client) registered() bool {
	return c.socket != nil
}

// signaling send signaling message to client.
func (c *client) signaling(evt invar.Event, msg string) error {
	if !c.registered() {
		return invar.ErrInvalidState
	}
	c.socket.Emit(string(evt), msg)
	logger.I("Signlaing", evt, "to client", c.id, "with msg", msg)
	return nil
}
