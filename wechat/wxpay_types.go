package wechat

// Wechat pay APIv3 apis, theme just for internal using
const (
	WxpApisDomain = "https://api.mch.weixin.qq.com"

	// Wechat certificate download and upload medias
	wxpApiDownCert    = "/v3/certificates"
	wxpApiUploadImage = "/v3/merchant/media/upload"
	wxpApiUploadVideo = "v3/merchant/media/video_upload"

	// Wechat Direct connected merchants
	wxpApiDrH5       = "/v3/pay/transactions/h5"
	wxpApiDrJS       = "/v3/pay/transactions/jsapi"
	wxpApiDrApp      = "/v3/pay/transactions/app"
	wxpApiDrIDQuery  = "/v3/pay/transactions/id/%s?mchid=%s"
	wxpApiDrNoQuery  = "/v3/pay/transactions/out-trade-no/%s?mchid=%s"
	wxpApiDrClose    = "/v3/pay/transactions/out-trade-no/%s/close"
	wxpApiDrRefund   = "/v3/refund/domestic/refunds"
	wxpApiDrRefQuery = "/v3/refund/domestic/refunds/%s"

	// Wechat SP, PF payments and merch account change
	wxpApiMchH5      = "/v3/pay/partner/transactions/h5"
	wxpApiMchJS      = "/v3/pay/partner/transactions/jsapi"
	wxpApiMchApp     = "/v3/pay/partner/transactions/app"
	WxpApiMchIDQuery = "/v3/pay/partner/transactions/id/%s"
	WxpApiMchNoQuery = "/v3/pay/partner/transactions/out-trade-no/%s?sp_mchid=%s&sub_mchid=%s"
	WxpApiMchClose   = "/v3/pay/partner/transactions/out-trade-no/%s/close"
	WxpApiMchAccMod  = "/v3/apply4sub/sub_merchants/%s/modify-settlement"
	WxpApiMchMQuery  = "/v3/apply4sub/sub_merchants/%s/settlement"

	// Wechat Service provider
	WxpApiSPMchReg     = "/v3/applyment4sub/applyment/"
	WxpApiSPMchRegCode = "/v3/applyment4sub/applyment/business_code/%s"
	WxpApiSPMchRegID   = "/v3/applyment4sub/applyment/applyment_id/%s"

	// Wechat E-commerce platform URL
	WxpApiPFMchReg      = "/v3/ecommerce/applyments/"
	WxpApiPFMchRIDQuery = "/v3/ecommerce/applyments/%s"
	WxpApiPFMchRNoQuery = "/v3/ecommerce/applyments/out-request-no/%s"

	WxpApiPFBalance  = "/v3/ecommerce/fund/balance/%s"
	WxpApiPFEndDay   = "/v3/ecommerce/fund/enddaybalance/%s"
	WxpApiPFWithdraw = "/v3/ecommerce/fund/withdraw"
	WxpApiPFWIDQuery = "/v3/ecommerce/fund/withdraw/%s"
	WxpApiPFWNoQuery = "/v3/ecommerce/fund/withdraw/out-request-no/%s?sub_mchid=%s"

	WxpApiPFDividing   = "/v3/ecommerce/profitsharing/orders"
	WxpApiPFDiviRefund = "/v3/ecommerce/profitsharing/returnorders?sub_mchid=%s&transaction_id=%s&out_order_no=%s"
	WxpApiPFDiviClose  = "/v3/ecommerce/profitsharing/finish-order"
	WxpApiPFDiviQuery  = "/v3/ecommerce/profitsharing/returnorders?sub_mchid=%s&transaction_id=%s&out_order_no=%s"
	WxpApiPFDRefQuery  = "/v3/ecommerce/profitsharing/returnorders?sub_mchid=%s&out_order_no=%s&out_return_no=%s"
	WxpApiPFDiviRecAdd = "/v3/ecommerce/profitsharing/receivers/add"
	WxpApiPFDiviRecDel = "/v3/ecommerce/profitsharing/receivers/delete"

	WxpApiPFRefund   = "/v3/ecommerce/refunds/apply"
	WxpApiPFRIDQuery = "/v3/ecommerce/refunds/id/%s"
	WxpApiPFRNoQuery = "/v3/ecommerce/refunds/out-refund-no/%s?sub_mchid=%s"

	// For wechat pay APIv2 refund
	WxpDrRefund   = "/secapi/pay/refund"
	WxpDrRefQuery = "/pay/refundquery"
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

// Merchant secure informations
//
// - see more
//
// [Certificate Usage](https://pay.weixin.qq.com/wiki/doc/apiv3/wechatpay/wechatpay3_0.shtml),
// [APIv3 Key](https://pay.weixin.qq.com/wiki/doc/apiv3/wechatpay/wechatpay3_2.shtml),
// [APIv3 Key Usecase](https://pay.weixin.qq.com/wiki/doc/apiv3/wechatpay/wechatpay5_1.shtml)
type WxMerch struct {
	MchID    string `description:"merchant id of wechat"`
	SerialNo string `description:"merchant certificate serial number"`
	PriPem   string `description:"merchant certificate private pem file, such as apiclient_key.pem"`
	PubPem   string `description:"merchant certificate public pem file, such as apiclient_cert.pem"`
	APIv3Key string `description:"merchant APIv3 secure key, use for decrypt response signtrue data from wechat"`
}

// Wechat merchant pay platform certificates datas
//
// - see more
//
// [Pay Platform Certificate](https://pay.weixin.qq.com/wiki/doc/apiv3/wechatpay/wechatpay3_1.shtml#part-2)
type WxPayPlat struct {
	CertPem  string `description:"wechat pay platform certificate public pem file, such as wxp_cert.pem"`
	SerialNo string `description:"wechat pay platform certificate serial number"`
	Expire   int64  `description:"unix seconds time for next refresh certificate"`
}

// Wxpay casher, contain pay provider merchant, pay platform secures, APIv3 access agent
type WxCashier struct {
	SKey   string      `description:"cashier custom uniqe key, use for cache or seach with map"`
	AppID  string      `description:"which app to handle the trade transaction"`
	RootCa string      `description:"TLS authenticate certificate file, such as Root_CA.pem"`
	Agent  *WxPayAgent `description:"Wxpay agent to access wechat pay APIv3 apis"`
}

// -------- For Response Error

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
//
// [Common Errors](https://pay.weixin.qq.com/wiki/doc/apiv3/wxpay/Share/error_code.shtml),
// [H5 Pay Errors](https://pay.weixin.qq.com/wiki/doc/apiv3/wechatpay/wechatpay2_1.shtml),
// [Error Struct](https://pay.weixin.qq.com/wiki/doc/apiv3/wechatpay/wechatpay2_0.shtml#part-7)
type WxRetErr struct {
	Code    string      `json:"code"    description:"response code, such as : PARAM_ERROR"`
	Message string      `json:"message" description:"response result message"`
	Detail  WxErrDetail `json:"detail"  description:"response error details"`
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

// DiscountDetail discount detail for request pay
type DiscountDetail struct {
	CostPrice   int64        `json:"cost_price,omitempty"   description:"original order price"`
	InvoiceID   string       `json:"invoice_id,omitempty"   description:"merchant trade invoice id"`
	GoodsDetail *GoodsDetail `json:"goods_detail,omitempty" description:"goods details"`
}

// MetaData image or video media data
type MetaData struct {
	FileName string `json:"filename" validate:"required" description:"upload media file name with suffix that must be JPG, JPEG, BMP, PNG on ignore char case"`
	HashCode string `json:"sha256"   validate:"required" description:"upload media hash code by sha256"`
}

// RefundAmount amount settle information
type RefundFrom struct {
	Account string `json:"account" validate:"required" description:"contribution account type"`
	Amount  int64  `json:"amount"  validate:"gt=0"     description:"contribution amount from this account"`
}

// RefundAmount amount settle information
type RefundAmount struct {
	RefundTotal int64         `json:"refund"   validate:"gt=0"     description:"refund amount of request"`
	Froms       []*RefundFrom `json:"from,omitempty"               description:"refund contribution account and amount"`
	Total       int64         `json:"total"    validate:"gt=0"     description:"total amount of original trade ticket"`
	Currency    string        `json:"currency" validate:"required" description:"currentcy type"`
}

// RefundGoods goods detail of refund
type RefundGoods struct {
	MerchGID  string `json:"merchant_goods_id" validate:"required" description:"goods id of merchant"`
	WxpayGID  string `json:"wechatpay_goods_id,omitempty"          description:"goods id of wechat pay platform"`
	GoodsName string `json:"goods_name,omitempty"                  description:"goods name"`
	UnitPrice int64  `json:"unit_price"        validate:"gt=0"     description:"unit price, from Fen or cent"`
	Amount    int64  `json:"refund_amount"     validate:"gt=0"     description:"refund amount of goods price"`
	Quantity  int64  `json:"refund_quantity"   validate:"gt=0"     description:"goods quantify of refund action"`
}

// -------- For Response, * is required

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

// RefRetAmount amount settle information for refund response
type RefRetAmount struct {
	Total        int64         `json:"total"             description:"* total amount"`
	RefundTotal  int64         `json:"refund"            description:"* refund amount"`
	Froms        []*RefundFrom `json:"from"              description:"refund contribution account and amount"`
	PayerTotal   int64         `json:"payer_total"       description:"* pay amount from player"`
	PayerRefund  int64         `json:"payer_refund"      description:"* refund amount to player"`
	SettleRefund int64         `json:"settlement_refund" description:"* settlement refund amount"`
	SettleTotal  int64         `json:"settlement_total"  description:"* settlement total"`
	Discount     int64         `json:"discount_refund"   description:"* discount refund amount"`
	Currency     string        `json:"currency"          description:"* currentcy type"`
}

// RefRetPromot refund promotion details
type RefRetPromot struct {
	PromotID     string         `json:"promotion_id"   description:"* promotion unique id"`
	Scope        string         `json:"scope"          description:"* promotion scope"`
	Type         string         `json:"type"           description:"* promotion type"`
	Amount       int64          `json:"amount"         description:"* promotion amount"`
	RefundAmount int64          `json:"refund_amount"  description:"* refund amount of promotion"`
	CouponDetail []*RefundGoods `json:"goods_detail"   description:"goods coupon detail list"`
}

// -------- For Agent Input

// WxDrH5 Request input data of H5 direct pay
type WxDrH5 struct {
	AppID      string          `json:"appid"        validate:"required" description:"official account of wechat"`
	MchID      string          `json:"mchid"        validate:"required" description:"merchant id of wechat"`
	Desc       string          `json:"description"  validate:"required" description:"goods description"`
	TradeNo    string          `json:"out_trade_no" validate:"required" description:"trade number of service provider system"`
	TimeExpire string          `json:"time_expire,omitempty"            description:"ticket expire time as unix seconds string"`
	Attach     string          `json:"attach,omitempty"                 description:"attach information"`
	NotifyURL  string          `json:"notify_url"   validate:"required" description:"result notify url send from wechat pay platform"`
	GoodsTag   string          `json:"goods_tag,omitempty"              description:"goods order discount mark"`
	Amount     *Amount         `json:"amount"       validate:"required" description:"trade amount information"`
	Detail     *DiscountDetail `json:"detail,omitempty"                 description:"promotion detail"`
	Scene      *SceneInfo      `json:"scene_info"   validate:"required" description:"trade scene"`
	Settle     *SettleInfo     `json:"settle_info,omitempty"            description:"settlement information"`
}

// WxDrApp Request input data of app direct pay
type WxDrApp struct {
	AppID      string          `json:"appid"        validate:"required" description:"official account of wechat"`
	MchID      string          `json:"mchid"        validate:"required" description:"merchant id of wechat"`
	Desc       string          `json:"description"  validate:"required" description:"goods description"`
	TradeNo    string          `json:"out_trade_no" validate:"required" description:"trade number of service provider system"`
	TimeExpire string          `json:"time_expire,omitempty"            description:"ticket expire time as unix seconds string"`
	Attach     string          `json:"attach,omitempty"                 description:"attach information"`
	NotifyURL  string          `json:"notify_url"   validate:"required" description:"result notify url send from wechat pay platform"`
	GoodsTag   string          `json:"goods_tag,omitempty"              description:"goods order discount mark"`
	Amount     *Amount         `json:"amount"       validate:"required" description:"trade amount information"`
	Detail     *DiscountDetail `json:"detail,omitempty"                 description:"promotion detail"`
	Scene      *SceneInfo      `json:"scene_info,omitempty"             description:"trade scene"`
	Settle     *SettleInfo     `json:"settle_info,omitempty"            description:"settlement information"`
}

// WxDrJS Request input data of JSAPI direct pay
type WxDrJS struct {
	AppID      string          `json:"appid"        validate:"required" description:"official account of wechat"`
	MchID      string          `json:"mchid"        validate:"required" description:"merchant id of wechat"`
	Desc       string          `json:"description"  validate:"required" description:"goods description"`
	TradeNo    string          `json:"out_trade_no" validate:"required" description:"trade number of service provider system"`
	TimeExpire string          `json:"time_expire,omitempty"            description:"ticket expire time as unix seconds string"`
	Attach     string          `json:"attach,omitempty"                 description:"attach information"`
	NotifyURL  string          `json:"notify_url"   validate:"required" description:"result notify url send from wechat pay platform"`
	GoodsTag   string          `json:"goods_tag,omitempty"              description:"goods order discount mark"`
	Amount     *Amount         `json:"amount"       validate:"required" description:"trade amount information"`
	Payer      *Payer          `json:"payer"        validate:"required" description:"trade payer"`
	Detail     *DiscountDetail `json:"detail,omitempty"                 description:"promotion detail"`
	Scene      *SceneInfo      `json:"scene_info,omitempty"             description:"trade scene"`
	Settle     *SettleInfo     `json:"settle_info,omitempty"            description:"settlement information"`
}

// WxDrRefund Request input data of direct refund
type WxDrRefund struct {
	TranID    string         `json:"transaction_id"                    description:"transaction id of wechat pay platform"`
	TradeNo   string         `json:"out_trade_no"  validate:"required" description:"trade number of service provider system"`
	RefundNo  string         `json:"out_refund_no" validate:"required" description:"refund transaction number of service provider system"`
	Reason    string         `json:"reason,omitempty"                  description:"refund action reason"`
	NotifyURL string         `json:"notify_url,omitempty"              description:"refund notify url send from wechat pay platform"`
	FundsAcc  string         `json:"funds_account,omitempty"           description:"refund funds account, such as : AVAILABLE"`
	Amount    *RefundAmount  `json:"amount"        validate:"required" description:"amount settle information"`
	Goods     []*RefundGoods `json:"goods_detail,omitempty"            description:"goods details of refund"`
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

// WxMchCerts wechat merchant certificates
type WxMchCerts struct {
	Datas []Certificate `json:"data"`
}

// WxRetUpload uploaded media id of wechat pay platform
type WxRetUpload struct {
	MediaID string `json:"media_id" description:"uploaded media id of wechat pay platform"`
}

// WxRetRefund refund result of trade informations
type WxRetRefund struct {
	RefundID    string          `json:"refund_id"             description:"* refund transaction id of wechat pay platform"`
	RefundNo    string          `json:"out_refund_no"         description:"* refund transaction number of service provider system"`
	TranID      string          `json:"transaction_id"        description:"* original transaction id of wechat pay platform"`
	TradeNo     string          `json:"out_trade_no"          description:"* orifinal transaction number of service provider system"`
	Channel     string          `json:"channel"               description:"* refund channel, such as : ORIGINAL, BALANCE, OTHER_BALANCE, OTHER_BANKCARD"`
	PayerAcc    string          `json:"user_received_account" description:"* payer received money account"`
	SuccessTime string          `json:"success_time"          description:"success refund time, format as YYYY-MM-DDTHH:mm:ss+TIMEZONE"`
	CreateTime  string          `json:"create_time"           description:"* create refund time, format as YYYY-MM-DDTHH:mm:ss+TIMEZONE"`
	RefundState string          `json:"status"                description:"* refund status, such as : SUCCESS, CLOSED, PROCESSING, ABNORMAL"`
	FundsAcc    string          `json:"funds_account"         description:"refund funds account, such as : AVAILABLE"`
	Amount      *RefRetAmount   `json:"amount"                description:"* amount settle information of refund"`
	Promotions  []*RefRetPromot `json:"promotion_detail"      description:"promotion details of refund"`
}
