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

// Response code and message from wechat pay
//
// see more such as
//
// - [Common error codes](https://pay.weixin.qq.com/wiki/doc/apiv3/wxpay/Share/error_code.shtml)
// - [H5 Pay error codes](https://pay.weixin.qq.com/wiki/doc/apiv3/wechatpay/wechatpay2_1.shtml)
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
// - see more
// [Struct Define](https://pay.weixin.qq.com/wiki/doc/apiv3/wechatpay/wechatpay2_0.shtml#part-7)
type WxRetErr struct {
	WxCodeMsg
	Detail WxErrDetail `json:"detail" description:"response error details"`
}

// Merchant secure informations of wechat pay platform
type WxMerch struct {
	MchID     string `description:"merchant id of wechat"`
	SerialNo  string `description:"merchant certificate serial number"`
	PriPem    string `description:"merchant certificate private pem file"`
	PayPlatSN string `description:"wechat pay platform serial number"`
}

// WxMchID wechat merchant id
type WxMchID struct {
	ID string `json:"mchid" description:"merchant id"`
}

// -------- For Request

// Payer amount payer
type Payer struct {
	OpenID string `json:"openid" validate:"required" description:"payer openid of wechat"`
}

// Amount settle information
type Amount struct {
	Total    int64  `json:"total" validate:"gt=0" description:"total order amount, unit is one cent CNY"`
	Currency string `json:"currency,omitempty"    description:"CNY: RMB. Domestic merchant account only supports RMB."`
}

// H5Info H5 scene information, used in H5 payment
type H5Info struct {
	Type     string `json:"type" validate:"required" description:"scene type, such as iOS, Android and Wap"`
	AppName  string `json:"app_name,omitempty"       description:"app name"`
	AppURL   string `json:"app_url,omitempty"        description:"app url "`
	BundleID string `json:"bundle_id,omitempty"      description:"bundle id of IOS system platform"`
	PkgName  string `json:"package_name,omitempty"   description:"package name of android system platform"`
}

// StoreInfo settle information
type StoreInfo struct {
	ID       int64  `json:"id" validate:"required" description:"store number of merchant"`
	Name     string `json:"name,omitempty"         description:"store name of merchant"`
	AreaCode string `json:"area_code,omitempty"    description:"area code"`
	Address  string `json:"address,omitempty"      description:"store address of merchant"`
}

// SceneInfo sales scene information
type SceneInfo struct {
	ClientIP  string     `json:"payer_client_ip" validate:"required" description:"payer client ip address"`
	DeviceID  string     `json:"device_id,omitempty"                 description:"merchant device number (store number or cashier ID)"`
	StoreInfo *StoreInfo `json:"store_info,omitempty"                description:"Merchant store information"`
	H5Info    *H5Info    `json:"h5_info,omitempty"                   description:"H5 scene information, only used in h5 payment"`
}

// SettleInfo settle information
type SettleInfo struct {
	ProfitSharing bool `json:"profit_sharing,omitempty" description:"settlement information"`
}

// GoodsDetail scene information
type GoodsDetail struct {
	MerchGID  string `json:"merchant_goods_id" validate:"required" description:"goods id of merchant"`
	WxpayGID  string `json:"wechatpay_goods_id,omitempty"          description:"goods id of wechat pay platform"`
	GoodsName string `json:"goods_name,omitempty"                  description:"goods name"`
	Guantity  int64  `json:"quantity"          validate:"gt=0"     description:"trade quantity"`
	UnitPrice int64  `json:"unit_price"        validate:"gt=0"     description:"unit price, from Fen or cent"`
}

// Detail discount information
type Detail struct {
	CostPrice   int64        `json:"cost_price,omitempty"   description:"original order price"`
	InvoiceID   string       `json:"invoice_id,omitempty"   description:"merchant trade invoice id"`
	GoodsDetail *GoodsDetail `json:"goods_detail,omitempty" description:"goods details"`
}

// -------- For Response

// PayAmount amount settle information
type PayAmount struct {
	Total         int64  `json:"total"          description:"total amount"`
	PayerTotal    int64  `json:"payer_total"    description:"pay amount from player"`
	Currency      string `json:"currency"       description:"currentcy type"`
	PayerCurrency string `json:"payer_currency" description:"payer currency type"`
}

// PayScene payment scene information
type PayScene struct {
	DeviceID string `json:"device_id" description:"device id of merchant, which device handle the trade"`
}

// CouponDetail promotion goods details
type CouponDetail struct {
	GoodsID        string `json:"goods_id"        description:"* goods id"`
	Quantity       int64  `json:"quantity"        description:"* goods quantity"`
	UnitPrice      int64  `json:"unit_price"      description:"* trade unit price, base Fen or cern"`
	DiscountAmount int64  `json:"discount_amount" description:"* goods discount amount"`
	GoodsRemark    string `json:"goods_remark"    description:"goods remark information"`
}

// PromotionDetail promotion details
type PromotionDetail struct {
	CouponID        string          `json:"coupon_id"            description:"* coupon unique id"`
	Name            string          `json:"name"                 description:"promotion action name"`
	Scope           string          `json:"scope"                description:"promotion scope"`
	Type            string          `json:"type"                 description:"promotion type"`
	Amount          int64           `json:"amount"               description:"* coupon amount"`
	StockID         string          `json:"stock_id"             description:"stock id"`
	WxpayContribute int64           `json:"wechatpay_contribute" description:"contribute from wechat pay platform"`
	MerchContribute int64           `json:"merchant_contribute"  description:"contribute from merchant"`
	OtherContribute int64           `json:"other_contribute"     description:"contribute from other organization"`
	Currency        string          `json:"currency"             description:"currency type, such as CNY"`
	CouponDetail    []*CouponDetail `json:"goods_detail"         description:"goods coupon detail list"`
}

// -------- For Agent Input

// WxDrH5 Request input data of H5 direct pay
type WxDrH5 struct {
	AppID      string      `json:"appid"        validate:"required" description:"official account of wechat"`
	MchID      string      `json:"mchid"        validate:"required" description:"merchant id of wechat"`
	Desc       string      `json:"description"  validate:"required" description:"goods description"`
	TradeNo    string      `json:"out_trade_no" validate:"required" description:"trade number of service provider system"`
	TimeExpire string      `json:"time_expire,omitempty"            description:"ticket expire time as unix seconds string"`
	Attach     string      `json:"attach,omitempty"                 description:"attach information"`
	NotifyURL  string      `json:"notify_url"   validate:"required" description:"result notify url send from wechat pay platform"`
	GoodsTag   string      `json:"goods_tag,omitempty"              description:"goods order discount mark"`
	Amount     *Amount     `json:"amount"       validate:"required" description:"trade amount information"`
	Detail     *Detail     `json:"detail,omitempty"                 description:"promotion detail"`
	Scene      *SceneInfo  `json:"scene_info"   validate:"required" description:"trade scene"`
	Settle     *SettleInfo `json:"settle_info,omitempty"            description:"settlement information"`
}

// WxDrApp Request input data of app direct pay
type WxDrApp struct {
	AppID      string      `json:"appid"        validate:"required" description:"official account of wechat"`
	MchID      string      `json:"mchid"        validate:"required" description:"merchant id of wechat"`
	Desc       string      `json:"description"  validate:"required" description:"goods description"`
	TradeNo    string      `json:"out_trade_no" validate:"required" description:"trade number of service provider system"`
	TimeExpire string      `json:"time_expire,omitempty"            description:"ticket expire time as unix seconds string"`
	Attach     string      `json:"attach,omitempty"                 description:"attach information"`
	NotifyURL  string      `json:"notify_url"   validate:"required" description:"result notify url send from wechat pay platform"`
	GoodsTag   string      `json:"goods_tag,omitempty"              description:"goods order discount mark"`
	Amount     *Amount     `json:"amount"       validate:"required" description:"trade amount information"`
	Detail     *Detail     `json:"detail,omitempty"                 description:"promotion detail"`
	Scene      *SceneInfo  `json:"scene_info,omitempty"             description:"trade scene"`
	Settle     *SettleInfo `json:"settle_info,omitempty"            description:"settlement information"`
}

// WxDrJS Request input data of JSAPI direct pay
type WxDrJS struct {
	AppID      string      `json:"appid"        validate:"required" description:"official account of wechat"`
	MchID      string      `json:"mchid"        validate:"required" description:"merchant id of wechat"`
	Desc       string      `json:"description"  validate:"required" description:"goods description"`
	TradeNo    string      `json:"out_trade_no" validate:"required" description:"trade number of service provider system"`
	TimeExpire string      `json:"time_expire,omitempty"            description:"ticket expire time as unix seconds string"`
	Attach     string      `json:"attach,omitempty"                 description:"attach information"`
	NotifyURL  string      `json:"notify_url"   validate:"required" description:"result notify url send from wechat pay platform"`
	GoodsTag   string      `json:"goods_tag,omitempty"              description:"goods order discount mark"`
	Amount     *Amount     `json:"amount"       validate:"required" description:"trade amount information"`
	Payer      *Payer      `json:"payer"        validate:"required" description:"trade payer"`
	Detail     *Detail     `json:"detail,omitempty"                 description:"promotion detail"`
	Scene      *SceneInfo  `json:"scene_info,omitempty"             description:"trade scene"`
	Settle     *SettleInfo `json:"settle_info,omitempty"            description:"settlement information"`
}

// -------- For Agent Response

// WxRetDrH5 Response result of JSAPI direct pay
type WxRetDrH5 struct {
	H5URL string `json:"h5_url" description:"redirect web url for H5 direct pay"`
}

// WxRetDrApp Response result of app direct pay
type WxRetDrApp struct {
	PrePayID string `json:"prepay_id" description:"pre-pay trade id"`
}

// WxRetDrJS Response result of JSAPI direct pay
type WxRetDrJS struct {
	PrePayID string `json:"prepay_id" description:"pre-pay trade id"`
}

// WxRetTicket query result of trade informations
type WxRetTicket struct {
	AppID       string             `json:"appid"            description:"* official account of wechat"`
	MchID       string             `json:"mchid"            description:"* merchant id of wechat"`
	TradeNo     string             `json:"out_trade_no"     description:"* trade number of service provider system"`
	TranID      string             `json:"transaction_id"   description:"transaction id of wechat pay platform"`
	TradeType   string             `json:"trade_type"       description:"transaction type, such as JSAPI, NATIVE, APP, MICROPAY, MWEB, FACEPAY"`
	TradeState  string             `json:"trade_state"      description:"* trade state, such as SUCCESS, REFUND, NOTPAY, CLOSED, REVOKED, USERPAYING, PAYERROR"`
	TradeDesc   string             `json:"trade_state_desc" description:"* trade result description"`
	BankType    string             `json:"bank_type"        description:"bank type, such as CMC and so on"`
	Attach      string             `json:"attach"           description:"attach information, it can be used as custom param on pay notify received"`
	SuccessTime string             `json:"success_time"     description:"success payment time, format as YYYY-MM-DDTHH:mm:ss+TIMEZONE"`
	Payer       *Payer             `json:"payer"            description:"* attach information"`
	Amount      *PayAmount         `json:"amount"           description:"amount settle information"`
	Scene       *PayScene          `json:"scene_info"       description:"trade scene"`
	Promotion   []*PromotionDetail `json:"promotion_detail" description:"promotion details"`
}

// --------
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
