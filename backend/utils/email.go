// /backend/utils/email.go
package utils

import (
	"crypto/tls"
	"fmt" // <-- 新增导入
	"log" // <-- 新增导入
	"net/smtp"
	"os"
	"strconv"

	"github.com/jordan-wright/email"
)

const (
	SMTP_HOST = "smtp.163.com"
	SMTP_PORT = 465
	SMTP_USER = "caigbro@163.com"
)

func SendVerificationCode(toEmail, code string) error {
	smtpPass := os.Getenv("SMTP_PASSWORD")
	if smtpPass == "" {
		log.Println("Error: SMTP_PASSWORD environment variable not set")
		return fmt.Errorf("email service is not configured")
	}

	e := email.NewEmail()
	e.From = "Your Blog <" + SMTP_USER + ">"
	e.To = []string{toEmail}
	e.Subject = "您的登录验证码"
	e.HTML = []byte("您好，您的登录验证码是：<b>" + code + "</b>。该验证码 5 分钟内有效。")

	addr := SMTP_HOST + ":" + strconv.Itoa(SMTP_PORT)

	return e.SendWithTLS(addr,
		smtp.PlainAuth("", SMTP_USER, smtpPass, SMTP_HOST),
		&tls.Config{ServerName: SMTP_HOST, InsecureSkipVerify: true},
	)
}