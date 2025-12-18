package services

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"uu/config"
)

// NativePayURL YunGouOS Native支付API地址
const NativePayURL = "https://api.pay.yungouos.com/api/pay/wxpay/nativePay"

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

// 验证回调签名
func VerifySign(params map[string]string, sign, apiKey string) bool {
	return PayGenerateSign(params, apiKey) == sign
}

// 发起Native支付
// outTradeNo: 商户订单号
// totalFee: 支付金额（元）
// body: 商品描述
// attach: 附加数据，回调时原样返回
func CreateNativePay(outTradeNo string, totalFee float64, body, attach string) (*NativePayResponse, error) {
	paymentConfig := config.CONFIG.Payment
	if paymentConfig == nil {
		return nil, fmt.Errorf("payment config not found")
	}

	params := map[string]string{
		"out_trade_no": outTradeNo,
		"total_fee":    fmt.Sprintf("%.2f", totalFee),
		"mch_id":       paymentConfig.MchId,
		"body":         body,
		"type":         "2", // type=2 返回二维码base64图片
		"notify_url":   paymentConfig.NotifyUrl,
	}
	if attach != "" {
		params["attach"] = attach
	}

	// 生成签名
	params["sign"] = PayGenerateSign(params, paymentConfig.ApiKey)

	// 发起POST请求
	formData := url.Values{}
	for k, v := range params {
		formData.Set(k, v)
	}

	config.Log.Infof("NativePay request: out_trade_no=%s, total_fee=%.2f, body=%s", outTradeNo, totalFee, body)

	resp, err := http.PostForm(NativePayURL, formData)
	if err != nil {
		config.Log.Errorf("NativePay request failed: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		config.Log.Errorf("NativePay read response failed: %v", err)
		return nil, err
	}

	config.Log.Infof("NativePay response: %s", string(respBody))

	var result NativePayResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		config.Log.Errorf("NativePay parse response failed: %v", err)
		return nil, err
	}

	return &result, nil
}
