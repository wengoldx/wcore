package wechat

import (
	"fmt"
)

type WxPayAgent struct {
}

const (
	WxpApisDomain = "https://api.mch.weixin.qq.com"

	// Wechat certificate and upload medias
	WxpApiCert  = "/v3/certificates"
	WxpApiUpImg = "/v3/merchant/media/upload"

	// Wechat Direct connected merchants
	WxpDrApp     = "/v3/pay/transactions/app"
	WxpDrJS      = "/v3/pay/transactions/jsapi"
	WxpDrNative  = "/v3/pay/transactions/native"
	WxpDrH5      = "/v3/pay/transactions/h5"
	WxpDrIDQuery = "/v3/pay/transactions/id/%s"
	WxpDrNoQuery = "/v3/pay/transactions/out-trade-no/%s"
	WxpDrClose   = "/v3/pay/transactions/out-trade-no/%s/close"

	WxpDrRefund   = "/secapi/pay/refund"
	WxpDrRefQuery = "/pay/refundquery"

	// Wechat SP, PF payments and merch account change
	WxpMchApp     = "/v3/pay/partner/transactions/app"
	WxpMchJS      = "/v3/pay/partner/transactions/jsapi"
	WxpMchH5      = "/v3/pay/partner/transactions/h5"
	WxpMchIDQuery = "/v3/pay/partner/transactions/id/%s"
	WxpMchNoQuery = "/v3/pay/partner/transactions/out-trade-no/%s"
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
	WxpPFWNoQuery = "/v3/ecommerce/fund/withdraw/out-request-no/%s"

	WxpPFDividing   = "/v3/ecommerce/profitsharing/orders"
	WxpPFDiviRefund = "/v3/ecommerce/profitsharing/returnorders"
	WxpPFDiviClose  = "/v3/ecommerce/profitsharing/finish-order"
	WxpPFDiviRecAdd = "/v3/ecommerce/profitsharing/receivers/add"
	WxpPFDiviRecDel = "/v3/ecommerce/profitsharing/receivers/delete"

	WxpPFRefund   = "/v3/ecommerce/refunds/apply"
	WxpPFRIDQuery = "/v3/ecommerce/refunds/id/%s"
	WxpPFRNoQuery = "/v3/ecommerce/refunds/out-refund-no/%s"
)

// Append wechat pay APIv3 domain and params combined string
//	@param formatpath "API path, it maybe have format keyword"
//	@param param      "Default dynamic params to insert key value into formatpath"
//	@return - string "Full url link with param if have"
func (w *WxPayAgent) Url(formatpath string, param ...string) string {
	path := formatpath
	if num := len(param); num > 0 {
		path = fmt.Sprintf(formatpath, param[0])
	}
	return WxpApisDomain + path
}
