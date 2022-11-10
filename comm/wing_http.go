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
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/wengoldx/wing/invar"
	"github.com/wengoldx/wing/logger"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
)

// SetRequest use for set http request before execute http.Client.Do,
// you can use this middle-ware to set auth as username and passord, and so on.
//	@param req Http requester
//	@return - bool If current request ignore TLS verify or not, false is verify by default.
//			- error Exception message
type SetRequest func(req *http.Request) (bool, error)

const (
	// ContentTypeJson json content type
	ContentTypeJson = "application/json;charset=UTF-8"

	// ContentTypeForm form content type
	ContentTypeForm = "application/x-www-form-urlencoded"
)

// Logout http request url when flag is on, and logger level lowwer than debug.
var DebugPrintRequest = true

// Logout http response datas when flag is on, and logger level lowwer than debug.
var DebugPrintResponse = true

// readResponse read response body after executed request, it should return
// invar.ErrInvalidState when response code is not http.StatusOK.
func readResponse(resp *http.Response) ([]byte, error) {
	if resp.StatusCode != http.StatusOK {
		logger.E("Failed http client, status:", resp.StatusCode)
		return nil, invar.ErrInvalidState
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.E("Failed read response, err:", err)
		return nil, err
	}

	if DebugPrintResponse {
		logger.D("Response:", string(body))
	}
	return body, nil
}

// unmarshalResponse unmarshal response body after execute request,
// it may not check the given body if empty.
func unmarshalResponse(body []byte, out interface{}) error {
	if err := json.Unmarshal(body, out); err != nil {
		logger.E("Unmarshal body to struct err:", err)
		return err
	}

	if DebugPrintResponse {
		logger.D("Response struct:", out)
	}
	return nil
}

// httpPostJson http post method, you can set post data as json struct.
func httpPostJson(tagurl string, postdata interface{}) ([]byte, error) {
	params, err := json.Marshal(postdata)
	if err != nil {
		logger.E("Marshal post data err:", err)
		return nil, err
	}

	resp, err := http.Post(tagurl, ContentTypeJson, bytes.NewReader(params))
	if err != nil {
		logger.E("Http post json err:", err)
		return nil, err
	}

	defer resp.Body.Close()
	return readResponse(resp)
}

// httpPostForm http post method, you can set post data as url.Values.
func httpPostForm(tagurl string, postdata url.Values) ([]byte, error) {
	resp, err := http.PostForm(tagurl, postdata)
	if err != nil {
		logger.E("Http post form err:", err)
		return nil, err
	}

	defer resp.Body.Close()
	return readResponse(resp)
}

// --------------------------------------------------

// GetIP get just ip not port from controller.Ctx.Request.RemoteAddr of beego
func GetIP(remoteaddr string) string {
	ip, _, _ := net.SplitHostPort(remoteaddr)
	if ip == "::1" {
		ip = "127.0.0.1"
	}
	logger.I("Got ip [", ip, "] from [", remoteaddr, "]")
	return ip
}

// Get all the loacl IP of current deploy machine
func GetLocalIPs() ([]string, error) {
	netfaces, err := net.Interfaces()
	if err != nil {
		logger.E("Get ip interfaces err:", err)
		return nil, err
	}

	ips := []string{}
	for _, netface := range netfaces {
		addrs, err := netface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				if v.IP.IsGlobalUnicast() {
					ips = append(ips, v.IP.String())
				}
			}
		}
	}

	// Check the result list is empty
	if len(ips) == 0 {
		return nil, invar.ErrNotFound
	}

	return ips, nil
}

// EncodeUrl encode url params
func EncodeUrl(rawurl string) string {
	enurl, err := url.Parse(rawurl)
	if err != nil {
		logger.E("Encode urlm err:", err)
		return rawurl
	}
	enurl.RawQuery = enurl.Query().Encode()
	return enurl.String()
}

// HttpGet handle http get method
func HttpGet(tagurl string, params ...interface{}) ([]byte, error) {
	if len(params) > 0 {
		tagurl = fmt.Sprintf(tagurl, params...)
	}

	rawurl := EncodeUrl(tagurl)
	if DebugPrintRequest {
		logger.D("Http Get:", rawurl)
	}

	resp, err := http.Get(rawurl)
	if err != nil {
		logger.E("Failed http get, err:", err)
		return nil, err
	}
	defer resp.Body.Close()
	return readResponse(resp)
}

// HttpPost handle http post method, you can set content type as
// comm.ContentTypeJson or comm.ContentTypeForm, or other you need set.
//
// ---
//
//	// set post data as json string
//	data := struct {"key": "Value", "id": "123"}
//	resp, err := comm.HttpPost(tagurl, data)
//
//	// set post data as form string
//	data := "key=Value&id=123"
//	resp, err := comm.HttpPost(tagurl, data, comm.ContentTypeForm)
func HttpPost(tagurl string, postdata interface{}, contentType ...string) ([]byte, error) {
	ct := ContentTypeJson
	if len(contentType) > 0 {
		ct = contentType[0]
	}

	if DebugPrintRequest {
		logger.D("Http Post:", tagurl, "ContentType:", ct)
	}

	switch ct {
	case ContentTypeJson:
		return httpPostJson(tagurl, postdata)
	case ContentTypeForm:
		return httpPostForm(tagurl, postdata.(url.Values))
	}
	return nil, invar.ErrInvalidParams
}

// HttpGetString call HttpGet and trim " char both begin and end
func HttpGetString(tagurl string, params ...interface{}) (string, error) {
	resp, err := HttpGet(tagurl, params...)
	if err != nil {
		return "", err
	}
	return strings.Trim(string(resp), "\""), nil
}

// HttpPostString call HttpPost and trim " char both begin and end.
func HttpPostString(tagurl string, postdata interface{}, contentType ...string) (string, error) {
	resp, err := HttpPost(tagurl, postdata, contentType...)
	if err != nil {
		return "", err
	}
	return strings.Trim(string(resp), "\""), nil
}

// HttpGetStruct handle http get method and unmarshal data to struct object
func HttpGetStruct(tagurl string, out interface{}, params ...interface{}) error {
	body, err := HttpGet(tagurl, params...)
	if err != nil {
		return err
	}
	return unmarshalResponse(body, out)
}

// HttpPostStruct handle http post method and unmarshal data to struct object
func HttpPostStruct(tagurl string, postdata, out interface{}, contentType ...string) error {
	body, err := HttpPost(tagurl, postdata, contentType...)
	if err != nil {
		return err
	}
	return unmarshalResponse(body, out)
}

// ==================================================

// HttpClientGet handle http get by http.Client, you can set request headers or
// ignore TLS verfiy of https url by setRequstFunc middle-ware function as :
//
// ---
//
//	comm.HttpClientGet(tagurl, func(req *http.Request) (bool, error) {
//			req.Header.Set("Content-Type", "application/json;charset=UTF-8")
//			req.SetBasicAuth("username", "password") // set auther header
//			return true, nil  // true is ignore TLS verify of https url
//		}, "same-params") ([]byte, error) {
//		// TODO do samething
//	}
func HttpClientGet(tagurl string, setRequestFunc SetRequest, params ...interface{}) ([]byte, error) {
	if len(params) > 0 {
		tagurl = fmt.Sprintf(tagurl, params...)
	}

	rawurl := EncodeUrl(tagurl)
	if DebugPrintRequest {
		logger.D("Http Client Get:", rawurl)
	}

	// generate new request instanse
	req, err := http.NewRequest(http.MethodGet, rawurl, http.NoBody)
	if err != nil {
		logger.E("Create http request err:", err)
		return nil, err
	}

	return httpClientDo(req, setRequestFunc)
}

// HttpClientPost handle https post by http.Client, you can set request headers or
// ignore TLS verfiy of https url by setRequstFunc middle-ware function as :
//
// ---
//
//	comm.HttpClientPost(tagurl, func(req *http.Request) (bool, error) {
//			req.Header.Set("Content-Type", "application/json;charset=UTF-8")
//			req.SetBasicAuth("username", "password") // set auther header
//			return true, nil  // true is ignore TLS verify of https url
//		}, "post-data") ([]byte, error) {
//		// TODO do samething
//	}
func HttpClientPost(tagurl string, setRequestFunc SetRequest, postdata ...interface{}) ([]byte, error) {
	var body io.Reader
	if len(postdata) > 0 {
		params, err := json.Marshal(postdata[0])
		if err != nil {
			logger.E("Marshal post data err:", err)
			return nil, err
		}
		body = bytes.NewReader(params)
	} else {
		body = http.NoBody
	}

	if DebugPrintRequest {
		logger.D("Http Client Post:", tagurl)
	}

	// generate new request instanse
	req, err := http.NewRequest(http.MethodPost, tagurl, body)
	if err != nil {
		logger.E("Create http request err:", err)
		return nil, err
	}

	// set json as default content type
	req.Header.Set("Content-Type", ContentTypeJson)
	return httpClientDo(req, setRequestFunc)
}

// HttpClientGetStruct handle http get method and unmarshal data to struct object
func HttpClientGetStruct(tagurl string, setRequestFunc SetRequest, out interface{}, params ...interface{}) error {
	body, err := HttpClientGet(tagurl, setRequestFunc, params...)
	if err != nil {
		return err
	}
	return unmarshalResponse(body, out)
}

// HttpClientPostStruct handle http post method and unmarshal data to struct object
func HttpClientPostStruct(tagurl string, setRequestFunc SetRequest, out interface{}, postdata ...interface{}) error {
	body, err := HttpClientPost(tagurl, setRequestFunc, postdata...)
	if err != nil {
		return err
	}
	return unmarshalResponse(body, out)
}

// httpClientDo handle http client DO method, and return response.
func httpClientDo(req *http.Request, setRequestFunc SetRequest) ([]byte, error) {
	client := &http.Client{}

	// use middle-ware to set request header
	if setRequestFunc != nil {
		ignoreTLS, err := setRequestFunc(req)
		if err != nil {
			logger.E("Set http request err:", err)
			return nil, err
		}

		logger.I("httpClientDo: ignore TLS:", ignoreTLS)
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: ignoreTLS,
			},
		}
	}

	// execute http request
	resp, err := client.Do(req)
	if err != nil {
		logger.E("Execute client DO, err:", err)
		return nil, err
	}

	defer resp.Body.Close()
	return readResponse(resp)
}
