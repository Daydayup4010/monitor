package utils

import (
	"crypto/rand"
	"fmt"
	"gopkg.in/gomail.v2"
	"math/big"
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
	m.SetHeader("Subject", "【CS Goods】")
	body := fmt.Sprintf(`<h2>您好！</h2>
        <p>您的邮箱验证码是：<strong>%s</strong></p>
        <p>请在10分钟内完成验证。</p>
        <p><small>如非本人操作，请忽略此邮件。</small></p>`, code)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(es.SMTPHost, es.SMTPPort, es.FromEmail, es.FromPassword)

	err := d.DialAndSend(m)
	if err != nil {
		return ErrCodeSendEmailCode
	}
	return SUCCESS
}

// SendVIPNotification 发送VIP开通/续费通知邮件
func (es *EmailService) SendVIPNotification(toEmail string, months int, expiryDate string) int {
	m := gomail.NewMessage()
	m.SetHeader("From", es.FromEmail)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", "【CS Goods】VIP会员开通成功")

	body := fmt.Sprintf(`
		<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px;">
			<h2 style="color: #1890ff;">🎉 恭喜您成为CS Goods VIP会员！</h2>
			<div style="background: #f5f5f5; padding: 20px; border-radius: 8px; margin: 20px 0;">
				<p style="margin: 10px 0;"><strong>会员时长：</strong>%d 个月</p>
				<p style="margin: 10px 0;"><strong>到期时间：</strong>%s</p>
			</div>
			<p>您现在可以享受以下VIP特权：</p>
			<ul style="color: #666;">
				<li>📊 完整饰品涨跌榜数据</li>
				<li>💰 搬砖利润分析工具</li>
				<li>📈 饰品走势图表</li>
				<li>🔔 更多专业功能</li>
			</ul>
			<div style="text-align: center; margin: 30px 0;">
				<a href="https://www.csgoods.com.cn" style="display: inline-block; background: #1890ff; color: #fff; padding: 12px 30px; border-radius: 6px; text-decoration: none; font-weight: 500;">立即访问 CS Goods</a>
			</div>
			<p style="text-align: center; color: #666;">网站地址：<a href="https://www.csgoods.com.cn" style="color: #1890ff;">www.csgoods.com.cn</a></p>
			<p style="margin-top: 20px;">感谢您的支持！如有任何问题，请联系我们。</p>
			<p style="color: #999; font-size: 12px; margin-top: 30px;">
				此邮件由系统自动发送，请勿直接回复。<br>
				Email: goods.monitor@foxmail.com | QQ: 401026211
			</p>
		</div>
	`, months, expiryDate)

	m.SetBody("text/html", body)

	d := gomail.NewDialer(es.SMTPHost, es.SMTPPort, es.FromEmail, es.FromPassword)

	err := d.DialAndSend(m)
	if err != nil {
		return ErrCodeSendEmailCode
	}
	return SUCCESS
}

// SendErrorAlert 发送错误告警邮件
func (es *EmailService) SendErrorAlert(recipients []string, subject, body string) error {
	if len(recipients) == 0 {
		return fmt.Errorf("no recipients specified")
	}

	m := gomail.NewMessage()
	m.SetHeader("From", es.FromEmail)
	m.SetHeader("To", recipients...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(es.SMTPHost, es.SMTPPort, es.FromEmail, es.FromPassword)

	return d.DialAndSend(m)
}
