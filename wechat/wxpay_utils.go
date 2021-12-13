package wechat

import (
	"encoding/json"
	"errors"
	"github.com/wengoldx/wing/logger"
	"github.com/wengoldx/wing/secure"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Response code and message from wechat pay
//
// see more such as
//
// * [Common error codes](https://pay.weixin.qq.com/wiki/doc/apiv3/wxpay/Share/error_code.shtml)
// * [H5 Pay error codes](https://pay.weixin.qq.com/wiki/doc/apiv3/wechatpay/wechatpay2_1.shtml)
type WxCodeMsg struct {
	Code    string `json:"code"    description:"response code, such as : PARAM_ERROR"`
	Message string `json:"message" description:"response result message"`
}

// Response error detail from wechat pay
type WxErrDetail struct {
	Field    string `json:"field"    description:"error field, such as : /amount/currency"`
	Value    string `json:"value"    description:"error value, such as : XYZ"`
	Issue    string `json:"issue"    description:"error description, such as : Currency code is invalid"`
	Location string `json:"location" description:"error position, such as : body"`
}

// Http response error informations from wechat pay
//
// see more
// [Struct Define](https://pay.weixin.qq.com/wiki/doc/apiv3/wechatpay/wechatpay2_0.shtml#part-7)
type WxRetErr struct {
	WxCodeMsg
	Detail WxErrDetail `json:"detail" description:"response error details"`
}

// Merchant information of wechat pay platform
type WxMerch struct {
	MchID     string `description:"merchant cash transaction account"`
	SerialNo  string `description:"merchant certificate serial number"`
	PriPem    string `description:"merchant certificate private pem file"`
	PayPlatSN string `description:"wechat pay platform serial number"`
}

// Generate wechat signature string packet
//	@param method Get or POST http method
//	@param URL API access path
//	@param timestamp Signature timestamp string as unix int format
//	@param nonce Nonce string
//	@param body Signature requst datas
//	@return - string Signatrue packet string
//
// `WARNING` :
//
//	`DO NOT change the order of the signature strings`
//
// see more
// [Generate Signature String](https://pay.weixin.qq.com/wiki/doc/apiv3/wechatpay/wechatpay4_0.shtml#part-1)
func SignPacket(method, URL, timestamp, nonce, body string) string {
	packet := ""
	packet += method + "\n"
	packet += URL + "\n"
	packet += timestamp + "\n"
	packet += nonce + "\n"
	packet += body + "\n"
	return packet
}

// Generate wechat authorization string packet
//	@param mchid Merchant ID
//	@param nonce Nonce string
//	@param timesstemp Signature timestamp string as unix int format
//	@param serialno Certificate serial number
//	@param signture Signature secure datas
//	@return - string Auth packet string
//
// `WARNING` :
//
//	`DO NOT change the order of the signature strings`
//
// see more
// [Set Http Auth Header](https://pay.weixin.qq.com/wiki/doc/apiv3/wechatpay/wechatpay4_0.shtml#part-3)
func AuthPacket(mchid, nonce, timestamp, serialno, signature string) string {
	packet := ""
	packet += "WECHATPAY2-SHA256-RSA2048 "
	packet += "mchid=\"" + mchid + "\","
	packet += "nonce_str=\"" + nonce + "\","
	packet += "timestamp=\"" + timestamp + "\","
	packet += "serial_no=\"" + serialno + "\","
	packet += "signature=\"" + signature + "\""
	return packet
}

// Generate wechat notification string packet
//	@param timestamp Signature timestamp string as unix int format
//	@param nonce Nonce string
//	@param body Signature requst datas
//	@return - string Notify packet string
//
// `WARNING` :
//
//	`DO NOT change the order of the signature strings`
func NotifyPacket(timestamp, nonce, body string) string {
	packet := ""
	packet += timestamp + "\n"
	packet += nonce + "\n"
	packet += body + "\n"
	return packet
}

// Encrpty request signture string by given private pem key file
// as RSA PCS#8 format.
//	@param prifile RSA PCS#8 private pem file path
//	@param signstr To be encript signture string
//	@return - string Encrpty string
//			- error Handle result
func EncrpySign(prifile, signstr string) (string, error) {
	return secure.RSA2Sign4FB64(prifile, []byte(signstr))
}

// wechatV3Http use RSA2 signed []byte data, and hex to string,
//
// see more
//
// * [Wechat Access Rules](https://pay.weixin.qq.com/wiki/doc/apiv3/wechatpay/wechatpay2_0.shtml)
// * [Wechat Authenticate](https://pay.weixin.qq.com/wiki/doc/apiv3/wechatpay/wechatpay4_0.shtml)
//
//	@return - http.status, error
func WxpayAPIv3Http(method, urlpath, body string, resp interface{}, ms WxMerch) error {
	url := WxpApisDomain + urlpath
	logger.D("Request wechat APIv3 ["+method+"]:", urlpath)

	// Step 1. generate nonce and timestamp strings
	nonceStr := secure.GenNonce()
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	// Step 2. sign request packet datas,
	// https://pay.weixin.qq.com/wiki/doc/apiv3/wechatpay/wechatpay4_0.shtml#part-1
	signStr := SignPacket(method, urlpath, timestamp, nonceStr, body)

	// Step 3. generate the signature string by rsa256,
	// https://pay.weixin.qq.com/wiki/doc/apiv3/wechatpay/wechatpay4_0.shtml#part-2
	signature, err := EncrpySign(ms.PriPem, signStr)
	if err != nil {
		logger.E("Faild to encripty signture, err:", err)
		return err
	}

	// Step 4. auth string,
	// https://pay.weixin.qq.com/wiki/doc/apiv3/wechatpay/wechatpay4_0.shtml#part-3
	authStr := AuthPacket(ms.MchID, nonceStr, timestamp, ms.SerialNo, signature)

	// Step 5. generate request client and setup hearder
	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		logger.E("http NewRequest error, err:", err)
		return err
	}

	// https://pay.weixin.qq.com/wiki/doc/apiv3/wechatpay/wechatpay2_0.shtml#part-0
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// https://pay.weixin.qq.com/wiki/doc/apiv3/wechatpay/wechatpay3_1.shtml#part-3
	req.Header.Set("Wechatpay-Serial", ms.PayPlatSN)

	// https://pay.weixin.qq.com/wiki/doc/apiv3/wechatpay/wechatpay2_0.shtml#part-8
	req.Header.Set("User-Agent", "Go-http-client/1.1")

	// https://pay.weixin.qq.com/wiki/doc/apiv3/wechatpay/wechatpay3_3.shtml
	req.Header.Set("Authorization", authStr)

	// Step 6. send the request
	client := http.Client{}
	ret, err := client.Do(req)
	if err != nil {
		logger.E("Failed request to wechat, err:", err)
		return err
	}
	defer ret.Body.Close()
	logger.D("Response status code:", ret.StatusCode)

	// Step 7. read the response from wechat
	retbody, err := ioutil.ReadAll(ret.Body)
	if err != nil {
		logger.E("Failed read wechat response body, err:", err)
		return err
	}

	// Step 8. check response status
	// https://pay.weixin.qq.com/wiki/doc/apiv3/wechatpay/wechatpay2_1.shtml
	if ret.StatusCode != 200 && ret.StatusCode != 204 {
		errinfo := &WxRetErr{}
		if err := json.Unmarshal(retbody, errinfo); err != nil {
			logger.E("Unmarhsal error message err:", err)
			return err
		}
		return errors.New(errinfo.Code + "-" + errinfo.Message)
	}

	// Step 9. parse return datas if have
	if resp != nil {
		if err = json.Unmarshal(retbody, resp); err != nil {
			logger.E("Success request wechat APIv3, but unmarhsal response data err:", err)
			return err
		}
		logger.D("Unmarchsal wechat APIv3 response")
	}
	return nil
}
