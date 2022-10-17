package mail

import (
	"context"
	"fmt"
	"net/smtp"

	"kp-management/internal/pkg/conf"
)

func SendReportEmail(ctx context.Context, toEmail string) error {
	host := conf.Conf.SMTP.Host
	port := conf.Conf.SMTP.Port
	email := conf.Conf.SMTP.Email
	password := conf.Conf.SMTP.Password

	header := make(map[string]string)
	header["From"] = "test" + "<" + email + ">"
	header["To"] = toEmail
	header["Subject"] = "报告通知"
	header["Content-Type"] = "text/html; charset=UTF-8"
	body := "我是一封测试电子邮件!"
	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body
	auth := smtp.PlainAuth(
		"",
		email,
		password,
		host,
	)
	return SendMailUsingTLS(
		fmt.Sprintf("%s:%d", host, port),
		auth,
		email,
		[]string{toEmail},
		[]byte(message),
	)
}
