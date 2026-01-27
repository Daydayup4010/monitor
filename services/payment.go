package services

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"uu/config"
	"uu/utils"
)

// YunGouOS API 基础地址
const YunGouOSBaseURL = "https://api.pay.yungouos.com"

// Native支付路径
const NativePayPath = "/api/pay/wxpay/nativePay"

// 小程序支付路径
const MinAppPayPath = "/api/pay/wxpay/minAppPay"

// 支付请求客户端
var payClient = utils.CreateClient(YunGouOSBaseURL)

// NativePayResponse Native支付响应
type NativePayResponse struct {
	Code int    `json:"code"` // 0：成功，其他：失败
	Msg  string `json:"msg"`  // 返回消息
	Data string `json:"data"` // type=2时返回二维码base64图片
}

// MinAppPayData 小程序支付数据
type MinAppPayData struct {
	TimeStamp string `json:"timeStamp"` // 时间戳
	NonceStr  string `json:"nonceStr"`  // 随机字符串
	Package   string `json:"package"`   // prepay_id=xxx
	SignType  string `json:"signType"`  // 签名类型
	PaySign   string `json:"paySign"`   // 签名
}

// MinAppPayResponse 小程序支付响应
type MinAppPayResponse struct {
	Code int           `json:"code"` // 0：成功，其他：失败
	Msg  string        `json:"msg"`  // 返回消息
	Data MinAppPayData `json:"data"` // 支付参数
}

// PayGenerateSign 生成签名
// 根据 https://open.pay.yungouos.com/#/api/sign
// 1. 将所有参数按参数名ASCII码从小到大排序
// 2. 拼接成 key1=value1&key2=value2... 的格式
// 3. 在最后拼接 &key=商户密钥
// 4. 对拼接后的字符串进行MD5加密
// 5. 将MD5结果转换为大写
func PayGenerateSign(order map[string]string, key string) string {
	data := url.Values{}
	for k, v := range order {
		data.Add(k, v)
	}
	keys := make([]string, 0, 0)
	for key := range data {
		if data.Get(key) != "" {
			keys = append(keys, key)
		}
	}
	sort.Strings(keys)
	body := data.Encode()
	d, _ := url.QueryUnescape(body)
	d += "&key=" + key
	h := md5.New()
	h.Write([]byte(d))
	return strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
}

// VerifySign 验证回调签名
func VerifySign(params map[string]string, sign, apiKey string) bool {
	return PayGenerateSign(params, apiKey) == sign
}

// CreateNativePay 发起Native支付
// outTradeNo: 商户订单号
// totalFee: 支付金额（元）
// body: 商品描述
// attach: 附加数据，回调时原样返回
func CreateNativePay(outTradeNo string, totalFee float64, body, attach string) (*NativePayResponse, error) {
	paymentConfig := config.CONFIG.Payment
	if paymentConfig == nil {
		return nil, fmt.Errorf("payment config not found")
	}

	// 参与签名的参数（根据文档：out_trade_no, total_fee, mch_id, body, attach 参与签名）
	signParams := map[string]string{
		"out_trade_no": outTradeNo,
		"total_fee":    fmt.Sprintf("%.2f", totalFee),
		"mch_id":       paymentConfig.MchId,
		"body":         body,
	}

	// 生成签名
	sign := PayGenerateSign(signParams, paymentConfig.ApiKey)

	// 完整请求参数（包含不参与签名的参数：type, notify_url）
	requestParams := map[string]string{
		"out_trade_no": outTradeNo,
		"total_fee":    fmt.Sprintf("%.2f", totalFee),
		"mch_id":       paymentConfig.MchId,
		"body":         body,
		"type":         "2",                     // type=2 返回二维码base64图片（不参与签名）
		"notify_url":   paymentConfig.NotifyUrl, // 不参与签名
		"sign":         sign,
		"attach":       attach,
	}
	if attach != "" {
		requestParams["attach"] = attach
	}

	config.Log.Infof("NativePay request: out_trade_no=%s, total_fee=%.2f, body=%s", outTradeNo, totalFee, body)

	var result NativePayResponse
	resp, err := payClient.DoRequest("POST", NativePayPath, utils.RequestOptions{
		FormData: requestParams,
		Result:   &result,
	})

	if err != nil {
		config.Log.Errorf("NativePay request failed: %v", err)
		return nil, err
	}

	config.Log.Infof("NativePay response: %s", resp.String())

	return &result, nil
}

// CreateMinAppPay 发起小程序支付
// outTradeNo: 商户订单号
// totalFee: 支付金额（元）
// body: 商品描述
// openId: 用户openid
// attach: 附加数据，回调时原样返回
func CreateMinAppPay(outTradeNo string, totalFee float64, body, openId, attach string) (*MinAppPayResponse, error) {
	paymentConfig := config.CONFIG.Payment
	if paymentConfig == nil {
		return nil, fmt.Errorf("payment config not found")
	}

	// 获取小程序AppID
	wechatConfig := config.CONFIG.Wechat
	if wechatConfig == nil {
		return nil, fmt.Errorf("wechat config not found")
	}

	// 参与签名的参数（根据YunGouOS文档：out_trade_no, total_fee, mch_id, body, openId, app_id 参与签名）
	// 注意：attach, notify_url 不参与签名
	signParams := map[string]string{
		"out_trade_no": outTradeNo,
		"total_fee":    fmt.Sprintf("%.2f", totalFee),
		"mch_id":       paymentConfig.MchId,
		"body":         body,
		"open_id":      openId,
		"app_id":       wechatConfig.AppID,
	}

	// 生成签名
	sign := PayGenerateSign(signParams, paymentConfig.ApiKey)

	// 完整请求参数
	requestParams := map[string]string{
		"out_trade_no": outTradeNo,
		"total_fee":    fmt.Sprintf("%.2f", totalFee),
		"mch_id":       paymentConfig.MchId,
		"body":         body,
		"open_id":      openId,
		"app_id":       wechatConfig.AppID,
		"notify_url":   paymentConfig.NotifyUrl,
		"sign":         sign,
	}
	if attach != "" {
		requestParams["attach"] = attach
	}

	config.Log.Infof("MinAppPay request: out_trade_no=%s, total_fee=%.2f, body=%s, openId=%s", outTradeNo, totalFee, body, openId)

	var result MinAppPayResponse
	resp, err := payClient.DoRequest("POST", MinAppPayPath, utils.RequestOptions{
		FormData: requestParams,
		Result:   &result,
	})

	if err != nil {
		config.Log.Errorf("MinAppPay request failed: %v", err)
		return nil, err
	}

	config.Log.Infof("MinAppPay response: %s", resp.String())

	return &result, nil
}
