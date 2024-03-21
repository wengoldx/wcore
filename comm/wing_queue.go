// Copyright (c) 2018-2028 Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2019/05/22   yangping       New version
// -------------------------------------------------------------------

package comm

import (
	"container/list"
	"fmt"
	"sync"

	"github.com/wengoldx/wing/invar"
)

// Queue the type of queue with sync lock
//
//				--------- <- Head
//	Quere Top : |   1   | -> Pop
//				+-------+
//				|   2   |
//				+-------+
//				|  ...  |
//				+-------+
//		Push -> |   n   | : Queue Back (or Bottom)
//				+-------+
type Queue struct {
	list  *list.List
	mutex sync.Mutex
}

// GenQueue generat a new queue instance
func GenQueue() *Queue {
	return &Queue{list: list.New()}
}

// Push push a data to queue back if the data not nil
func (q *Queue) Push(data any) {
	if data == nil {
		return
	}
	q.mutex.Lock()
	defer q.mutex.Unlock()
	q.list.PushBack(data)
}

// Pop pick and remove the front data of queue,
// it will return invar.ErrEmptyData error if the queue is empty
func (q *Queue) Pop() (any, error) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if e := q.list.Front(); e != nil {
		q.list.Remove(e)
		return e.Value, nil
	}
	return nil, invar.ErrEmptyData
}

// Pick pick but not remove the front data of queue,
// it will return invar.ErrEmptyData error if the queue is empty
func (q *Queue) Pick() (any, error) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if e := q.list.Front(); e != nil {
		return e.Value, nil
	}
	return nil, invar.ErrEmptyData
}

// Clear clear the queue all data
func (q *Queue) Clear() {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	for e := q.list.Front(); e != nil; {
		en := e.Next()
		q.list.Remove(e)
		e = en
	}
}

// Len return the length of queue
func (q *Queue) Len() int {
	return q.list.Len()
}

// Fetch queue nodes, use callback returns to remove node or interupt fetch.
// Notice that DO NOT do heavy performence codes in callback, exist a lock here!
//
// For example caller code like (Events is a Queue object):
//
//	// Try fetch waiting requests to send to idle clusters
//	h.Events.Fetch(func(request any) (bool, bool) {
//		if clusterid := h.getIdleCluster(); clusterid != "" {
//			go h.sendRequest(clusterid, request)
//			return true /* remove */, false /* continue fetch */
//		}
//
//		// Remain request node and interupt fetch
//		return false, true
//	})
func (q *Queue) Fetch(callback func(value any) (bool, bool)) {
	if callback != nil {
		q.mutex.Lock()
		defer q.mutex.Unlock()

		for e := q.list.Front(); e != nil; e = e.Next() {
			remove, interupt := callback(e.Value)
			if remove { // Delete or remain the node
				q.list.Remove(e)
			}
			if interupt { // Interupt fetch or continue
				return
			}
		}
	}
}

// Dump print out the queue data.
// this method maybe just use for debug to out put queue items
func (q *Queue) Dump() {
	fmt.Println("-- dump the queue: (front -> back)")
	for e := q.list.Front(); e != nil; e = e.Next() {
		logs := fmt.Sprintf("   : %v", e.Value)
		fmt.Println(logs)
	}
}
