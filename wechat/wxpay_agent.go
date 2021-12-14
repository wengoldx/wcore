package wechat

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/wengoldx/wing/invar"
	"github.com/wengoldx/wing/logger"
	"github.com/wengoldx/wing/secure"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type WxPayAgent struct {
}

// Append wechat pay APIv3 domain and params combined string
//	@param formatpath API path, it maybe have format keyword
//	@param param Default dynamic params to insert key value into formatpath
//	@return - string Full url link with param if have
func (w *WxPayAgent) Url(formatpath string, param ...string) string {
	path := formatpath
	if num := len(param); num > 0 {
		path = fmt.Sprintf(formatpath, param[0])
	}
	return WxpApisDomain + path
}

// Check the given pay state if valid defined
//	@param state Pay state to check
//	@return - bool True is defined pay state, false is undefined
func (w *WxPayAgent) IsValidState(state PayState) bool {
	return state >= WXP_NOTPAY && state <= WXP_PAYING
}

// Get pay state name
//	@param state Pay state value
//	@return - string Pay state name string
func (w *WxPayAgent) State(state PayState) string {
	switch state {
	case WXP_NOTPAY:
		return "NOTPAY"
	case WXP_SUCCESS:
		return "SUCCESS"
	case WXP_CLOSED:
		return "CLOSED"
	case WXP_REFUND:
		return "REFUND"
	case WXP_ERROR:
		return "ERROR"
	case WXP_REVOKED:
		return "REVOKED"
	case WXP_PAYING:
		return "PAYING"
	}
	return ""
}

// -----------------------------------------------------------
// For Common Functions
// -----------------------------------------------------------

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
// - see more
// [Generate Signature String](https://pay.weixin.qq.com/wiki/doc/apiv3/wechatpay/wechatpay4_0.shtml#part-1)
func (w *WxPayAgent) SignPacket(method, URL, timestamp, nonce, body string) string {
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
// - see more
// [Set Http Auth Header](https://pay.weixin.qq.com/wiki/doc/apiv3/wechatpay/wechatpay4_0.shtml#part-3)
func (w *WxPayAgent) AuthPacket(mchid, nonce, timestamp, serialno, signature string) string {
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
func (w *WxPayAgent) NotifyPacket(timestamp, nonce, body string) string {
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
func (w *WxPayAgent) EncrpySign(prifile, signstr string) (string, error) {
	return secure.RSA2Sign4FB64(prifile, []byte(signstr))
}

// -----------------------------------------------------------
// For Certificate Update
// -----------------------------------------------------------

func (w *WxPayAgent) UpdateCert(resp *WxRetCert, ms *WxMerch) error {
	return w.getWxV3Http(wxpApiCert, resp, ms)
}

// -----------------------------------------------------------
// For Direct Pay
// -----------------------------------------------------------

// Request direct H5 pay action by using wechat pay APIv3
//
// - see more
// [Wechat APIv3](https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter3_3_1.shtml)
func (w *WxPayAgent) DrH5Pay(params *WxDrH5, resp *WxRetDrH5, ms *WxMerch) error {
	return w.postWxV3Http(wxpDrH5, params, resp, ms)
}

// Request direct app pay action by using wechat pay APIv3
//
// - see more
// [Wechat APIv3](https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter3_2_1.shtml)
func (w *WxPayAgent) DrAppPay(params *WxDrApp, resp *WxRetDrApp, ms *WxMerch) error {
	return w.postWxV3Http(wxpDrApp, params, resp, ms)
}

// Request direct JSAPI pay action by using wechat pay APIv3
//
// - see more
// [Wechat APIv3](https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter3_5_1.shtml)
func (w *WxPayAgent) DrJSPay(params *WxDrJS, resp *WxRetDrJS, ms *WxMerch) error {
	return w.postWxV3Http(wxpDrJS, params, resp, ms)
}

// Request close direct pay action by using wechat pay APIv3
//	@param tno Merchant transaction number of service provider
//
// - see more
// [Wechat APIv3](https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter3_3_3.shtml)
func (w *WxPayAgent) DrClose(tno string, ms *WxMerch) error {
	if ms == nil || len(ms.MchID) == 0 {
		logger.E("Null merch data or empty merch id!")
		return invar.ErrInvalidParams
	}
	return w.postWxV3Http(fmt.Sprintf(wxpDrClose, tno), &WxMchID{ID: ms.MchID}, nil, ms)
}

// Request query direct pay ticket with wechat trade id by using wechat pay APIv3
//	@param tid Transaction id of wechat pay platform
//	@return - resp Trade tickey details
//
// Dest URL format as:
//
//	https://api.mch.weixin.qq.com/v3/pay/transactions/id/1217752501201407033233368018?mchid=1230000109
//
// Notice that use agent.DrTNoQuery() will return same response datas
//
// - see more
// [Wechat APIv3 - H5](https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter3_3_2.shtml) ;
// [Wechat APIv3 - App](https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter3_2_2.shtml) ;
// [Wechat APIv3 - JSAPI](https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter3_5_2.shtml) ;
func (w *WxPayAgent) DrTIDQuery(tid string, resp *WxRetTicket, ms *WxMerch) error {
	if ms == nil || len(ms.MchID) == 0 {
		logger.E("Null merch data or empty merch id!")
		return invar.ErrInvalidParams
	}
	return w.getWxV3Http(fmt.Sprintf(wxpDrIDQuery, tid, ms.MchID), resp, ms)
}

// Request query direct pay ticket with merchant trade no by using wechat pay APIv3
//	@param tno Merchant transaction number of service provider
//	@return - resp Trade tickey details
//
// Dest URL format as:
//
//	https://api.mch.weixin.qq.com/v3/pay/transactions/out-trade-no/1217752501201407033233368018?mchid=1230000109
//
// Notice that use agent.DrTIDQuery() will return same response datas
//
// - see more
// [Wechat APIv3 - H5](https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter3_3_2.shtml) ;
// [Wechat APIv3 - App](https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter3_2_2.shtml) ;
// [Wechat APIv3 - JSAPI](https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter3_5_2.shtml) ;
func (w *WxPayAgent) DrTNoQuery(tno string, resp *WxRetTicket, ms *WxMerch) error {
	if ms == nil || len(ms.MchID) == 0 {
		logger.E("Null merch data or empty merch id!")
		return invar.ErrInvalidParams
	}
	return w.getWxV3Http(fmt.Sprintf(wxpDrNoQuery, tno, ms.MchID), resp, ms)
}

// -----------------------------------------------------------
// For Merchant Pay
// -----------------------------------------------------------

// Request register a new merchant by using wechat pay APIv3
func (w *WxPayAgent) PFRegistry(body string, resp interface{}, ms WxMerch) error {
	return w.postWxV3Http(WxpPFMchReg, body, resp, &ms)
}

// Request change merchant bank by using wechat pay APIv3
func (w *WxPayAgent) PFChangBank(body string, ms WxMerch) error {
	return w.postWxV3Http(fmt.Sprintf(WxpMchAccMod, ms.MchID), body, nil, &ms)
}

// Request merchant H5 pay action by using wechat pay APIv3
func (w *WxPayAgent) PFH5Pay(body string, resp interface{}, ms WxMerch) error {
	return w.postWxV3Http(WxpMchH5, body, resp, &ms)
}

// Request merchant app pay action by using wechat pay APIv3
func (w *WxPayAgent) PFAppPay(body string, resp interface{}, ms WxMerch) error {
	return w.postWxV3Http(WxpMchApp, body, resp, &ms)
}

// Request merchant JSAPI pay action by using wechat pay APIv3
func (w *WxPayAgent) PFJSPay(body string, resp interface{}, ms WxMerch) error {
	return w.postWxV3Http(WxpMchJS, body, resp, &ms)
}

// Request merchant pay refund action by using wechat pay APIv3
func (w *WxPayAgent) PFPayRefund(body string, resp interface{}, ms WxMerch) error {
	return w.postWxV3Http(WxpPFRefund, body, resp, &ms)
}

// Request merchant withdraw action by using wechat pay APIv3
func (w *WxPayAgent) PFWithdraw(body string, resp interface{}, ms WxMerch) error {
	return w.postWxV3Http(WxpPFWithdraw, body, resp, &ms)
}

// Request merchant dividing action by using wechat pay APIv3
func (w *WxPayAgent) PFDividing(body string, resp interface{}, ms WxMerch) error {
	return w.postWxV3Http(WxpPFDividing, body, resp, &ms)
}

// Request merchant dividing refund action by using wechat pay APIv3
func (w *WxPayAgent) PFDiviRefund(body string, resp interface{}, ms WxMerch) error {
	return w.postWxV3Http(WxpPFDividing, body, resp, &ms)
}

// Request merchant close dividing action by using wechat pay APIv3
func (w *WxPayAgent) PFDiviClose(body string, resp interface{}, ms WxMerch) error {
	return w.postWxV3Http(WxpPFDiviClose, body, resp, &ms)
}

// Request merchant registry result by using wechat pay APIv3
func (w *WxPayAgent) PFRegQuery(regno string, resp interface{}, ms WxMerch) error {
	return w.getWxV3Http(fmt.Sprintf(WxpPFMchRNoQuery, regno), resp, &ms)
}

// Request merchant change bank result by using wechat pay APIv3
func (w *WxPayAgent) PFChgQuery(smid string, resp interface{}, ms WxMerch) error {
	return w.getWxV3Http(fmt.Sprintf(WxpMchMQuery, smid), resp, &ms)
}

// Request merchant query trade result by using wechat pay APIv3
func (w *WxPayAgent) PFQuery(ps string, resp interface{}, ms WxMerch) error {
	return w.getWxV3Http(fmt.Sprintf(WxpMchNoQuery, ps), resp, &ms)
}

// Request merchant query refund result by using wechat pay APIv3
func (w *WxPayAgent) PFRefQuery(rno, ps string, resp interface{}, ms WxMerch) error {
	return w.getWxV3Http(fmt.Sprintf(WxpPFRNoQuery, rno)+ps, resp, &ms)
}

// -----------------------------------------------------------

// Handle GET http method to access wechat pay APIv3
func (w *WxPayAgent) getWxV3Http(urlpath string, resp interface{}, ms *WxMerch) error {
	return w.wxpayAPIv3Http("GET", urlpath, "", resp, ms)
}

// Handle POST http method to access wechat pay APIv3
func (w *WxPayAgent) postWxV3Http(urlpath string, params, resp interface{}, ms *WxMerch) error {
	if params == nil || resp == nil || ms == nil || len(urlpath) == 0 {
		return invar.ErrInvalidParams
	}

	body, err := json.Marshal(params)
	if err != nil {
		logger.E("Marshal input params err:", err)
		return err
	}
	return w.wxpayAPIv3Http("POST", urlpath, string(body), resp, ms)
}

// Http getter or poster to access wechat pay APIv3
//	@return - error Handled result
//
// see more
//
// - [Wechat Access Rules](https://pay.weixin.qq.com/wiki/doc/apiv3/wechatpay/wechatpay2_0.shtml)
// - [Wechat Authenticate](https://pay.weixin.qq.com/wiki/doc/apiv3/wechatpay/wechatpay4_0.shtml)
func (w *WxPayAgent) wxpayAPIv3Http(method, urlpath, body string, resp interface{}, ms *WxMerch) error {
	url := WxpApisDomain + urlpath
	logger.D("Request wechat APIv3 ["+method+"]:", urlpath)

	// Step 1. generate nonce and timestamp strings
	nonceStr := secure.GenNonce()
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	// Step 2. sign request packet datas,
	// https://pay.weixin.qq.com/wiki/doc/apiv3/wechatpay/wechatpay4_0.shtml#part-1
	signStr := w.SignPacket(method, urlpath, timestamp, nonceStr, body)

	// Step 3. generate the signature string by rsa256,
	// https://pay.weixin.qq.com/wiki/doc/apiv3/wechatpay/wechatpay4_0.shtml#part-2
	signature, err := w.EncrpySign(ms.PriPem, signStr)
	if err != nil {
		logger.E("Faild to encripty signture, err:", err)
		return err
	}

	// Step 4. auth string,
	// https://pay.weixin.qq.com/wiki/doc/apiv3/wechatpay/wechatpay4_0.shtml#part-3
	authStr := w.AuthPacket(ms.MchID, nonceStr, timestamp, ms.SerialNo, signature)

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
