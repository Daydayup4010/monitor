package services

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
	"uu/config"
	"uu/utils"
)

// YunGouOS API 基础地址
const YunGouOSBaseURL = "https://api.pay.yungouos.com"

// Native支付路径
const NativePayPath = "/api/pay/wxpay/nativePay"

// 支付请求客户端
var payClient = utils.CreateClient(YunGouOSBaseURL)

// NativePayResponse Native支付响应
type NativePayResponse struct {
	Code int    `json:"code"` // 0：成功，其他：失败
	Msg  string `json:"msg"`  // 返回消息
	Data string `json:"data"` // type=2时返回二维码base64图片
}

// PayGenerateSign 生成签名
// 根据 https://open.pay.yungouos.com/#/api/sign
// 1. 将所有参数按参数名ASCII码从小到大排序
// 2. 拼接成 key1=value1&key2=value2... 的格式
// 3. 在最后拼接 &key=商户密钥
// 4. 对拼接后的字符串进行MD5加密
// 5. 将MD5结果转换为大写
func PayGenerateSign(params map[string]string, apiKey string) string {
	// 1. 参数按ASCII码排序
	keys := make([]string, 0, len(params))
	for k := range params {
		if params[k] != "" { // 空值不参与签名
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	// 2. 拼接参数
	var buf strings.Builder
	for i, k := range keys {
		if i > 0 {
			buf.WriteString("&")
		}
		buf.WriteString(k)
		buf.WriteString("=")
		buf.WriteString(params[k])
	}

	// 3. 拼接密钥
	buf.WriteString("&key=")
	buf.WriteString(apiKey)

	// 4. MD5加密并转大写
	hash := md5.Sum([]byte(buf.String()))
	return strings.ToUpper(hex.EncodeToString(hash[:]))
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
	if attach != "" {
		signParams["attach"] = attach
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
