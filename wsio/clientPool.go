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
	"errors"
	sio "github.com/googollee/go-socket.io"
	"github.com/wengoldx/wing/invar"
	"github.com/wengoldx/wing/logger"
	"sort"
	"sync"
)

// ClientPool client pool
type ClientPool struct {
	lock    sync.Mutex         // Mutex sync lock
	clients map[string]*client // Client map
	s2c     map[string]string  // Socket id to Client id
	wait    map[string]int     // Wait map
}

// clientPool singletone instance
var clientPool *ClientPool

func init() {
	clientPool = &ClientPool{
		clients: make(map[string]*client),
		s2c:     make(map[string]string),
		wait:    make(map[string]int),
	}
}

// getClient return client if exist.
func (cp *ClientPool) getClient(cid string) *client {
	return cp.clients[cid]
}

// GetClientPool return ClientPool object
func GetClientPool() *ClientPool {
	return clientPool
}

// Register bind socket to client.
func (cp *ClientPool) Register(cid string, sc sio.Socket) error {
	cp.lock.Lock()
	defer cp.lock.Unlock()

	if err := cp.registerLocked(cid, sc); err != nil {
		logger.E("Regisger client err:", err.Error())
		return err
	}
	return nil
}

// Deregister unbind socket to client.
func (cp *ClientPool) Deregister(sc sio.Socket) {
	cp.lock.Lock()
	defer cp.lock.Unlock()

	cp.deregisterLocked(sc)
}

// IsExist check client if exist.
func (cp *ClientPool) IsExist(cid string) bool {
	_, ok := cp.clients[cid]
	return ok
}

// Sid2Cid return client uuid from socket id.
func (cp *ClientPool) Sid2Cid(sid string) string {
	cp.lock.Lock()
	defer cp.lock.Unlock()

	return cp.s2c[sid]
}

// Signaling send signaling to indicate client.
func (cp *ClientPool) Signaling(cid, data string, evt invar.Event) error {
	if c, ok := cp.clients[cid]; ok {
		return c.signaling(evt, data)
	}
	return invar.ErrTagOffline
}

//Waiting client enter wait queue
func (cp *ClientPool) Waiting(cid string) {
	cp.lock.Lock()
	defer cp.lock.Unlock()

	cp.waitLocked(cid)
}

//LeaveWait client leave wait queue
func (cp *ClientPool) LeaveWait(cid string) {
	cp.lock.Lock()
	defer cp.lock.Unlock()

	cp.leaveWaitLock(cid)

}

// registerLocked register the client without acquiring the lock.
func (cp *ClientPool) registerLocked(cid string, sc sio.Socket) error {
	var newOne *client
	sid := sc.Id()

	if oldOne, ok := cp.clients[cid]; ok {
		oldOneID := oldOne.socket.Id()
		if oldOneID == sid {
			logger.W("Client", cid, "already bind socket", sid)
			return nil
		}
		logger.I("Drop socket", oldOneID, "of client", cid)
		delete(cp.s2c, oldOneID)
		oldOne.deregister()
		newOne = oldOne
	} else {
		newOne = newClient(cid)
	}

	// bind client with socket
	if err := newOne.register(sc); err != nil {
		return err
	}

	logger.I("Add client", cid, "bind socket", sid)
	cp.clients[cid] = newOne
	cp.s2c[sid] = cid // same as uuid
	return nil
}

// deregisterLocked unregister the client without acquiring the lock.
func (cp *ClientPool) deregisterLocked(sc sio.Socket) {
	sid := sc.Id()
	if cid := cp.s2c[sid]; cid != "" {
		delete(cp.s2c, sid)
		if c := cp.clients[cid]; c != nil {
			logger.I("Remove client", cid)
			// service.DelClient(models.Clients, cid)
			delete(cp.clients, cid)
			if _, ok := cp.wait[cid]; ok {
				delete(cp.wait, cid)
			}
			c.deregister()
			return
		}
		logger.I("Unbind socket", sid, "of client", cid)
	}
	sc.Disconnect()
}

//waitLocked waiting clients enter the waiting queue
func (cp *ClientPool) waitLocked(cid string) {
	if _, ok := cp.wait[cid]; ok {
		cp.wait[cid] = cp.wait[cid] + 1
		return
	}
	cp.wait[cid] = 1
	logger.I("Add client", cid, "enter wait")

}

//leaveWaitLock clients no longer waiting leave the waiting queue
func (cp *ClientPool) leaveWaitLock(cid string) {
	if _, ok := cp.wait[cid]; ok {
		logger.I("client", cid, "leave wait")
		delete(cp.wait, cid)
		return
	}
	logger.I("client", cid, "is not exist wait")
}

// GetWaitMax return wait map
func (cp *ClientPool) GetWaitMax() ([]string, error) {
	type uc struct {
		uuid  string
		times int
	}
	var max []uc
	if len(cp.wait) == 0 {
		return nil, errors.New("no wait client")
	}
	for k, v := range cp.wait {
		max = append(max, uc{uuid: k, times: v})
	}
	sort.Slice(max, func(i, j int) bool {
		return max[i].times > max[j].times
	})
	// add each client's uuid to string array
	var ids []string
	for _, v := range max {
		ids = append(ids, v.uuid)
	}
	return ids, nil
}
