// Copyright (c) 2018-2028 Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2019/05/22   yangping       New version
// -------------------------------------------------------------------

package mvc

import (
	"github.com/wengoldx/wing/logger"
	"strings"
)

// WAuthController the extend controller base on WingController to support auth account
// from http headers, the client caller must append tow headers before post http request.
//
// * Authoration : It must fixed keyword as WENGOLD or WENGOLD-NOSECURE
//
// * Token : Authenticate JWT token responsed by login success
//
//
// `USAGE` :
//
// The validator register code of input params struct see WingController description,
// but the restful api router function like follow.
//
// ---
//
// `controller.go`
//
//	// define custom controller using header auth function
//	type AccController struct {
//		mvc.WAuthController
//	}
//
//	func init() {
//		mvc.GAuthHandlerFunc = func(token string) (string, string) {
//			// decode and verify token string, than return indecated
//			// account uuid and password.
//			return "account uuid", "account password"
//		}
//	}
//
// `USECASE 1. Auth account and Parse input params`
//
//	//	@Description Restful api bind with /login on POST method
//	//	@Param Authoration header string true "WENGOLD"
//	//	@Param Token       header string true "Authentication token"
//	//	@Param data body types.Accout true "input param description"
//	//	@Success 200 {string} "response data description"
//	//	@router /login [post]
//	func (c *AccController) AccLogin() {
//		ps := &types.Accout{}
//		c.DoAfterValidated(ps, func(uuid, pwd string) (int, interface{}) {
//			// do same business with input no-empty account uuid,
//			// directe use c and ps param in this methed.
//			// ...
//			return http.StatusOK, "Done business"
//		} , false /* not limit error message even code is 40x */)
//	}
//
// `USECASE 2. Auth account on GET http method`
//
//	//	@Description Restful api bind with /info on GET method
//	//	@Param Authoration header string true "WENGOLD"
//	//	@Param Token       header string true "Authentication token"
//	//	@Success 200 {types.AccInfo} "response data description"
//	//	@router /info [get]
//	func (c *AccController) AccInfo() {
//		if uuid, _ := c.AuthRequestHeader(); uuid != "" {
//			// use c.BindValue("fieldkey", out) parse params from url
//			return service.AccInfo()
//		}
//	}
//
// `USECASE 3. Not auth but Parse input params`
//
//	//	@Description Restful api bind with /login on POST method
//	//	@Param Authoration header string true "WENGOLD-NOSECURE"
//	//	@Param data body types.UserInfo true "input param description"
//	//	@Success 200 {string} "response data description"
//	//	@router /login [post]
//	func (c *AccController) AccLogin() {
//		ps := &types.UserInfo{}
//		c.DoAfterValidated(ps, func(uuid, pwd string) (int, interface{}) {
//			// do same business with input empty account and pwd,
//			// directe use c and ps param in this methed.
//			// ...
//			return http.StatusOK, "Done business"
//		} , false /* not limit error message even code is 40x */)
//	}
type WAuthController struct {
	WingController
}

// AuthFunc auth request token from http header and returen account secures.
type AuthFunc func(token string) (string, string)

// Global handler function to auth token from http header
var GAuthHandlerFunc AuthFunc

// Get authoration and token from http header, than verify it and return account secures.
func (c *WAuthController) AuthRequestHeader() (string, string) {
	if GAuthHandlerFunc == nil {
		c.E403Denind("Controller not set global auth hander!")
		return "", ""
	}

	// check authoration secure key
	authoration := strings.ToUpper(c.Ctx.Request.Header.Get("Authoration"))
	if authoration != "WENGOLD" && authoration != "WENGOLD-NOSECURE" {
		c.E401Unauthed("Invalid header authoration: " + authoration)
		return "", ""
	} else if authoration == "WENGOLD-NOSECURE" {
		// FIXME :
		// Here means that current controller router method
		// no-need to auth header token, just return empty infos
		// and pass auth.
		return authoration, ""
	}

	// get token from header and verify it
	if token := c.Ctx.Request.Header.Get("Token"); token != "" {
		if uuid, pwd := GAuthHandlerFunc(token); uuid != "" {
			logger.D("Authenticated account:", uuid)
			return uuid, pwd
		}
	}

	// token is empty or invalid, response unauthed
	c.E401Unauthed("Unauthed header token!")
	return "", ""
}

// DoAfterValidated do bussiness action after success validate the given json data.
func (c *WAuthController) DoAfterValidated(ps interface{}, nextFunc2 NextFunc2, option ...interface{}) {
	if uuid, pwd := c.AuthRequestHeader(); uuid != "" {
		isprotect := !(len(option) > 0 && !option[0].(bool))
		c.doAfterParsedOrValidated(nextFunc2, "json", ps, uuid, pwd, true, isprotect)
	}
}

// DoAfterUnmarshal do bussiness action after success unmarshaled the given json data.
func (c *WAuthController) DoAfterUnmarshal(ps interface{}, nextFunc2 NextFunc2, option ...interface{}) {
	if uuid, pwd := c.AuthRequestHeader(); uuid != "" {
		isprotect := !(len(option) > 0 && !option[0].(bool))
		c.doAfterParsedOrValidated(nextFunc2, "json", ps, uuid, pwd, false, isprotect)
	}
}

// DoAfterValidatedXml do bussiness action after success validate the given xml data.
func (c *WAuthController) DoAfterValidatedXml(ps interface{}, nextFunc2 NextFunc2, option ...interface{}) {
	if uuid, pwd := c.AuthRequestHeader(); uuid != "" {
		isprotect := !(len(option) > 0 && !option[0].(bool))
		c.doAfterParsedOrValidated(nextFunc2, "xml", ps, uuid, pwd, true, isprotect)
	}
}

// DoAfterUnmarshalXml do bussiness action after success unmarshaled the given xml data.
func (c *WAuthController) DoAfterUnmarshalXml(ps interface{}, nextFunc2 NextFunc2, option ...interface{}) {
	if uuid, pwd := c.AuthRequestHeader(); uuid != "" {
		isprotect := !(len(option) > 0 && !option[0].(bool))
		c.doAfterParsedOrValidated(nextFunc2, "xml", ps, uuid, pwd, false, isprotect)
	}
}
