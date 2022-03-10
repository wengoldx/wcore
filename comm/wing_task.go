// Copyright (c) 2019-2029 DY All Rights Reserved.
//
// Author : yangping
// Email  : youhei_yp@163.com
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2019/05/22   yangping       New version
// -------------------------------------------------------------------

package comm

import (
	"github.com/wengoldx/wing/logger"
	"time"
)

/* ----------------------------------------------------------------- */
/* WARNING :                                                         */
/* The Task functions should restruct as TaskPool and new type Task, */
/* please do not use the old Task functions.                         */
/* ----------------------------------------------------------------- */

// Task the type of task, it support execute mutilpe jobs with job datas.
type Task struct {
	queue     *Queue        // jobs data queue of task
	interrupt bool          // enable interrupt task
	interval  time.Duration // sleep interval betwwen tow task job
	executing bool          // task executing status

}

// Notice that all generated tasks will using this chan
var chexe = make(chan string) // tasks chan

// TaskCallback task callback
type TaskCallback func(data interface{}) error

// GenTask generat a new task instance, you can set the interval duration
// and interrupt flag as the follow format:
//
// ---
//
//	interrupt := 1  // interrupt to execut the remain task jobs when case error
//	interval := 500 // sleep interval between task jobs in millisecond
//	task := comm.GenTask(callback, interrupt, interval)
//	task.Post(jobdata)
func GenTask(callback TaskCallback, options ...int) *Task {
	// generat the task and fill default options
	task := &Task{
		queue: GenQueue(), interrupt: false, interval: 0, executing: false,
	}

	// set task options from given data
	if optlen := len(options); optlen > 0 {
		task.interrupt = options[0] > 0
		if optlen > 1 {
			task.interval = time.Duration(options[1]) * time.Millisecond
		}
	}

	// start task channel to listen
	go task.innerTaskExecuter(callback)
	logger.I("Generat task:{interrupt:", task.interrupt, ", interval:", task.interval, "}")
	return task
}

// Post post a job to queue back
func (t *Task) Post(jobdata interface{}) {
	if jobdata == nil {
		logger.E("Invalid job data, abort post")
		return
	}
	t.queue.Push(jobdata)
	t.innerPostFor("Post Action")
}

// SetInterrupt set interrupt flag
func (t *Task) SetInterrupt(interrupt bool) {
	t.interrupt = interrupt
}

// setInterval set wait interval between tasks in microseconds, and it must > 0.
func (t *Task) SetInterval(interval int) {
	t.interval = time.Duration(interval) * time.Millisecond
}

// innerPostFor start runtime to post action
func (t *Task) innerPostFor(action string) {
	logger.I("Start runtime to post action:", action)
	go func() { chexe <- action }()
}

// innerTaskExecuter task execute monitor to listen tasks
func (t *Task) innerTaskExecuter(callback TaskCallback) {
	for {
		select {
		case action := <-chexe:
			logger.I("Received request from:", action)
			if callback == nil {
				logger.E("Nil task callback, abort request")
				break
			}

			// check current if executing status
			if t.executing {
				logger.W("Bussying now, try the next time...")
				break
			}

			// flag on executing and popup the topmost task to execte
			t.executing = true
			taskdata, err := t.queue.Pop()
			if err != nil {
				t.executing = false
				logger.I("Executed all tasks")
				break
			}

			if err := callback(taskdata); err != nil {
				logger.E("Execute task callback err:", err)
				if t.interrupt {
					logger.I("Interrupted tasks when case error")
					t.executing = false
					break
				}
			}
			if t.interval > 0 {
				logger.I("Waiting to next task after:", t.interval)
				time.Sleep(t.interval)
			}
			t.executing = false
			t.innerPostFor("Next Action")
		}
	}
}
