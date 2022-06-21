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
	"encoding/json"
	"encoding/xml"
	"github.com/astaxie/beego"
	"github.com/go-playground/validator/v10"
	"github.com/wengoldx/wing/invar"
	"github.com/wengoldx/wing/logger"
	"strings"
)

// WingController the base controller to support bee http functions.
type WingController struct {
	beego.Controller
}

// WAuthController the extend bee controller to support http request
// agent and auth token to verify.
//
// `USAGE` :
//
//	// You can define the custom controller like:
//	type CustomController struct {
//		mvc.WAuthController
//	}
//
//	// Rest4Method custom RESTful APIs
//	// ----------------------------------------------
//	// @Description The restfull api method bind with router /xx
//	// @Param Authoration header string true "WENGOLD"
//	// @Param Token       header string true "Authentication token"
//	// @Param data body types.InParams true "input params"
//	// @Success 200 {string} "response string data"
//	// @Failure 400 parse input param error
//	// @router /xx [post]
//	func(c * CustomController) Rest4Method() {
//		ps := &types.InParams{}
//		c.DoAfterValidated(ps, func() (int, interface{}) {
//			logger.I("Rest4Method: validated input params:", ps)
//			return services.FuncImplMethod(ps)
//		}, services.IsProtect())
//	}
//
//	// Http request agent and token auth handler
//	func (c *CustomController) AuthHeaderFunc(token string) bool {
//		// TODO. Auth the input token, and return result status
//		return true
//	}
type WAuthController struct {
	WingController
	WAuthInterface // http request agent and token auth interface
}

// WAuthInterface is an interface to auth http request agent and token handler.
type WAuthInterface interface {
	AuthHeaderFunc(token string) bool
}

// NextFunc do action after input params validated.
type NextFunc func() (int, interface{})

// Validator use for verify the input params on struct level
var Validator *validator.Validate

// ensureValidatorGenerated generat the validator instance if need
func ensureValidatorGenerated() {
	if Validator == nil {
		Validator = validator.New()
	}
}

// RegisterValidators register struct field validators from given map
func RegisterValidators(valmap map[string]validator.Func) {
	for tag, valfunc := range valmap {
		RegisterFieldValidator(tag, valfunc)
	}
}

// RegisterFieldValidator register validators on struct field level
func RegisterFieldValidator(tag string, valfunc validator.Func) {
	ensureValidatorGenerated()
	if err := Validator.RegisterValidation(tag, valfunc); err != nil {
		logger.E("Register validator:"+tag+", err:", err)
		return
	}
	logger.I("Registered validator:", tag)
}

// Check authenticate from request header
func (c *WAuthController) Prepare() {
	authoration := c.Ctx.Request.Header.Get("Authoration")
	if authoration != "WENGOLD" {
		logger.E("Invalid header authoration:", authoration)

		// TODO. comment out when use token auth
		// c.E401Unauthed("Invalid header authoration:" + authoration)
		return
	}

	token := c.Ctx.Request.Header.Get("Token")
	if token == "" || (c.WAuthInterface != nil && !c.WAuthInterface.AuthHeaderFunc(token)) {
		logger.E("Unauthed request token:", token)

		// TODO. comment out when use token auth
		// c.E401Unauthed("Unauthed request token!")
		return
	}

	logger.D("Authoration:", authoration, "token:", token)
}

// responCheckState check respon state and print out log, the datatype must
// range in ['json', 'jsonp', 'xml', 'yaml'], if outof range current controller
// just return blank string to close http connection.
func (c *WingController) responCheckState(datatype string, needCheck bool, state int, data ...interface{}) {
	if state != invar.StatusOK {
		if needCheck {
			c.ErrorState(state)
			return
		}

		errmsg := invar.StatusText(state)
		ctl, act := c.GetControllerAndAction()
		logger.E("Respone "+strings.ToUpper(datatype)+" error:", state, ">", ctl+"."+act, errmsg)
	} else {
		ctl, act := c.GetControllerAndAction()
		logger.I("Respone state:OK-"+strings.ToUpper(datatype)+" >", ctl+"."+act)
	}

	c.Ctx.Output.Status = state
	if len(data) > 0 {
		c.Data[datatype] = data[0]
	}

	switch datatype {
	case "json":
		c.ServeJSON()
	case "jsonp":
		c.ServeJSONP()
	case "xml":
		c.ServeXML()
	case "yaml":
		c.ServeYAML()
	default:
		// just return blank string to close http connection
		logger.W("Invalid response data tyep:" + datatype)
		c.Ctx.ResponseWriter.Write([]byte(""))
	}
}

// ResponJSON sends a json response to client,
// it may not send data if the state is not status ok
func (c *WingController) ResponJSON(state int, data ...interface{}) {
	c.responCheckState("json", true, state, data...)
}

// ResponJSONUncheck sends a json response to client witchout uncheck state code.
func (c *WingController) ResponJSONUncheck(state int, dataORerr ...interface{}) {
	c.responCheckState("json", false, state, dataORerr...)
}

// ResponJSONP sends a jsonp response to client,
// it may not send data if the state is not status ok
func (c *WingController) ResponJSONP(state int, data ...interface{}) {
	c.responCheckState("jsonp", true, state, data...)
}

// ResponJSONPUncheck sends a jsonp response to client witchout uncheck state code.
func (c *WingController) ResponJSONPUncheck(state int, dataORerr ...interface{}) {
	c.responCheckState("jsonp", false, state, dataORerr...)
}

// ResponXML sends xml response to client,
// it may not send data if the state is not status ok
func (c *WingController) ResponXML(state int, data ...interface{}) {
	c.responCheckState("xml", true, state, data...)
}

// ResponXMLUncheck sends xml response to client witchout uncheck state code.
func (c *WingController) ResponXMLUncheck(state int, dataORerr ...interface{}) {
	c.responCheckState("xml", false, state, dataORerr...)
}

// ResponYAML sends yaml response to client,
// it may not send data if the state is not status ok
func (c *WingController) ResponYAML(state int, data ...interface{}) {
	c.responCheckState("yaml", true, state, data...)
}

// ResponYAML sends yaml response to client witchout uncheck state code.
func (c *WingController) ResponYAMLUncheck(state int, dataORerr ...interface{}) {
	c.responCheckState("yaml", false, state, dataORerr...)
}

// ResponData sends YAML, XML OR JSON, depending on the value of the Accept header,
// it may not send data if the state is not status ok
func (c *WingController) ResponData(state int, data ...map[interface{}]interface{}) {
	if state != invar.StatusOK {
		c.ErrorState(state)
		return
	}

	ctl, act := c.GetControllerAndAction()
	logger.I("Respone state:OK-DATA >", ctl+"."+act)
	c.Ctx.Output.Status = state
	if len(data) > 0 {
		c.Data = data[0]
	}
	c.ServeFormatted()
}

// ResponOK sends a empty success response to client
func (c *WingController) ResponOK() {
	ctl, act := c.GetControllerAndAction()
	logger.I("Respone state:OK >", ctl+"."+act)

	w := c.Ctx.ResponseWriter
	w.WriteHeader(invar.StatusOK)
	// FIXME here maybe not set content type when response error
	// w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(""))
}

// ErrorState response error state to client
func (c *WingController) ErrorState(state int, err ...string) {
	ctl, act := c.GetControllerAndAction()
	errmsg := invar.StatusText(state)
	if len(err) > 0 {
		errmsg += ", " + err[0]
	}
	logger.E("Respone error:", state, ">", ctl+"."+act, errmsg)

	w := c.Ctx.ResponseWriter
	w.WriteHeader(state)
	// FIXME here maybe not set content type when response error
	// w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(""))
}

// E400Params response 400 invalid params error state to client
func (c *WingController) E400Params(err ...string) {
	c.ErrorState(invar.E400ParseParams, err...)
}

// E400rUnmarshal response 400 unmarshal params error state to client
func (c *WingController) E400Unmarshal(err ...string) {
	c.ErrorState(invar.E400ParseParams, err...)
}

// E400Validate response 400 invalid params error state to client, then print
// the params data and validate error
func (c *WingController) E400Validate(ps interface{}, err ...string) {
	logger.E("Invalid input params:", ps)
	c.ErrorState(invar.E400ParseParams, err...)
}

// E401Unauthed response 401 unauthenticated error state to client
func (c *WingController) E401Unauthed(err ...string) {
	c.ErrorState(invar.E401Unauthorized, err...)
}

// E403Denind response 403 permission denind error state to client
func (c *WingController) E403Denind(err ...string) {
	c.ErrorState(invar.E403PermissionDenied, err...)
}

// E404Exception response 404 not found error state to client
func (c *WingController) E404Exception(err ...string) {
	c.ErrorState(invar.E404Exception, err...)
}

// E405Disabled response 405 function disabled error state to client
func (c *WingController) E405Disabled(err ...string) {
	c.ErrorState(invar.E405FuncDisabled, err...)
}

// E406Input response 406 invalid inputs error state to client
func (c *WingController) E406Input(err ...string) {
	c.ErrorState(invar.E406InputParams, err...)
}

// E409Duplicate response 409 duplicate error state to client
func (c *WingController) E409Duplicate(err ...string) {
	c.ErrorState(invar.E409Duplicate, err...)
}

// E410Gone response 410 gone error state to client
func (c *WingController) E410Gone(err ...string) {
	c.ErrorState(invar.E410Gone, err...)
}

// Get and check valid token from header by key 'token', it may
// response error to client when token not found or empty.
func (c *WingController) ViaAuthToken() string {
	token := c.Ctx.Request.Header.Get("token")
	if token == "" {
		c.E401Unauthed("Not found token from header!")
	}
	return token
}

// ClientFrom return client ip from who requested
func (c *WingController) ClientFrom() string {
	return c.Ctx.Request.RemoteAddr
}

// BindValue bind value with key from url, the dest container must pointer
func (c *WingController) BindValue(key string, dest interface{}) error {
	if err := c.Ctx.Input.Bind(dest, key); err != nil {
		logger.E("Parse", key, "from url, err:", err)
		return invar.ErrInvalidData
	}
	return nil
}

// doAfterValidatedInner do bussiness action after success unmarshal params or
// validate the unmarshaled json data.
func (c *WingController) doAfterParsedOrValidated(datatype string, ps interface{},
	nextFunc NextFunc, isvalidate, isprotect bool) {

	// unmarshal the input params
	switch datatype {
	case "json":
		if err := json.Unmarshal(c.Ctx.Input.RequestBody, ps); err != nil {
			c.E400Unmarshal(err.Error())
			return
		}
	case "xml":
		if err := xml.Unmarshal(c.Ctx.Input.RequestBody, ps); err != nil {
			c.E400Unmarshal(err.Error())
			return
		}
	default: // current not support the jsonp and yaml parse
		c.E404Exception("Invalid data type:" + datatype)
		return
	}

	// validate input params if need
	if isvalidate {
		ensureValidatorGenerated()
		if err := Validator.Struct(ps); err != nil {
			c.E400Validate(ps, err.Error())
			return
		}
	}

	// execute business function after unmarshal and validated
	if status, resp := nextFunc(); resp != nil {
		c.responCheckState(datatype, isprotect, status, resp)
	} else {
		c.responCheckState(datatype, isprotect, status)
	}
}

// DoAfterValidated do bussiness action after success validate the given json data,
// notice that you should register the field level validator for the input data's struct,
// then use it in struct describetion label as validate target.
//
// ---
//
// `types.go`
//
//	type struct Accout {
//		Acc string `json:"acc" validate:"required,IsVaildUuid"`
//		PWD string `json:"pwd" validate:"required_without"`
//		Num int    `json:"num"`
//	}
//
//	// define custom validator on struct field level
//	func isVaildUuid(fl validator.FieldLevel) bool {
//		m, _ := regexp.Compile("^[0-9a-zA-Z]*$")
//		str := fl.Field().String()
//		return m.MatchString(str)
//	}
//
//	func init() {
//		mvc.RegisterFieldValidator("IsVaildUuid", isVaildUuid)
//	}
//
// ---
//
// `controller.go`
//
//	func (c *AccController) AccLogin() {
//		ps := &types.Accout{}
//		c.DoAfterValidated(ps, func() (int, interface{}) {
//			// do same business function
//			// directe use c and ps param in this methed.
//			// ...
//			return http.StatusOK, "Done business"
//		} /** , false /* not limit error message even code is 40x */ */)
//	}
func (c *WingController) DoAfterValidated(ps interface{}, nextFunc NextFunc, option ...interface{}) {
	isprotect := !(len(option) > 0 && !option[0].(bool))
	c.doAfterParsedOrValidated("json", ps, nextFunc, true, isprotect)
}

// DoAfterUnmarshal do bussiness action after success unmarshaled the given json data.
//	see DoAfterValidated
func (c *WingController) DoAfterUnmarshal(ps interface{}, nextFunc NextFunc, option ...interface{}) {
	isprotect := !(len(option) > 0 && !option[0].(bool))
	c.doAfterParsedOrValidated("json", ps, nextFunc, false, isprotect)
}

// DoAfterValidatedXml do bussiness action after success validate the given xml data.
//	see DoAfterValidated
func (c *WingController) DoAfterValidatedXml(ps interface{}, nextFunc NextFunc, option ...interface{}) {
	isprotect := !(len(option) > 0 && !option[0].(bool))
	c.doAfterParsedOrValidated("xml", ps, nextFunc, true, isprotect)
}

// DoAfterUnmarshalXml do bussiness action after success unmarshaled the given xml data.
//	see DoAfterValidated, DoAfterValidatedXml
func (c *WingController) DoAfterUnmarshalXml(ps interface{}, nextFunc NextFunc, option ...interface{}) {
	isprotect := !(len(option) > 0 && !option[0].(bool))
	c.doAfterParsedOrValidated("xml", ps, nextFunc, false, isprotect)
}
