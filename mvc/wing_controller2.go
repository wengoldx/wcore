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
// * Authoration : It must fixed keyword as WENGOLD
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
//	//	@Description Restful api bind with /login on POST method
//	//	@Param Authoration header string true "WENGOLD"
//	//	@Param Token       header string true "Authentication token"
//	//	@Param data body types.Accout true "input param description"
//	//	@Success 200 {string} "response data description"
//	//	@router /login [post]
//	func (c *AccController) AccLogin() {
//		ps := &types.Accout{}
//		c.DoAfterValidated(ps, func(uuid string) (int, interface{}) {
//			// do same business with input no-empty account uuid
//			// directe use c and ps param in this methed.
//			// ...
//			return http.StatusOK, "Done business"
//		} , false /* not limit error message even code is 40x */)
//	}
//
// `OR auth account on GET http method`
//
//	//	@Description Restful api bind with /info on GET method
//	//	@Param Authoration header string true "WENGOLD"
//	//	@Param Token       header string true "Authentication token"
//	//	@Success 200 {types.AccInfo} "response data description"
//	//	@router /info [get]
//	func (c *AccController) AccInfo() {
//		if uuid := c.AuthRequestHeader(); uuid != "" {
//			// use c.BindValue("fieldkey", out) parse params from url
//			return service.AccInfo()
//		}
//	}
type WAuthController struct {
	WingController
}

// AuthFunc auth request token from http header and returen account uuid.
type AuthFunc func(token string) string

// Global handler function to auth token from http header
var GAuthHandlerFunc AuthFunc

func (c *WAuthController) AuthRequestHeader() string {
	if GAuthHandlerFunc == nil {
		c.E403Denind("Controller not set global auth hander!")
		return ""
	}

	// check authoration secure key
	authoration := c.Ctx.Request.Header.Get("Authoration")
	if strings.ToUpper(authoration) != "WENGOLD" {
		c.E401Unauthed("Invalid header authoration: " + authoration)
		return ""
	}

	// get token from header and verify it
	if token := c.Ctx.Request.Header.Get("Token"); token != "" {
		if uuid := GAuthHandlerFunc(token); uuid != "" {
			logger.D("Authenticated account:", uuid)
			return uuid
		}
	}

	// token is empty or invalid, response unauthed
	c.E401Unauthed("Unauthed header token!")
	return ""
}

// DoAfterValidated do bussiness action after success validate the given json data,
// notice that you should register the field level validator for the input data's struct,
// then use it in struct describetion label as validate target.
//	see WAuthController
func (c *WAuthController) DoAfterValidated(ps interface{}, nextFunc2 NextFunc2, option ...interface{}) {
	if uuid := c.AuthRequestHeader(); uuid != "" {
		isprotect := !(len(option) > 0 && !option[0].(bool))
		c.doAfterParsedOrValidated("json", ps, nextFunc2, uuid, true, isprotect)
	}
}

// DoAfterUnmarshal do bussiness action after success unmarshaled the given json data.
//	see DoAfterValidated
func (c *WAuthController) DoAfterUnmarshal(ps interface{}, nextFunc2 NextFunc2, option ...interface{}) {
	if uuid := c.AuthRequestHeader(); uuid != "" {
		isprotect := !(len(option) > 0 && !option[0].(bool))
		c.doAfterParsedOrValidated("json", ps, nextFunc2, uuid, false, isprotect)
	}
}

// DoAfterValidatedXml do bussiness action after success validate the given xml data.
//	see DoAfterValidated
func (c *WAuthController) DoAfterValidatedXml(ps interface{}, nextFunc2 NextFunc2, option ...interface{}) {
	if uuid := c.AuthRequestHeader(); uuid != "" {
		isprotect := !(len(option) > 0 && !option[0].(bool))
		c.doAfterParsedOrValidated("xml", ps, nextFunc2, uuid, true, isprotect)
	}
}

// DoAfterUnmarshalXml do bussiness action after success unmarshaled the given xml data.
//	see DoAfterValidated, DoAfterValidatedXml
func (c *WAuthController) DoAfterUnmarshalXml(ps interface{}, nextFunc2 NextFunc2, option ...interface{}) {
	if uuid := c.AuthRequestHeader(); uuid != "" {
		isprotect := !(len(option) > 0 && !option[0].(bool))
		c.doAfterParsedOrValidated("xml", ps, nextFunc2, uuid, false, isprotect)
	}
}
