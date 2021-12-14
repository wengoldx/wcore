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

// Upload image file to wechat platform by APIv3,
// it just support suffix in jpg, png, jpeg, bmp, and file size must be samll 2MB.
//	@param file Upload file content
//	@param header Upload file header information
//	@param ms Merchant secret informations
//	@return - string Media ID of wechat pay platform
//			- error Handled result
//
//	// use beego controller to get file and header
//	file, header, err := ctrl.GetFile("img")
//
// see more
//
// - [Wechat Image Upload](https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter2_1_1.shtml)
func (w *WxPayAgent) UploadImage(file io.Reader, header *multipart.FileHeader, ms *WxMerch) (string, error) {
	if header == nil || header.Size > (2*1024*1024) {
		logger.E("Null file header or file size oversized")
		return "", invar.ErrInvalidParams
	}

	filename := header.Filename
	suffix := strings.TrimLeft(strings.ToLower(path.Ext(filename)), ".")
	if len(suffix) == 0 || !(suffix == "jpg" || suffix == "png" || suffix == "jpeg" || suffix == "bmp") {
		logger.E("Invalid image file type, must be in jpg, jpeg, bmp, png")
		return "", invar.ErrInvalidParams
	}

	return w.wxpayAPIv3Upload(wxpApiUpImage, filename, suffix, file, ms)
}

// Upload video file to wechat platform by APIv3,
// it only support suffix in mp4, avi, wmv, mpeg, mov, mkv, flv, f4v, m4v, rmvb,
// and file size must be samll 5MB.
//	@param file Upload file content
//	@param header Upload file header information
//	@param ms Merchant secret informations
//	@return - string Media ID of wechat pay platform
//			- error Handled result
//
//	// use beego controller to get file and header
//	file, header, err := ctrl.GetFile("video")
//
// see more
//
// - [Wechat Video Upload](https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter2_1_2.shtml)
func (w *WxPayAgent) UploadVideo(file io.Reader, header *multipart.FileHeader, ms *WxMerch) (string, error) {
	if header == nil || header.Size > (5*1024*1024) {
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

	return w.wxpayAPIv3Upload(wxpApiUpVideo, filename, suffix, file, ms)
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
// [Wechat APIv3 - H5](https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter3_3_2.shtml),
// [Wechat APIv3 - App](https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter3_2_2.shtml),
// [Wechat APIv3 - JSAPI](https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter3_5_2.shtml)
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
// [Wechat APIv3 - H5](https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter3_3_2.shtml),
// [Wechat APIv3 - App](https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter3_2_2.shtml),
// [Wechat APIv3 - JSAPI](https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter3_5_2.shtml)
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
func (w *WxPayAgent) PFRegistry(body string, resp interface{}, ms *WxMerch) error {
	return w.postWxV3Http(WxpPFMchReg, body, resp, ms)
}

// Request change merchant bank by using wechat pay APIv3
func (w *WxPayAgent) PFChangBank(mid, body string, ms *WxMerch) error {
	return w.postWxV3Http(fmt.Sprintf(WxpMchAccMod, mid), body, nil, ms)
}

// Request merchant H5 pay action by using wechat pay APIv3
func (w *WxPayAgent) PFH5Pay(body string, resp interface{}, ms *WxMerch) error {
	return w.postWxV3Http(WxpMchH5, body, resp, ms)
}

// Request merchant app pay action by using wechat pay APIv3
func (w *WxPayAgent) PFAppPay(body string, resp interface{}, ms *WxMerch) error {
	return w.postWxV3Http(WxpMchApp, body, resp, ms)
}

// Request merchant JSAPI pay action by using wechat pay APIv3
func (w *WxPayAgent) PFJSPay(body string, resp interface{}, ms *WxMerch) error {
	return w.postWxV3Http(WxpMchJS, body, resp, ms)
}

// Request merchant pay refund action by using wechat pay APIv3
func (w *WxPayAgent) PFPayRefund(body string, resp interface{}, ms *WxMerch) error {
	return w.postWxV3Http(WxpPFRefund, body, resp, ms)
}

// Request merchant withdraw action by using wechat pay APIv3
func (w *WxPayAgent) PFWithdraw(body string, resp interface{}, ms *WxMerch) error {
	return w.postWxV3Http(WxpPFWithdraw, body, resp, ms)
}

// Request merchant dividing action by using wechat pay APIv3
func (w *WxPayAgent) PFDividing(body string, resp interface{}, ms *WxMerch) error {
	return w.postWxV3Http(WxpPFDividing, body, resp, ms)
}

// Request merchant dividing refund action by using wechat pay APIv3
func (w *WxPayAgent) PFDiviRefund(body string, resp interface{}, ms *WxMerch) error {
	return w.postWxV3Http(WxpPFDiviRefund, body, resp, ms)
}

// Request merchant close dividing action by using wechat pay APIv3
func (w *WxPayAgent) PFDiviClose(body string, resp interface{}, ms *WxMerch) error {
	return w.postWxV3Http(WxpPFDiviClose, body, resp, ms)
}

// Request merchant registry result by using wechat pay APIv3
func (w *WxPayAgent) PFRegQuery(regno string, resp interface{}, ms *WxMerch) error {
	return w.getWxV3Http(fmt.Sprintf(WxpPFMchRNoQuery, regno), resp, ms)
}

// Request merchant change bank result by using wechat pay APIv3
func (w *WxPayAgent) PFChgQuery(smid string, resp interface{}, ms *WxMerch) error {
	return w.getWxV3Http(fmt.Sprintf(WxpMchMQuery, smid), resp, ms)
}

// Request merchant query trade result by using wechat pay APIv3
func (w *WxPayAgent) PFQuery(tno, spid, smid string, resp interface{}, ms *WxMerch) error {
	return w.getWxV3Http(fmt.Sprintf(WxpMchNoQuery, tno, spid, smid), resp, ms)
}

// Request merchant query refund result by using wechat pay APIv3
func (w *WxPayAgent) PFRefQuery(rno, smid string, resp interface{}, ms *WxMerch) error {
	return w.getWxV3Http(fmt.Sprintf(WxpPFRNoQuery, rno, smid), resp, ms)
}

// Request merchant query balance result by using wechat pay APIv3
func (w *WxPayAgent) PFBalQuery(smid, acctype string, resp interface{}, ms *WxMerch) error {
	if acctype != "" {
		return w.getWxV3Http(fmt.Sprintf(WxpPFBalance, smid)+"?account_type="+acctype, resp, ms)
	}
	return w.getWxV3Http(fmt.Sprintf(WxpPFBalance, smid), resp, ms)
}

// Request merchant query balance end date by using wechat pay APIv3
func (w *WxPayAgent) PFEndQuery(smid, enddate string, resp interface{}, ms *WxMerch) error {
	if enddate != "" {
		return w.getWxV3Http(fmt.Sprintf(WxpPFEndDay, smid)+"?date="+enddate, resp, ms)
	}
	return w.getWxV3Http(fmt.Sprintf(WxpPFEndDay, smid), resp, ms)
}

// Request merchant query withdraw result by using wechat pay APIv3
func (w *WxPayAgent) PFWithdrawQuery(wno, smid string, resp interface{}, ms *WxMerch) error {
	return w.getWxV3Http(fmt.Sprintf(WxpPFWNoQuery, wno, smid), resp, ms)
}

// Request merchant query deviding result by using wechat pay APIv3
func (w *WxPayAgent) PFDiviQuery(smid, tid, dno string, resp interface{}, ms *WxMerch) error {
	return w.getWxV3Http(fmt.Sprintf(WxpPFDiviQuery, smid, tid, dno), resp, ms)
}

// Request merchant query deviding refund by using wechat pay APIv3
func (w *WxPayAgent) PFDRefQuery(smid, rno, tno string, resp interface{}, ms *WxMerch) error {
	return w.getWxV3Http(fmt.Sprintf(WxpPFDRefQuery, smid, tno, rno), resp, ms)
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

// Sign request auth header to access wechat pay APIv3
//	@param method Http method of GET, POST
//	@param urlpath Wechat pay platform APIv3 api
//	@param body Request data, mashaled to json string
//	@return - string Authentication header
//			- error Handled result
func (w *WxPayAgent) signAuthHeader(method, urlpath, body string, ms *WxMerch) (string, error) {
	// Step 1. generate nonce and timestamp strings
	noncestr := secure.GenNonce()
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	// Step 2. sign request packet datas,
	// https://pay.weixin.qq.com/wiki/doc/apiv3/wechatpay/wechatpay4_0.shtml#part-1
	signstr := w.SignPacket(method, urlpath, timestamp, noncestr, body)

	// Step 3. generate the signature string by rsa256,
	// https://pay.weixin.qq.com/wiki/doc/apiv3/wechatpay/wechatpay4_0.shtml#part-2
	signature, err := w.EncrpySign(ms.PriPem, signstr)
	if err != nil {
		logger.E("Faild to encripty signture, err:", err)
		return "", err
	}

	// Step 4. auth string,
	// https://pay.weixin.qq.com/wiki/doc/apiv3/wechatpay/wechatpay4_0.shtml#part-3
	authstr := w.AuthPacket(ms.MchID, noncestr, timestamp, ms.SerialNo, signature)
	return authstr, nil
}

// Generate request body data for upload image and
// video to wechat platform.
//
// see more
//
// - [Wechat Upload Image](https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter2_1_1.shtml#part-4)
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
//	@return - error Handled result
//
// see more
//
// - [Wechat Access Rules](https://pay.weixin.qq.com/wiki/doc/apiv3/wechatpay/wechatpay2_0.shtml),
// - [Wechat Authenticate](https://pay.weixin.qq.com/wiki/doc/apiv3/wechatpay/wechatpay4_0.shtml)
func (w *WxPayAgent) wxpayAPIv3Http(method, urlpath, body string, resp interface{}, ms *WxMerch) error {
	logger.D("Request wechat APIv3 ["+method+"]:", urlpath)
	url := WxpApisDomain + urlpath

	// Step 1. sign request authentication hearder
	authstr, err := w.signAuthHeader(method, urlpath, body, ms)
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
	req.Header.Set("Wechatpay-Serial", ms.PayPlatSN)

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
//	@param file Upload file content
//	@param filename Upload file name with suffix
//	@param suffix Upload file suffix
//	@return - string Media ID of wechat pay platform
//			- error Handled result
//
// see more
//
// - [Wechat Upload Image](https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter2_1_1.shtml),
// - [Wechat Upload Video](https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter2_1_2.shtml)
func (w *WxPayAgent) wxpayAPIv3Upload(urlpath, filename, suffix string, file io.Reader, ms *WxMerch) (string, error) {
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
	authstr, err := w.signAuthHeader(method, urlpath, string(filemeta), ms)
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
