package utils

import (
	"crypto/rand"
	"fmt"
	"gopkg.in/gomail.v2"
	"math/big"
	"uu/config"
)

func GenerateVerificationCode(length int) string {
	const digits = "0123456789"
	code := make([]byte, length)
	max := big.NewInt(int64(len(digits)))

	for i := range code {
		num, err := rand.Int(rand.Reader, max)
		if err != nil {
			return ""
		}
		code[i] = digits[num.Int64()]
	}
	return string(code)
}

type EmailService struct {
	SMTPHost     string `yaml:"host"`
	SMTPPort     int    `yaml:"port"`
	FromEmail    string `yaml:"email"`
	FromPassword string `yaml:"password"`
}

func (es *EmailService) SendVerificationCode(toEmail, code string) int {
	m := gomail.NewMessage()
	m.SetHeader("From", es.FromEmail)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", "【Monitor】")
	body := fmt.Sprintf(`<h2>您好！</h2>
        <p>您的邮箱验证码是：<strong>%s</strong></p>
        <p>请在10分钟内完成验证。</p>
        <p><small>如非本人操作，请忽略此邮件。</small></p>`, code)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(es.SMTPHost, es.SMTPPort, es.FromEmail, es.FromPassword)

	err := d.DialAndSend(m)
	if err != nil {
		config.Log.Errorf("send code fail: %v", err)
		return ErrCodeSendEmailCode
	}
	return SUCCESS
}
