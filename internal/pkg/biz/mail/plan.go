package mail

import (
	"context"
	"fmt"
	"net/smtp"

	"kp-management/internal/pkg/conf"
	"kp-management/internal/pkg/dal/model"
)

const (
	planHTMLTemplate = `<!DOCTYPE html>
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
            width: 908px;
            /* height: 135px; */
            background-color: #f8f8f8;
            border-radius: 15px;
            display: flex;
            flex-direction: column;
            align-items: center;
            margin-top: 36px;
            padding-top: 24px;
            box-sizing: border-box;
            /* overflow: hidden; */
        }

        .email-body > .plan-name {
            font-family: 'PingFang SC';
            font-style: normal;
            font-weight: 600;
            font-size: 16px;
            color: #000;
        }

        .report-list {
            display: flex;
            flex-direction: column;
            align-items: center;
            justify-content: center;
            position: relative;
            margin-top: 13px;
            padding-bottom: 41px;
        }

        .line {
            width: 817px;
            height: 10px;
            background-color: #1A1A1D;
            border-radius: 4.5px;
            position: absolute;
            top: 15px;
        }

        .list {
            display: flex;
            flex-direction: column;
            padding-bottom: 20px;
            margin-top: 20px;
            background-color: #FEFEFE;
            border-width: 0px 2px 2px 2px;
            border-style: solid;
            border-color: #EC663C;
            width: 805px;
            max-height: 400px;
            overflow-y: scroll;
            box-sizing: border-box;
            z-index: 20;
        }

        .list-item {
            display: flex;
            box-sizing: border-box;
            justify-content: space-between;
            padding: 10px 0;
            font-size: 12px;
            margin: 0 26px;
            border-bottom: 1px solid #000;
            cursor: pointer;
        }

        .list-item > p:nth-child(1) {
            font-size: 14px !important;
        }

        .list-item:hover {
            color: #EC663C;
            border-color: #EC663C;
        }


        .team {
            font-size: 20px;
            margin-top: 36px;
        }

        a {
            text-decoration: none;
            color: #000;
        }
    </style>
</head>

<body>
    <div class="email">
        <img class="logo" src="https://apipost.oss-cn-beijing.aliyuncs.com/kunpeng/logo_black.png" alt="">
        <p class="title">性能测试平台</p>
        <p class="slogn">预见未来, 轻松上线</p>
        <p class="team">【%s】</p>
        <div class="email-body">
            <p class="plan-name">【%s】By %s</p>
            <div class="report-list">
                <div class="line"></div>
                <div class="list">
                    %s
                </div>
            </div>
        </div>
    </div>
</body>

</html>`

	reportListHTMLTemplate = `<a href="%s">
                        <div class="list-item">
                            <p>【%s】</p>
                            <p>执行者: %s</p>
                            <p>%s</p>
                        </div>
                    </a>`
)

func SendPlanEmail(ctx context.Context, toEmail string, planName, teamName, userName string, reports []*model.Report, runUsers []*model.User) error {
	host := conf.Conf.SMTP.Host
	port := conf.Conf.SMTP.Port
	email := conf.Conf.SMTP.Email
	password := conf.Conf.SMTP.Password

	header := make(map[string]string)
	header["From"] = "RunnerGo" + "<" + email + ">"
	header["To"] = toEmail
	header["Subject"] = fmt.Sprintf("测试报告 【%s】的【%s】给您发送了【%s】的测试报告，点击查看", teamName, userName, planName)
	header["Content-Type"] = "text/html; charset=UTF-8"

	memo := make(map[int64]*model.User, 0)
	for _, user := range runUsers {
		memo[user.ID] = user
	}

	var r string
	for _, report := range reports {

		r += fmt.Sprintf(reportListHTMLTemplate, conf.Conf.Base.Domain+"#/email/report?report_id="+fmt.Sprintf("%d", report.ID)+"&team_id="+fmt.Sprintf("%d", report.TeamID), report.SceneName, memo[report.RunUserID].Nickname, report.CreatedAt.Format("2006-01-02 15:04:05"))
	}
	body := fmt.Sprintf(planHTMLTemplate, teamName, planName, userName, r)
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
