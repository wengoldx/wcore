// Copyright (c) 2018-2028 Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2019/05/22   yangping       New version
// 00002       2022/03/26   yangping       Using toolbox.Task
// -------------------------------------------------------------------

package comm

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/toolbox"
	"github.com/wengoldx/wing/logger"
	"time"
)

// Task datas for multipe generate
type WTask struct {
	Name    string           // monitor task name
	Func    toolbox.TaskFunc // monitor task execute function
	Spec    string           // monitor task interval
	ForProd bool             // indicate the task only for prod mode, default no limit
}

// Add a single monitor task to list
func AddTask(tname, spec string, f toolbox.TaskFunc) {
	monitor := toolbox.NewTask(tname, spec, f)
	monitor.ErrLimit = 0

	logger.I("Create task:", tname, "and add to list")
	toolbox.AddTask(tname, monitor)
}

// Generate tasks and start them as monitors.
func StartTasks(monitors []*WTask) {
	for _, m := range monitors {
		if m.ForProd && beego.BConfig.RunMode != "prod" {
			logger.W("Filter out task:", m.Name, "on dev mode")
			continue
		}
		AddTask(m.Name, m.Spec, m.Func)
	}

	toolbox.StartTask()
	logger.I("Started all monitors")
}

// Return task if exist, or nil when unexist
func GetTask(tname string) *toolbox.Task {
	if tasker, ok := toolbox.AdminTaskList[tname]; ok {
		return tasker.(*toolbox.Task)
	}
	return nil
}

// Task the type of task
type TTask struct {
	queue     *Queue
	interrupt bool
	interval  time.Duration
	executing bool
}

var chexe = make(chan string)

// TaskCallback task callback
type TaskCallback func(data interface{}) error

// GenTask generat a new task instance, you can set the interval duration
// and interrupt flag as the follow format:
// [CODE:]
//   interrupt := 1  // interrupt to execut the remain tasks when case error
//   interval := 500 // sleep interval between tasks in microseconds
//   task := comm.GenTask(callback, interrupt, interval)
//   task.Post(taskdata)
// [CODE]
func GenTask(callback TaskCallback, configs ...int) *TTask {
	// generat the task and fill default configs
	task := &TTask{
		queue: GenQueue(), interrupt: false, interval: 0, executing: false,
	}

	// set task configs from given data
	if configs != nil {
		task.interrupt = len(configs) > 0 && configs[0] > 0
		if len(configs) > 1 && configs[1] > 0 {
			task.interval = time.Duration(configs[1] * 1000)
		}
	}

	// start task channel to listen
	go task.innerTaskExecuter(callback)
	logger.I("Generat a task:{interrupt:", task.interrupt, ", interval:", task.interval, "}")
	return task
}

// Post post a task to tasks queue back
func (t *TTask) Post(taskdata interface{}, check bool) {
	if check && t.queue.Len() > 5000 {
		logger.E("Task queue too busy now!")
		return
	}

	if taskdata == nil {
		logger.E("Invalid task data, abort post")
		return
	}

	t.queue.Push(taskdata)
	t.innerPostFor("Post Action")
}

// SetInterrupt set interrupt flag
func (t *TTask) SetInterrupt(interrupt bool) {
	t.interrupt = interrupt
}

// setInterval set wait interval between tasks in microseconds, and it must > 0.
func (t *TTask) SetInterval(interval int) {
	if interval > 0 {
		t.interval = time.Duration(interval * 1000)
	}
}

// innerPostFor start runtime to post action
func (t *TTask) innerPostFor(action string) {
	logger.I("Start runtime to post action:", action)
	go func() { chexe <- action }()
}

// innerTaskExecuter task execute monitor to listen tasks
func (t *TTask) innerTaskExecuter(callback TaskCallback) {
	for {
		logger.I("Blocking for task requir select...")
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
