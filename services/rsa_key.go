package services

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

type RsaKey struct {
	Public  string `json:"public"`
	Private string `json:"private"`
	AppKey  string `json:"appKey"`
}

// SignParams 签名参数结构体
type SignParams struct {
	Timestamp string `json:"timestamp"`
	Sign      string `json:"sign"`
	AppKey    string `json:"appKey"`
}

// GenerateRSAKeys 生成RSA密钥对
func GenerateRSAKeys(bits int) (privateKeyPEM, publicKeyPEM string, err error) {
	// 1. 生成私钥
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return "", "", err
	}

	// 2. 编码私钥为PEM格式
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}
	privateKeyPEM = string(pem.EncodeToMemory(privateKeyBlock))

	// 3. 编码公钥为PEM格式
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return "", "", err
	}
	publicKeyBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	}
	publicKeyPEM = string(pem.EncodeToMemory(publicKeyBlock))

	return privateKeyPEM, publicKeyPEM, nil
}

func LoadRsaKey() *RsaKey {
	data, _ := os.ReadFile("rsa_key.json")
	var rsaKey RsaKey
	_ = json.Unmarshal(data, &rsaKey)

	return &rsaKey
}

// LoadPrivateKeyFromBase64 从Base64格式加载私钥（rsa_key.json中的格式）
func LoadPrivateKeyFromBase64(privateKeyBase64 string) (*rsa.PrivateKey, error) {
	// 移除 Base64 字符串中的换行符
	privateKeyBase64 = strings.ReplaceAll(privateKeyBase64, "\n", "")
	privateKeyBase64 = strings.ReplaceAll(privateKeyBase64, "\r", "")

	// 解码Base64
	derBytes, err := base64.StdEncoding.DecodeString(privateKeyBase64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 private key: %v", err)
	}

	// 尝试解析PKCS1格式
	privKey, err := x509.ParsePKCS1PrivateKey(derBytes)
	if err == nil {
		return privKey, nil
	}

	// 尝试解析PKCS8格式
	key, err := x509.ParsePKCS8PrivateKey(derBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %v", err)
	}

	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("private key is not RSA type")
	}

	return rsaKey, nil
}

// GetTimestamp 获取当前时间戳（格式：2006-01-02 15:04:05）
func GetTimestamp() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// GenerateSign 生成签名
// body: 请求体参数（map格式）
// timestamp: 时间戳
// appKey: 应用密钥
// privateKey: RSA私钥
func GenerateSign(body map[string]interface{}, timestamp string, appKey string, privateKey *rsa.PrivateKey) (string, error) {
	// 1. 构建待签名字符串
	signStr := buildSignString(body, timestamp, appKey)

	// 2. 计算SHA256哈希
	hash := sha256.Sum256([]byte(signStr))

	// 3. 使用私钥进行RSA签名
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hash[:])
	if err != nil {
		return "", fmt.Errorf("failed to sign: %v", err)
	}

	// 4. Base64编码
	return base64.StdEncoding.EncodeToString(signature), nil
}

// buildSignString 构建待签名字符串
// 按照悠悠有品的签名规则：将参数按字母顺序排序，拼接成 key + json.dumps(value) 的形式
func buildSignString(body map[string]interface{}, timestamp string, appKey string) string {
	// 复制body参数，添加timestamp和appKey
	params := make(map[string]interface{})
	for k, v := range body {
		params[k] = v
	}
	params["timestamp"] = timestamp
	params["appKey"] = appKey

	// 获取所有key并排序
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 构建签名字符串：key + json.dumps(value)，无分隔符
	var builder strings.Builder
	for _, k := range keys {
		builder.WriteString(k)
		// 使用 json.Marshal 将 value 转换为 JSON 字符串（紧凑格式）
		// 注意：需要禁用 HTML 转义，否则 & 会被转义为 \u0026
		jsonBytes := jsonMarshalNoEscape(params[k])
		builder.Write(jsonBytes)
	}

	return builder.String()
}

// jsonMarshalNoEscape JSON序列化，不转义HTML特殊字符（&, <, >）
func jsonMarshalNoEscape(v interface{}) []byte {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	_ = encoder.Encode(v)
	// Encode 会在末尾添加换行符，需要去掉
	result := buf.Bytes()
	if len(result) > 0 && result[len(result)-1] == '\n' {
		result = result[:len(result)-1]
	}
	return result
}

// GetSignParams 获取签名参数（timestamp和sign）
// body: 请求体参数
// 返回 SignParams 包含 timestamp, sign, appKey
func GetSignParams(body map[string]interface{}) (*SignParams, error) {
	// 加载RSA密钥
	rsaKey := LoadRsaKey()

	// 从Base64加载私钥
	privateKey, err := LoadPrivateKeyFromBase64(rsaKey.Private)
	if err != nil {
		return nil, fmt.Errorf("failed to load private key: %v", err)
	}

	// 获取时间戳
	timestamp := GetTimestamp()

	// 生成签名
	sign, err := GenerateSign(body, timestamp, rsaKey.AppKey, privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to generate sign: %v", err)
	}

	return &SignParams{
		Timestamp: timestamp,
		Sign:      sign,
		AppKey:    rsaKey.AppKey,
	}, nil
}

// GetSignParamsForJSON 获取JSON格式请求体的签名参数
// jsonBody: JSON格式的请求体
func GetSignParamsForJSON(jsonBody []byte) (*SignParams, error) {
	var body map[string]interface{}
	if err := json.Unmarshal(jsonBody, &body); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json body: %v", err)
	}
	return GetSignParams(body)
}
