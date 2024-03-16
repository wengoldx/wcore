// Copyright (c) 2018-2028 Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2019/05/22   yangping       New version
// -------------------------------------------------------------------

package logger

import (
	"runtime"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

const (
	logConfigLevel   = "logger::level"   // configs key of logger level
	logConfigMaxDays = "logger::maxdays" // configs key of logger max days

	// LevelDebug debug level of logger
	LevelDebug = "debug"

	// LevelInfo info level of logger
	LevelInfo = "info"

	// LevelWarn warn level of logger
	LevelWarn = "warn"

	// LevelError error level of logger
	LevelError = "error"
)

// init initialize app logger
//
// `NOTICE` : you must config logger params in /conf/app.config file as:
//
// ---
//
//	[logger]
//	level = "debug"
//	maxdays = "7"
//
// ---
//
// - the level values range in : [debug, info, warn, error], default is info.
//
// - maxdays is the max days to hold logs cache, default is 7 days.
func init() {
	config := readLoggerConfigs()
	beego.SetLogger(logs.AdapterFile, config)
	beego.SetLogFuncCall(true)
	logs.SetLogFuncCallDepth(5)
	logs.Async(3) // allow asynchronous

	// set application logger level
	switch beego.AppConfig.String(logConfigLevel) {
	case LevelDebug:
		beego.SetLevel(beego.LevelDebug)
	case LevelInfo:
		beego.SetLevel(beego.LevelInformational)
	case LevelWarn:
		beego.SetLevel(beego.LevelWarning)
	case LevelError:
		beego.SetLevel(beego.LevelError)
	}
}

// readLoggerConfigs get logger configs
func readLoggerConfigs() string {
	app := beego.BConfig.AppName
	if app == "" || app == "beego" {
		app = "wing"
	}

	maxdays := beego.AppConfig.String(logConfigMaxDays)
	if maxdays == "" {
		maxdays = "7"
	}
	return "{\"filename\":\"logs/" + app + ".log\", \"daily\":true, \"maxdays\":" + maxdays + "}"
}

// appendFuncName append runtime calling function name start log prefix, it format as :
// ------------------------------------------------------------------------------------
// 2023/05/31 10:56:36.609 [I] [code_file.go:89]  FuncName() Log output message string
// ------------------------------------------------------------------------------------
func appendFuncName(v ...interface{}) []interface{} {
	/* Fixed the call skipe on 1 to filter current function name */
	if pc, _, _, ok := runtime.Caller(2); ok {
		if funcptr := runtime.FuncForPC(pc); funcptr != nil {
			if funname := funcptr.Name(); funname != "" {
				fns := strings.SplitAfter(funname, ".")
				v = append([]interface{}{fns[len(fns)-1] + "()"}, v...)
			}
		}
	}
	return v
}

// SetOutputLogger close console logger on prod mode and only remain file logger.
func SetOutputLogger() {
	if beego.BConfig.RunMode != "dev" && GetLevel() != LevelDebug {
		beego.BeeLogger.DelLogger(logs.AdapterConsole)
	}
}

// GetLevel return current logger output level
func GetLevel() string {
	switch beego.BeeLogger.GetLevel() {
	case beego.LevelDebug:
		return LevelDebug
	case beego.LevelInformational:
		return LevelInfo
	case beego.LevelWarning:
		return LevelWarn
	case beego.LevelError:
		return LevelError
	default:
		return ""
	}
}

// EM logs a message at emergency level.
func EM(v ...interface{}) {
	beego.Emergency(appendFuncName(v...)...)
}

// AL logs a message at alert level.
func AL(v ...interface{}) {
	beego.Alert(appendFuncName(v...)...)
}

// CR logs a message at critical level.
func CR(v ...interface{}) {
	beego.Critical(appendFuncName(v...)...)
}

// E logs a message at error level.
func E(v ...interface{}) {
	beego.Error(appendFuncName(v...)...)
}

// W logs a message at warning level.
func W(v ...interface{}) {
	beego.Warning(appendFuncName(v...)...)
}

// N logs a message at notice level.
func N(v ...interface{}) {
	beego.Notice(appendFuncName(v...)...)
}

// I logs a message at info level.
func I(v ...interface{}) {
	beego.Informational(appendFuncName(v...)...)
}

// D logs a message at debug level.
func D(v ...interface{}) {
	beego.Debug(appendFuncName(v...)...)
}
