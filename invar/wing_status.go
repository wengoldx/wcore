// Copyright (c) 2018-2019 WING All Rights Reserved.
//
// Author : yangping
// Email  : youhei_yp@163.com
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2019/05/22   yangping       New version
// 00002       2019/06/30   zhaixing       Add function from godfs
// -------------------------------------------------------------------
package invar

import (
	"net/http"
)

const (
	StatusOK             = http.StatusOK
	E400ParseParams      = http.StatusBadRequest
	E401Unauthorized     = http.StatusUnauthorized
	E403PermissionDenied = http.StatusForbidden
	E404Exception        = http.StatusNotFound
	E405FuncDisabled     = http.StatusMethodNotAllowed
	E406InputParams      = http.StatusNotAcceptable
	E408Timeout          = http.StatusRequestTimeout
	E409Duplicate        = http.StatusConflict
	E410Gone             = http.StatusGone
	E412InvalidState     = http.StatusPreconditionFailed
	E424Locked           = http.StatusLocked

	// @deprecated: use the up define of ExxxXXXX
	StatusErrParseParams      = E400ParseParams
	StatusErrInputParams      = E406InputParams
	StatusErrUnauthorized     = E401Unauthorized
	StatusErrCaseException    = E404Exception
	StatusErrPermissionDenind = E403PermissionDenied
	StatusErrFuncDisabled     = E405FuncDisabled
	StatusErrTimeout          = E408Timeout
	StatusErrDuplicate        = E409Duplicate
	StatusErrInvalidState     = E412InvalidState
	StatusErrLocked           = E424Locked
	StatusErrGone             = E410Gone
)

var statusText = map[int]string{
	StatusOK:             "OK",
	E400ParseParams:      "Parse Input Params Error",
	E401Unauthorized:     "Unauthorized",
	E403PermissionDenied: "Permission Denied",
	E404Exception:        "Case Exception",
	E405FuncDisabled:     "Function Disabled",
	E406InputParams:      "Invalid Input Params",
	E408Timeout:          "Request Timeout",
	E409Duplicate:        "Duplicate Request",
	E410Gone:             "Gone",
	E412InvalidState:     "Invalid State",
	E424Locked:           "Resource Locked",
}

// StatusText returns a text for the HTTP status code. It returns the empty
// string if the code is unknown.
func StatusText(code int) string {
	codetext := statusText[code]
	if codetext == "" {
		return http.StatusText(code)
	}
	return codetext
}
