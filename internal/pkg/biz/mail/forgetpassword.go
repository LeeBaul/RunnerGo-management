package mail

import (
	"context"
	"fmt"
	"net/smtp"

	"kp-management/internal/pkg/conf"
)

const (
	forgetHTMLTemplate = `<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Document</title>
    <style>
        * {
            padding: 0;
            margin: 0;
        }
        .email {
            width: 100vw;
            height: 100vh;
            background-color: #f2f2f2;
            display: flex;
            flex-direction: column;
            align-items: center;
            justify-content: center;
        }

        .logo {
            width: 159px;
            height: 26px;
        }

        .title {
            margin-top: 30px;
            font-size: 30px;
            color: #000;
            font-family: 'PingFang SC';
            font-style: normal;
            font-weight: 600;
        }

        .slogn {
            margin-top: 30px;
            font-family: 'PingFang SC';
            font-style: normal;
            font-weight: 400;
            font-size: 18px;
            color: #999999;
        }

        .email-body {
            width: 386px;
            height: 135px;
            background-color: #f8f8f8;
            border-radius: 15px;
            display: flex;
            flex-direction: column;
            align-items: center;
            margin-top: 77px;
            padding-top: 24px;
            box-sizing: border-box;
        }

        .email-body > p {
            font-family: 'PingFang SC';
            font-style: normal;
            font-weight: 600;
            font-size: 16px;
            color: #000;
        }

        .email-body > button {
            background: #EC663C;
            border-radius: 5px;
            width: 335px;
            height: 41px;
            color: #fff;
            margin-top: 24px;
            border: none;
        }
    </style>
</head>

<body>
    <div class="email">
        <img class="logo" src="https://apipost.oss-cn-beijing.aliyuncs.com/kunpeng/logo_black.png" alt="">
        <p class="title">性能测试平台</p>
        <p class="slogn">预见未来, 轻松上线</p>
        <div class="email-body">
            <p>点击下方按钮重置密码</p>
            <button><a href="%s">重置密码</a></button>
        </div>
    </div>
</body>

</html>`
)

func SendForgetPasswordEmail(ctx context.Context, toEmail string, userID int64) error {
	host := conf.Conf.SMTP.Host
	port := conf.Conf.SMTP.Port
	email := conf.Conf.SMTP.Email
	password := conf.Conf.SMTP.Password

	header := make(map[string]string)
	header["From"] = "RunnerGo" + "<" + email + ">"
	header["To"] = toEmail
	header["Subject"] = "找回密码"
	header["Content-Type"] = "text/html; charset=UTF-8"
	body := fmt.Sprintf(forgetHTMLTemplate, conf.Conf.Base.Domain+"/reset?email="+toEmail)
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
