package wechat

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/wengoldx/wing/invar"
	"github.com/wengoldx/wing/logger"
	"github.com/wengoldx/wing/secure"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"path"
	"strconv"
	"strings"
	"time"
)

// A agent using wechat pay APIv3 to support pay as the follow ways
//
// #### Direct connection merchant pay
//
// The methods with Dr prefix abbreviation to provider H5, App, JSAPI ways to
// place a order, query trade informations, request refund, close trade, and
// update trade status when a trade pay success or failed.
//
// - see more
//
// [Direct Connection Merchat](https://pay.weixin.qq.com/wiki/doc/apiv3/index.shtml)
//
// #### Service provider merchant pay
//
// The methods with PF prefix abbreviation to provider payments for both service
// provider merchant, and sub merchant of custom shopping mall platform.
// the support functions not just as direct connection merchant pay, but extras
// contain dividing pay, dividing refund, dividing query, sub merchant register,
// registry query, balance query, sub merchant withdraw, and so on.
//
// - see more
//
// [Service Provider Merchat](https://pay.weixin.qq.com/wiki/doc/apiv3_partner/index.shtml)
//
// `USAGE` :
//
// Before use the WxPayAgent, you must register a 'official account of wechat',
// then config and download the valid account, merchant secret certificates,
// serial number, API keys and so on.
//
// ---
//
//	agent := &wechat.WxPayAgent{}
//	err := agent.DownCerts(out)			// download merchant certificates
//	mediaid, err := agent.UploadImage(file, header) // upload image to wechat platform
//	mediaid, err := agent.UploadVideo(file, header) // upload video to wechat platform
//	err := agent.DrH5Pay(ps, out)		// request pay for H5 way
//	err := agent.DrAppPay(ps, out)		// request pay for App way
//	err := agent.DrJSPay(ps, out)		// request pay for JS way
//	err := agent.DrClose(tno)			// close trade
//	err := agent.DrTIDQuery(tid, out)	// query trade ticket by transaction id
//	err := agent.DrTNoQuery(tno, out)	// query trade ticket by trade number of shopping mall platform
//	// ...
type WxPayAgent struct {
	Merch *WxMerch   `description:"merchant secure informations"`
	PPlat *WxPayPlat `description:"wechat pay platform secure informations"`
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
//
// [Generate Signature](https://pay.weixin.qq.com/wiki/doc/apiv3/wechatpay/wechatpay4_0.shtml#part-1)
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
//	@param signkey Secure signature string
//	@return - string Auth packet string
//
// `WARNING` :
//
//	`DO NOT change the order of the signature strings`
//
// - see more
//
// [Set Http Auth Header](https://pay.weixin.qq.com/wiki/doc/apiv3/wechatpay/wechatpay4_0.shtml#part-3)
func AuthPacket(mchid, nonce, timestamp, serialno, signkey string) string {
	packet := ""
	packet += "WECHATPAY2-SHA256-RSA2048 "
	packet += "mchid=\"" + mchid + "\","
	packet += "nonce_str=\"" + nonce + "\","
	packet += "timestamp=\"" + timestamp + "\","
	packet += "serial_no=\"" + serialno + "\","
	packet += "signature=\"" + signkey + "\""
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

// Encrpty request signature string by given private pem key file
// as RSA PCS#8 format.
//	@param prifile RSA PCS#8 private pem file path
//	@param signstr To be encript sign content string
//	@return - string Encrpty string
//			- error Exception message
func EncrpySign(prifile, signstr string) (string, error) {
	return secure.RSA2Sign4FB64(prifile, []byte(signstr))
}

// Verify response data from wechat when received result notification
//	@param pubfile Pay platform certificate (public pem file path)
//	@param signstr Need to verify content
//	@oaram signkey Secure signature string
//	@return - error Exception message
func VerifySign(pubfile, signstr string, signkey []byte) error {
	return secure.RSAVerify4F(pubfile, []byte(signstr), signkey)
}

// DecryptPacket decrypt response data from wechat by AES-256-GCM
//	@param ciphertext Certificate ciphertext
//	@param noncestr Nonce string
//	@param additional Addiional data
//	@param apiv3key Merchant pay platform APIv3 key
//	@return - string Decrypted response datas
//			- error Exception message
func DecryptPacket(ciphertext, noncestr, additional, apiv3key string) (string, error) {
	secretkey, additionalData := []byte(apiv3key), []byte(additional)
	return secure.GCMDecrypt(secretkey, ciphertext, noncestr, additionalData)
}

// Verify the request body if valid from wechat
//	@param req Http request
//	@param body Http request body
//	@return - error Exception message
func (w *WxPayAgent) VerifyRequest(req *http.Request, body string) error {
	timestamp := req.Header.Get("Wechatpay-Timestamp")
	noncestr := req.Header.Get("Wechatpay-Nonce")
	signb64key := req.Header.Get("Wechatpay-Signature")

	signkey, err := secure.Base64ToByte(signb64key)
	if err != nil {
		logger.E("Decode sinature by base64, err:", err)
		return err
	}

	signstr := NotifyPacket(timestamp, noncestr, body)
	if err = VerifySign(w.PPlat.CertPem, signstr, signkey); err != nil {
		logger.E("Verify request body err:", err)
		return err
	}
	return nil
}

// -----------------------------------------------------------
// For Certificate Download
// -----------------------------------------------------------

// Download the wechat merchant (who set as agent merch) all certificates
//	@param resp Wechat Merchant all certificates.
//	@return - error Exception message
//
// - see more
//
// [Get Merchant Certificate](https://pay.weixin.qq.com/wiki/doc/apiv3/wechatpay/wechatpay5_1.shtml),
// [Decrypt Certificate](https://pay.weixin.qq.com/wiki/doc/apiv3/wechatpay/wechatpay4_2.shtml)
func (w *WxPayAgent) DownCerts(resp *WxMchCerts) error {
	return w.getWxV3Http(wxpApiDownCert, resp)
}

// Upload image file to wechat platform by APIv3,
// it just support suffix in jpg, png, jpeg, bmp, and file size must be samll 2MB.
//	@param file Upload file content
//	@param header Upload file header information
//	@return - string Media ID of wechat pay platform
//			- error Exception message
//
//	// use beego controller to get file and header
//	file, header, err := ctrl.GetFile("img")
//
// - see more
//
// [Image Upload](https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter2_1_1.shtml)
func (w *WxPayAgent) UploadImage(file io.Reader, header *multipart.FileHeader) (string, error) {
	if header == nil || header.Size > (2*1024*1024) /* 2MB */ {
		logger.E("Null file header or file size oversized")
		return "", invar.ErrInvalidParams
	}

	filename := header.Filename
	suffix := strings.TrimLeft(strings.ToLower(path.Ext(filename)), ".")
	if len(suffix) == 0 || !(suffix == "jpg" || suffix == "png" || suffix == "jpeg" || suffix == "bmp") {
		logger.E("Invalid image file type, must be in jpg, jpeg, bmp, png")
		return "", invar.ErrInvalidParams
	}

	return w.wxpayAPIv3Upload(wxpApiUploadImage, filename, suffix, file)
}

// Upload video file to wechat platform by APIv3,
// it only support suffix in mp4, avi, wmv, mpeg, mov, mkv, flv, f4v, m4v, rmvb,
// and file size must be samll 5MB.
//	@param file Upload file content
//	@param header Upload file header information
//	@return - string Media ID of wechat pay platform
//			- error Exception message
//
//	// use beego controller to get file and header
//	file, header, err := ctrl.GetFile("video")
//
// - see more
//
// [Video Upload](https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter2_1_2.shtml)
func (w *WxPayAgent) UploadVideo(file io.Reader, header *multipart.FileHeader) (string, error) {
	if header == nil || header.Size > (5*1024*1024) /* 5MB */ {
		logger.E("Null file header or file size oversized")
		return "", invar.ErrInvalidParams
	}

	filename := header.Filename
	suffix := strings.TrimLeft(strings.ToLower(path.Ext(filename)), ".")
	if len(suffix) == 0 ||
		!(suffix == "mp4" || suffix == "avi" || suffix == "wmv" || suffix == "mpeg" || suffix == "mov" || suffix == "mkv" || suffix == "flv" || suffix == "f4v" || suffix == "m4v" || suffix == "rmvb") {
		logger.E("Invalid video file type, must be in mp4, avi, wmv, mpeg, mov, mkv, flv, f4v, m4v, rmvb")
		return "", invar.ErrInvalidParams
	}

	return w.wxpayAPIv3Upload(wxpApiUploadVideo, filename, suffix, file)
}

// -----------------------------------------------------------
// For Direct Pay
// -----------------------------------------------------------

// Request direct H5 pay action
//	@param params Request params for H5 pay APIv3.
//	@param resp Output request result.
//	@return - error Exception message
//
// - see more
//
// APIv3 [H5](https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter3_3_1.shtml)
func (w *WxPayAgent) DrH5Pay(params *WxDrH5, resp *WxRetDrH5) error {
	return w.postWxV3Http(wxpApiDrH5, params, resp)
}

// Request direct JSAPI pay action
//	@param params Request params for JSAPI pay APIv3.
//	@param resp Output request result.
//	@return - error Exception message
//
// - see more
//
// APIv3 [JSAPI](https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter3_5_1.shtml)
func (w *WxPayAgent) DrJSPay(params *WxDrJS, resp *WxRetDrJS) error {
	return w.postWxV3Http(wxpApiDrJS, params, resp)
}

// Request direct app pay action
//	@param params Request params for App pay APIv3.
//	@param resp Output request result.
//	@return - error Exception message
//
// - see more
//
// APIv3 [App](https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter3_2_1.shtml)
func (w *WxPayAgent) DrAppPay(params *WxDrApp, resp *WxRetDrApp) error {
	return w.postWxV3Http(wxpApiDrApp, params, resp)
}

// Request query direct pay ticket with wechat trade id by using wechat pay APIv3
//	@param tid Transaction id of wechat pay platform
//	@param resp Output request result.
//	@return - error Exception message
//
// Dest URL format as:
//
//	https://api.mch.weixin.qq.com/v3/pay/transactions/id/1217752501201407033233368018?mchid=1230000109
//
// Notice that use agent.DrTNoQuery() will return same response datas
//
// - see more
//
// APIv3 [H5](https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter3_3_2.shtml),
// [JSAPI](https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter3_1_2.shtml),
// [App](https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter3_2_2.shtml),
// [Wechat App](https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter3_5_2.shtml)
func (w *WxPayAgent) DrTIDQuery(tid string, resp *WxRetTicket) error {
	if w.Merch == nil || len(w.Merch.MchID) == 0 {
		logger.E("Null merch data or empty merch id!")
		return invar.ErrInvalidParams
	}
	return w.getWxV3Http(fmt.Sprintf(wxpApiDrIDQuery, tid, w.Merch.MchID), resp)
}

// Request query direct pay ticket with merchant trade no by using wechat pay APIv3
//	@param tno Merchant transaction number of service provider
//	@param resp Output request result.
//	@return - error Exception message
//
// Dest URL format as:
//
//	https://api.mch.weixin.qq.com/v3/pay/transactions/out-trade-no/1217752501201407033233368018?mchid=1230000109
//
// Notice that use agent.DrTIDQuery() will return same response datas
//
// - see more
//
// APIv3 [H5](https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter3_3_2.shtml),
// [JSAPI](https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter3_1_2.shtml),
// [App](https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter3_2_2.shtml),
// [Wechat App](https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter3_5_2.shtml)
func (w *WxPayAgent) DrTNoQuery(tno string, resp *WxRetTicket) error {
	if w.Merch == nil || len(w.Merch.MchID) == 0 {
		logger.E("Null merch data or empty merch id!")
		return invar.ErrInvalidParams
	}
	return w.getWxV3Http(fmt.Sprintf(wxpApiDrNoQuery, tno, w.Merch.MchID), resp)
}

// Request refund action by given trade number
//	@param param Request params for refund of APIv3.
//	@param resp Output request result.
//	@return - error Exception message
//
// - see more
//
// APIv3 [H5](https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter3_3_9.shtml),
// [JSAPI](https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter3_1_9.shtml),
// [App](https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter3_2_9.shtml),
// [Wechat App](https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter3_5_9.shtml)
func (w *WxPayAgent) DrRefund(params *WxDrRefund, resp *WxRetRefund) error {
	return w.postWxV3Http(wxpApiDrRefund, params, resp)
}

// Request query refund ticket with merchant trade no by using wechat pay APIv3
//	@param tno Merchant transaction number of service provider
//	@param resp Output request result.
//	@return - error Exception message
//
// Dest URL format as:
//
//	https://api.mch.weixin.qq.com/v3/refund/domestic/refunds/1217752501201407033233368018
//
// - see more
//
// APIv3 [H5](https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter3_3_10.shtml),
// [JSAPI](https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter3_1_10.shtml)，
// [App](https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter3_2_10.shtml),
// [Wechat App](https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter3_5_10.shtml)
func (w *WxPayAgent) DrRefundQuery(tno string, resp *WxRetRefund) error {
	return w.getWxV3Http(fmt.Sprintf(wxpApiDrRefQuery, tno), resp)
}

// Request close direct pay action by using wechat pay APIv3
//	@param tno Merchant transaction number of service provider
//	@return - error Exception message
//
// - see more
//
// [Close APIv3](https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter3_3_3.shtml)
func (w *WxPayAgent) DrClose(tno string) error {
	if w.Merch == nil || len(w.Merch.MchID) == 0 {
		logger.E("Null merch data or empty merch id!")
		return invar.ErrInvalidParams
	}
	return w.postWxV3Http(fmt.Sprintf(wxpApiDrClose, tno), &WxMchID{ID: w.Merch.MchID}, nil)
}

// -----------------------------------------------------------
// For Merchant Pay
// -----------------------------------------------------------

// Request merchant H5 pay action by using wechat pay APIv3
func (w *WxPayAgent) PFH5Pay(body string, resp interface{}) error {
	return w.postWxV3Http(wxpApiMchH5, body, resp)
}

// Request merchant JSAPI pay action by using wechat pay APIv3
func (w *WxPayAgent) PFJSPay(body string, resp interface{}) error {
	return w.postWxV3Http(wxpApiMchJS, body, resp)
}

// Request merchant app pay action by using wechat pay APIv3
func (w *WxPayAgent) PFAppPay(body string, resp interface{}) error {
	return w.postWxV3Http(wxpApiMchApp, body, resp)
}

// -----

// Request register a new merchant by using wechat pay APIv3
func (w *WxPayAgent) PFRegistry(body string, resp interface{}) error {
	return w.postWxV3Http(WxpApiPFMchReg, body, resp)
}

// Request change merchant bank by using wechat pay APIv3
func (w *WxPayAgent) PFChangBank(mid, body string) error {
	return w.postWxV3Http(fmt.Sprintf(WxpApiMchAccMod, mid), body, nil)
}

// Request merchant pay refund action by using wechat pay APIv3
func (w *WxPayAgent) PFPayRefund(body string, resp interface{}) error {
	return w.postWxV3Http(WxpApiPFRefund, body, resp)
}

// Request merchant withdraw action by using wechat pay APIv3
func (w *WxPayAgent) PFWithdraw(body string, resp interface{}) error {
	return w.postWxV3Http(WxpApiPFWithdraw, body, resp)
}

// Request merchant dividing action by using wechat pay APIv3
func (w *WxPayAgent) PFDividing(body string, resp interface{}) error {
	return w.postWxV3Http(WxpApiPFDividing, body, resp)
}

// Request merchant dividing refund action by using wechat pay APIv3
func (w *WxPayAgent) PFDiviRefund(body string, resp interface{}) error {
	return w.postWxV3Http(WxpApiPFDiviRefund, body, resp)
}

// Request merchant close dividing action by using wechat pay APIv3
func (w *WxPayAgent) PFDiviClose(body string, resp interface{}) error {
	return w.postWxV3Http(WxpApiPFDiviClose, body, resp)
}

// Request merchant registry result by using wechat pay APIv3
func (w *WxPayAgent) PFRegQuery(regno string, resp interface{}) error {
	return w.getWxV3Http(fmt.Sprintf(WxpApiPFMchRNoQuery, regno), resp)
}

// Request merchant change bank result by using wechat pay APIv3
func (w *WxPayAgent) PFChgQuery(smid string, resp interface{}) error {
	return w.getWxV3Http(fmt.Sprintf(WxpApiMchMQuery, smid), resp)
}

// Request merchant query trade result by using wechat pay APIv3
func (w *WxPayAgent) PFQuery(tno, spid, smid string, resp interface{}) error {
	return w.getWxV3Http(fmt.Sprintf(WxpApiMchNoQuery, tno, spid, smid), resp)
}

// Request merchant query refund result by using wechat pay APIv3
func (w *WxPayAgent) PFRefQuery(rno, smid string, resp interface{}) error {
	return w.getWxV3Http(fmt.Sprintf(WxpApiPFRNoQuery, rno, smid), resp)
}

// Request merchant query balance result by using wechat pay APIv3
func (w *WxPayAgent) PFBalQuery(smid, acctype string, resp interface{}) error {
	if acctype != "" {
		return w.getWxV3Http(fmt.Sprintf(WxpApiPFBalance, smid)+"?account_type="+acctype, resp)
	}
	return w.getWxV3Http(fmt.Sprintf(WxpApiPFBalance, smid), resp)
}

// Request merchant query balance end date by using wechat pay APIv3
func (w *WxPayAgent) PFEndQuery(smid, enddate string, resp interface{}) error {
	if enddate != "" {
		return w.getWxV3Http(fmt.Sprintf(WxpApiPFEndDay, smid)+"?date="+enddate, resp)
	}
	return w.getWxV3Http(fmt.Sprintf(WxpApiPFEndDay, smid), resp)
}

// Request merchant query withdraw result by using wechat pay APIv3
func (w *WxPayAgent) PFWithdrawQuery(wno, smid string, resp interface{}) error {
	return w.getWxV3Http(fmt.Sprintf(WxpApiPFWNoQuery, wno, smid), resp)
}

// Request merchant query deviding result by using wechat pay APIv3
func (w *WxPayAgent) PFDiviQuery(smid, tid, dno string, resp interface{}) error {
	return w.getWxV3Http(fmt.Sprintf(WxpApiPFDiviQuery, smid, tid, dno), resp)
}

// Request merchant query deviding refund by using wechat pay APIv3
func (w *WxPayAgent) PFDRefQuery(smid, rno, tno string, resp interface{}) error {
	return w.getWxV3Http(fmt.Sprintf(WxpApiPFDRefQuery, smid, tno, rno), resp)
}

// -----------------------------------------------------------

// Handle GET http method to access wechat pay APIv3
func (w *WxPayAgent) getWxV3Http(urlpath string, resp interface{}) error {
	return w.wxpayAPIv3Http("GET", urlpath, "", resp)
}

// Handle POST http method to access wechat pay APIv3
func (w *WxPayAgent) postWxV3Http(urlpath string, params, resp interface{}) error {
	if params == nil || resp == nil || len(urlpath) == 0 {
		return invar.ErrInvalidParams
	}

	body, err := json.Marshal(params)
	if err != nil {
		logger.E("Marshal input params err:", err)
		return err
	}
	return w.wxpayAPIv3Http("POST", urlpath, string(body), resp)
}

// Sign request auth header to access wechat pay APIv3
//	@param method Http method of GET, POST
//	@param urlpath Wechat pay platform APIv3 api
//	@param body Request data, mashaled to json string
//	@return - string Authentication header
//			- error Exception message
func (w *WxPayAgent) signAuthHeader(method, urlpath, body string) (string, error) {
	// Step 1. generate nonce and timestamp strings
	noncestr := secure.GenNonce()
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	// Step 2. sign request packet datas,
	// https://pay.weixin.qq.com/wiki/doc/apiv3/wechatpay/wechatpay4_0.shtml#part-1
	signstr := SignPacket(method, urlpath, timestamp, noncestr, body)

	// Step 3. generate the signature string by rsa256,
	// https://pay.weixin.qq.com/wiki/doc/apiv3/wechatpay/wechatpay4_0.shtml#part-2
	signkey, err := EncrpySign(w.Merch.PriPem, signstr)
	if err != nil {
		logger.E("Faild to encripty signature, err:", err)
		return "", err
	}

	// Step 4. auth string,
	// https://pay.weixin.qq.com/wiki/doc/apiv3/wechatpay/wechatpay4_0.shtml#part-3
	authstr := AuthPacket(w.Merch.MchID, noncestr, timestamp, w.Merch.SerialNo, signkey)
	return authstr, nil
}

// Generate request body data for upload image and video to wechat platform,
// the data formated example like below :
//
// ------------------------------------------------------------
//
//	--boundary
//	Content-Disposition: form-data; name="meta";
//	Content-Type: application/json
//
//	{ "filename": "filea.jpg", "sha256": "hjkahkjsjkfsjk78687dhjahdajhk" }
//	--boundary
//	Content-Disposition: form-data; name="file"; filename="filea.jpg";
//	Content-Type: image/jpg
//
//	pic1xxxbuffersxxx
//	--boundary--
//
// ------------------------------------------------------------
//
// - see more
//
// [Upload Image](https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter2_1_1.shtml#part-4)
func (w *WxPayAgent) genUploadBody(fn, suffix string, fm, buff []byte) (*bytes.Buffer, string, error) {
	// Step 1. creter and setup body informations
	body := &bytes.Buffer{}
	bodywriter := multipart.NewWriter(body)
	bodywriter.SetBoundary(wxpSignBoundary)

	// Step 2. set body content type
	header := make(textproto.MIMEHeader)
	header.Set("Content-Disposition", "form-data;name=\"meta\";")
	header.Set("Content-Type", "application/json;")
	bodypart, err := bodywriter.CreatePart(header)
	if err != nil {
		return nil, "", err
	}

	// Step 3. set file meta information
	if _, err = bodypart.Write(fm); err != nil {
		return nil, "", err
	}

	// Step 4. set file content type
	header.Set("Content-Disposition", "form-data;name=\"file\";filename=\""+fn+"\";")
	header.Set("Content-Type", "image/"+suffix+";")
	bodypart, err = bodywriter.CreatePart(header)
	if err != nil {
		return nil, "", err
	}

	// Step 5. set file content data to body
	if _, err = bodypart.Write(buff); err != nil {
		return nil, "", err
	}

	// Step 6. close body writer
	if err := bodywriter.Close(); err != nil {
		return nil, "", err
	}

	contenttype := bodywriter.FormDataContentType()
	return body, contenttype, nil
}

// Http getter or poster to access wechat pay APIv3
//	@param method Http method of GET, POST
//	@param urlpath Wechat pay platform APIv3 api
//	@param body Request data, mashaled to json string
//	@return - error Exception message
//
// `WARNING` :
//
// It maybe return invar.ErrInvalidClient error when agent merchant or
// pay platform fields not set, such as:
//
//	w.Merch.MchID
//	w.Merch.PriPem
//	w.Merch.SerialNo
//	w.PPlat.SerialNo
//
// - see more
//
// [Access Rules](https://pay.weixin.qq.com/wiki/doc/apiv3/wechatpay/wechatpay2_0.shtml),
// [Authenticate](https://pay.weixin.qq.com/wiki/doc/apiv3/wechatpay/wechatpay4_0.shtml)
func (w *WxPayAgent) wxpayAPIv3Http(method, urlpath, body string, resp interface{}) error {
	logger.D("Request wechat APIv3 ["+method+"]:", urlpath)
	url := WxpApisDomain + urlpath

	// check agent merchant and pay platform configs if valid
	if w.Merch.MchID == "" || w.Merch.PriPem == "" || w.Merch.SerialNo == "" || w.PPlat.SerialNo == "" {
		logger.E("Invalid merchant or pay platform values of agent!")
		return invar.ErrInvalidClient
	}

	// Step 1. sign request authentication hearder
	authstr, err := w.signAuthHeader(method, urlpath, body)
	if err != nil {
		return err
	}

	// Step 2. generate request client and setup hearder
	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		logger.E("New a http request err:", err)
		return err
	}

	// https://pay.weixin.qq.com/wiki/doc/apiv3/wechatpay/wechatpay2_0.shtml#part-0
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// https://pay.weixin.qq.com/wiki/doc/apiv3/wechatpay/wechatpay3_1.shtml#part-3
	req.Header.Set("Wechatpay-Serial", w.PPlat.SerialNo)

	// https://pay.weixin.qq.com/wiki/doc/apiv3/wechatpay/wechatpay2_0.shtml#part-8
	req.Header.Set("User-Agent", "Go-http-client/1.1")

	// https://pay.weixin.qq.com/wiki/doc/apiv3/wechatpay/wechatpay3_3.shtml
	req.Header.Set("Authorization", authstr)

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

// Upload image or video file to wechat platform by APIv3
//	@param urlpath Wechat APIv3 url path
//	@param filename Upload file name with suffix
//	@param suffix Upload file suffix
//	@param file Upload file content
//	@return - string Media ID of wechat pay platform
//			- error Exception message
//
// - see more
//
// Upload [Image](https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter2_1_1.shtml),
// [Video](https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter2_1_2.shtml)
func (w *WxPayAgent) wxpayAPIv3Upload(urlpath, filename, suffix string, file io.Reader) (string, error) {
	logger.D("Upload file "+filename+" to wechat by APIv3 [ POST ]:", urlpath)
	url, method := WxpApisDomain+urlpath, "POST"

	// Step 1. read file content bytes
	filebuff, err := ioutil.ReadAll(file)
	if err != nil {
		logger.E("Failed read file content, err:", err)
		return "", err
	}

	// Step 2. hash file by sha256, and meta json data
	// https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter2_1_1.shtml#part-4
	filehash := secure.HashSHA256Hex(filebuff)
	metadata := &MetaData{FileName: filename, HashCode: filehash}
	filemeta, err := json.Marshal(metadata)
	if err != nil {
		logger.E("Marsh upload file meta data err:", err)
		return "", err
	}

	// Step 3. sign request authentication hearder
	authstr, err := w.signAuthHeader(method, urlpath, string(filemeta))
	if err != nil {
		return "", err
	}

	// Step 4. generate request body data
	body, contenttype, err := w.genUploadBody(filename, suffix, filemeta, filebuff)
	if err != nil {
		logger.E("Generate upload request body err:", err)
		return "", err
	}

	// Step 5. generate request client and setup hearder
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		logger.E("New a http request err:", err)
		return "", err
	}

	req.Header.Set("Content-Type", contenttype)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Charset", "UTF-8")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3100.0 Safari/537.36")
	req.Header.Set("Authorization", authstr)

	// Step 6. send the request
	client := http.Client{}
	ret, err := client.Do(req)
	if err != nil {
		logger.E("Failed upload file, err:", err)
		return "", err
	}
	defer ret.Body.Close()
	logger.D("Response status code:", ret.StatusCode)

	// Step 7. read the response from wechat
	retbody, err := ioutil.ReadAll(ret.Body)
	if err != nil {
		logger.E("Failed read wechat response body, err:", err)
		return "", err
	}

	// Step 8. check response status
	// https://pay.weixin.qq.com/wiki/doc/apiv3/wechatpay/wechatpay2_1.shtml
	if ret.StatusCode != 200 && ret.StatusCode != 204 {
		errinfo := &WxRetErr{}
		if err := json.Unmarshal(retbody, errinfo); err != nil {
			logger.E("Unmarhsal error message err:", err)
			return "", err
		}
		return "", errors.New(errinfo.Code + "-" + errinfo.Message)
	}

	// Step 9. parse return datas if have
	resp := &WxRetUpload{}
	if err = json.Unmarshal(retbody, resp); err != nil {
		logger.E("Success request wechat APIv3, but unmarhsal response data err:", err)
		return "", err
	}
	logger.D("Upload file and received media id:", resp.MediaID)
	return resp.MediaID, nil
}
