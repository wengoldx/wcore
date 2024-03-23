// Copyright (c) 2018-Now Dunyu All Rights Reserved.
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
	"strings"

	"github.com/wengoldx/wing/logger"
)

// WAuthController the extend controller base on WingController to support
// auth account from http headers, the client caller must append two headers
// before post request if expect the controller method enable execute token
// authentication from header.
//
// * Authoration : It must fixed keyword as WENGOLD-V1.1
//
// * Token : Authenticate JWT token responsed by login success
//
// * Location : Optional value of client indicator, global location
//
// `USAGE` :
//
// The validator register code of input params struct see WingController description,
// but the restful auth api of router method as follow usecase 1 and 2.
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
//	//	@Param Authoration header string true "WENGOLD-V1.1"
//	//	@Param Token       header string true "Authentication token"
//	//	@Param data body types.Accout true "input param description"
//	//	@Success 200 {string} "response data description"
//	//	@router /login [post]
//	func (c *AccController) AccLogin() {
//		ps := &types.Accout{}
//		c.DoAfterValidated(ps, func(uuid string) (int, any) {
//		// Or get authed account password as :
//		// c.DoAfterAuthValidated(ps, func(uuid, pwd string) (int, any) {
//			// do same business with input NO-EMPTY account uuid,
//			// directe use c and ps param in this methed.
//			// ...
//			return http.StatusOK, "Done business"
//		} , false /* not limit error message even code is 40x */)
//	}
//
// `USECASE 2. Auth account on GET http method`
//
//	//	@Description Restful api bind with /detail on GET method
//	//	@Param Authoration header string true "WENGOLD-V1.1"
//	//	@Param Token       header string true "Authentication token"
//	//	@Success 200 {types.Detail} "response data description"
//	//	@router /detail [get]
//	func (c *AccController) AccDetail() {
//		if uuid := c.AuthRequestHeader(); uuid != "" {
//			// use c.BindValue("fieldkey", out) parse params from url
//			c.ResponJSON(service.AccDetail(uuid))
//		}
//	}
//
// `USECASE 3. No-Auth and Use WingController`
//
//	//	@Description Restful api bind with /update on POST method
//	//	@Param data body types.UserInfo true "input param description"
//	//	@Success 200
//	//	@router /update [post]
//	func (c *AccController) AccUpdate() {
//		ps := &types.UserInfo{}
//		c.WingController.DoAfterValidated(ps, func() (int, any) {
//			// directe use c and ps param in this methed.
//			// ...
//			return http.StatusOK, nil
//		} , false /* not limit error message even code is 40x */)
//	}
//
// `USECASE 4. No-Auth and Custom code`
//
//	//	@Description Restful api bind with /list on GET method
//	//	@Success 200 {object} []types.Account "response data description"
//	//	@router /list [get]
//	func (c *AccController) AccList() {
//		// do same business without auth and input params
//		c.ResponJSON(service.AccList())
//	}
type WAuthController struct {
	WingController
}

// NextFunc2 do action after input params validated, it decode token to get account uuid.
type NextFunc2 func(uuid string) (int, any)

// NextFunc3 do action after input params validated, it decode token to get account uuid and password.
type NextFunc3 func(uuid, pwd string) (int, any)

// AuthHandlerFunc auth request token from http header and returen account secures.
type AuthHandlerFunc func(token string) (string, string)

// RoleHandlerFunc verify role access permission from account service.
type RoleHandlerFunc func(sub, obj, act string) bool

// Global handler function to auth token from http header
var GAuthHandlerFunc AuthHandlerFunc

// Global handler function to verify role from http header
var GRoleHandlerFunc RoleHandlerFunc

// Get authoration and token from http header, than verify it and return account secures.
func (c *WAuthController) AuthRequestHeader() string {
	uuid, _ := c.innerAuthHeader()
	return uuid
}

// DoAfterValidated do bussiness action after success validate the given json data.
func (c *WAuthController) DoAfterValidated(ps any, nextFunc2 NextFunc2, option ...any) {
	if uuid, _ := c.innerAuthHeader(); uuid != "" {
		isprotect := !(len(option) > 0 && !option[0].(bool))
		c.doAfterValidatedInner("json", ps, nextFunc2, uuid, true, isprotect)
	}
}

// DoAfterUnmarshal do bussiness action after success unmarshaled the given json data.
func (c *WAuthController) DoAfterUnmarshal(ps any, nextFunc2 NextFunc2, option ...any) {
	if uuid, _ := c.innerAuthHeader(); uuid != "" {
		isprotect := !(len(option) > 0 && !option[0].(bool))
		c.doAfterValidatedInner("json", ps, nextFunc2, uuid, false, isprotect)
	}
}

// DoAfterValidatedXml do bussiness action after success validate the given xml data.
func (c *WAuthController) DoAfterValidatedXml(ps any, nextFunc2 NextFunc2, option ...any) {
	if uuid, _ := c.innerAuthHeader(); uuid != "" {
		isprotect := !(len(option) > 0 && !option[0].(bool))
		c.doAfterValidatedInner("xml", ps, nextFunc2, uuid, true, isprotect)
	}
}

// DoAfterUnmarshalXml do bussiness action after success unmarshaled the given xml data.
func (c *WAuthController) DoAfterUnmarshalXml(ps any, nextFunc2 NextFunc2, option ...any) {
	if uuid, _ := c.innerAuthHeader(); uuid != "" {
		isprotect := !(len(option) > 0 && !option[0].(bool))
		c.doAfterValidatedInner("xml", ps, nextFunc2, uuid, false, isprotect)
	}
}

// ------------------------------------------------------

// DoAfterAuthValidated do bussiness action after success validate the given json data.
func (c *WAuthController) DoAfterAuthValidated(ps any, nextFunc3 NextFunc3, option ...any) {
	if uuid, pwd := c.innerAuthHeader(); uuid != "" {
		isprotect := !(len(option) > 0 && !option[0].(bool))
		c.doAfterValidatedInner3("json", ps, nextFunc3, uuid, pwd, true, isprotect)
	}
}

// DoAfterAuthUnmarshal do bussiness action after success unmarshaled the given json data.
func (c *WAuthController) DoAfterAuthUnmarshal(ps any, nextFunc3 NextFunc3, option ...any) {
	if uuid, pwd := c.innerAuthHeader(); uuid != "" {
		isprotect := !(len(option) > 0 && !option[0].(bool))
		c.doAfterValidatedInner3("json", ps, nextFunc3, uuid, pwd, false, isprotect)
	}
}

// DoAfterAuthValidatedXml do bussiness action after success validate the given xml data.
func (c *WAuthController) DoAfterAuthValidatedXml(ps any, nextFunc3 NextFunc3, option ...any) {
	if uuid, pwd := c.innerAuthHeader(); uuid != "" {
		isprotect := !(len(option) > 0 && !option[0].(bool))
		c.doAfterValidatedInner3("xml", ps, nextFunc3, uuid, pwd, true, isprotect)
	}
}

// DoAfterAuthUnmarshalXml do bussiness action after success unmarshaled the given xml data.
func (c *WAuthController) DoAfterAuthUnmarshalXml(ps any, nextFunc3 NextFunc3, option ...any) {
	if uuid, pwd := c.innerAuthHeader(); uuid != "" {
		isprotect := !(len(option) > 0 && !option[0].(bool))
		c.doAfterValidatedInner3("xml", ps, nextFunc3, uuid, pwd, false, isprotect)
	}
}

// ------------------------------------------------------

// Get authoration and token from http header, than verify it and return account secures.
func (c *WAuthController) innerAuthHeader() (string, string) {
	if GAuthHandlerFunc == nil || GRoleHandlerFunc == nil {
		c.E405Disabled("Controller not set global handlers!")
		return "", ""
	}

	// check authoration secure key
	authoration := strings.ToUpper(c.Ctx.Request.Header.Get("Authoration"))
	if authoration != "WENGOLD-V1.1" {
		if strings.HasPrefix(authoration, "WENGOLD") {
			c.E426UpgradeRequired("Upgrade required to WENGOLD-V1.1, not " + authoration)
			return "", ""
		}

		c.E401Unauthed("Unsupport authoration: " + authoration)
		return "", ""
	}

	// get token from header and verify it and user role
	if token := c.Ctx.Request.Header.Get("Token"); token != "" {
		if uuid, pwd := GAuthHandlerFunc(token); uuid == "" {
			c.E401Unauthed("Unauthed header token!")
			return "", ""
		} else {
			if !GRoleHandlerFunc(uuid, c.Ctx.Input.URL(), c.Ctx.Request.Method) {
				c.E403Denind("Role permission denied for " + uuid)
				return "", ""
			}

			if !c.HideRespLogs {
				logger.D("Authenticated account:", uuid)
			}
			return uuid, pwd
		}
	}

	// token is empty or invalid, response unauthed
	c.E401Unauthed("Unauthed header token!")
	return "", ""
}

// doAfterValidatedInner do bussiness action after success unmarshal params or
// validate the unmarshaled json data.
func (c *WAuthController) doAfterValidatedInner(datatype string,
	ps any, nextFunc2 NextFunc2, uuid string, isvalidate, isprotect bool) {
	if !c.validatrParams(datatype, ps, isvalidate) {
		return
	}

	// execute business function after unmarshal and validated
	if status, resp := nextFunc2(uuid); resp != nil {
		c.responCheckState(datatype, isprotect, status, resp)
	} else {
		c.responCheckState(datatype, isprotect, status)
	}
}

// doAfterValidatedInner3 do bussiness action after success unmarshal params or
// validate the unmarshaled json data.
func (c *WAuthController) doAfterValidatedInner3(datatype string,
	ps any, nextFunc3 NextFunc3, uuid, pwd string, isvalidate, isprotect bool) {
	if !c.validatrParams(datatype, ps, isvalidate) {
		return
	}

	// execute business function after unmarshal and validated
	if status, resp := nextFunc3(uuid, pwd); resp != nil {
		c.responCheckState(datatype, isprotect, status, resp)
	} else {
		c.responCheckState(datatype, isprotect, status)
	}
}
