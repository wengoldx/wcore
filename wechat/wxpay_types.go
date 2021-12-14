package wechat

// Wechat pay APIv3 apis, theme just for internal using
const (
	WxpApisDomain = "https://api.mch.weixin.qq.com"

	// Wechat certificate and upload medias
	wxpApiCert    = "/v3/certificates"
	wxpApiUpImage = "/v3/merchant/media/upload"
	wxpApiUpVideo = "v3/merchant/media/video_upload"

	// Wechat Direct connected merchants
	wxpDrH5      = "/v3/pay/transactions/h5"
	wxpDrApp     = "/v3/pay/transactions/app"
	wxpDrJS      = "/v3/pay/transactions/jsapi"
	wxpDrNative  = "/v3/pay/transactions/native" // not suppport yet of agent
	wxpDrIDQuery = "/v3/pay/transactions/id/%s?mchid=%s"
	wxpDrNoQuery = "/v3/pay/transactions/out-trade-no/%s?mchid=%s"
	wxpDrClose   = "/v3/pay/transactions/out-trade-no/%s/close"

	WxpDrRefund   = "/secapi/pay/refund"
	WxpDrRefQuery = "/pay/refundquery"

	// Wechat SP, PF payments and merch account change
	WxpMchApp     = "/v3/pay/partner/transactions/app"
	WxpMchJS      = "/v3/pay/partner/transactions/jsapi"
	WxpMchH5      = "/v3/pay/partner/transactions/h5"
	WxpMchIDQuery = "/v3/pay/partner/transactions/id/%s"
	WxpMchNoQuery = "/v3/pay/partner/transactions/out-trade-no/%s?sp_mchid=%s&sub_mchid=%s"
	WxpMchClose   = "/v3/pay/partner/transactions/out-trade-no/%s/close"
	WxpMchAccMod  = "/v3/apply4sub/sub_merchants/%s/modify-settlement"
	WxpMchMQuery  = "/v3/apply4sub/sub_merchants/%s/settlement"

	// Wechat Service provider
	WxpSPMchReg     = "/v3/applyment4sub/applyment/"
	WxpSPMchRegCode = "/v3/applyment4sub/applyment/business_code/%s"
	WxpSPMchRegID   = "/v3/applyment4sub/applyment/applyment_id/%s"

	// Wechat E-commerce platform URL
	WxpPFMchReg      = "/v3/ecommerce/applyments/"
	WxpPFMchRIDQuery = "/v3/ecommerce/applyments/%s"
	WxpPFMchRNoQuery = "/v3/ecommerce/applyments/out-request-no/%s"

	WxpPFBalance  = "/v3/ecommerce/fund/balance/%s"
	WxpPFEndDay   = "/v3/ecommerce/fund/enddaybalance/%s"
	WxpPFWithdraw = "/v3/ecommerce/fund/withdraw"
	WxpPFWIDQuery = "/v3/ecommerce/fund/withdraw/%s"
	WxpPFWNoQuery = "/v3/ecommerce/fund/withdraw/out-request-no/%s?sub_mchid=%s"

	WxpPFDividing   = "/v3/ecommerce/profitsharing/orders"
	WxpPFDiviRefund = "/v3/ecommerce/profitsharing/returnorders?sub_mchid=%s&transaction_id=%s&out_order_no=%s"
	WxpPFDiviClose  = "/v3/ecommerce/profitsharing/finish-order"
	WxpPFDiviQuery  = "/v3/ecommerce/profitsharing/returnorders?sub_mchid=%s&transaction_id=%s&out_order_no=%s"
	WxpPFDRefQuery  = "/v3/ecommerce/profitsharing/returnorders?sub_mchid=%s&out_order_no=%s&out_return_no=%s"
	WxpPFDiviRecAdd = "/v3/ecommerce/profitsharing/receivers/add"
	WxpPFDiviRecDel = "/v3/ecommerce/profitsharing/receivers/delete"

	WxpPFRefund   = "/v3/ecommerce/refunds/apply"
	WxpPFRIDQuery = "/v3/ecommerce/refunds/id/%s"
	WxpPFRNoQuery = "/v3/ecommerce/refunds/out-refund-no/%s?sub_mchid=%s"
)

// Custom boundary string of wxpay agent
//
// `WARNING` :
//
//	DO NOT Modify this string if YOU NOT KNOWN how to change.
const (
	wxpSignBoundary = "wengoldboundary"
)

// Custom payment status for service provider merchant,
// the values not same as wechat pay platform.
type PayState int

// payment status
const (
	WXP_NOTPAY PayState = iota
	WXP_SUCCESS
	WXP_CLOSED
	WXP_REFUND
	WXP_ERROR
	WXP_REVOKED
	WXP_PAYING
)

// -------- For Response Error

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

// MetaData image or video media data
type MetaData struct {
	FileName string `json:"filename" validate:"required" description:"upload media file name with suffix that must be JPG, JPEG, BMP, PNG on ignore char case"`
	HashCode string `json:"sha256"   validate:"required" description:"upload media hash code by sha256"`
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

// EncryptCert scene information of certificate
type EncryptCert struct {
	Algorithm      string `json:"algorithm"       description:"algorithm keyword, such as : AEAD_AES_256_GCM"`
	Nonce          string `json:"nonce"           description:"nonce string"`
	AssociatedData string `json:"associated_data" description:"associated data, such as : certificate"`
	Ciphertext     string `json:"ciphertext"      description:"ciphertext content"`
}

// Certificate wechat certificate
type Certificate struct {
	SerialNo      string      `json:"serial_no"           description:"serial number"`
	EffectiveTime string      `json:"effective_time"      description:"effective time, such as : 2018-06-08T10:34:56+08:00"`
	ExpireTime    string      `json:"expire_time"         description:"expire time, such as : 2018-12-08T10:34:56+08:00"`
	EncryptCert   EncryptCert `json:"encrypt_certificate" description:"encrypt certificate"`
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

// WxRetCert wechat pay platform certificate updated response
type WxRetCert struct {
	Datas []Certificate `json:"data"`
}

// WxRetUpload uploaded media id of wechat pay platform
type WxRetUpload struct {
	MediaID string `json:"media_id" description:"uploaded media id of wechat pay platform"`
}
