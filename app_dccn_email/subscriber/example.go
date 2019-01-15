package subscriber

import (
	"context"
	"log"

	mail "github.com/Ankr-network/dccn-hub/app_dccn_email/proto/v1"

	gomail "gopkg.in/gomail.v2"
)

var (
	sender = "994336359@qq.com"
	passwd = "espgstzviouubfjh"
)

func SendMail(e *mail.MailEvent) error {
	m := gomail.NewMessage()
	m.SetAddressHeader("From", sender, "Ankr Network")
	m.SetHeader("To",
		m.FormatAddress(e.Address, e.Name),
	)
	m.SetHeader("Subject", "Gomail")
	m.SetBody("text/html", e.Message)

	d := gomail.NewPlainDialer("smtp.qq.com", 465, sender, passwd) // 发送邮件服务器、端口、发件人账号、发件人密码
	if err := d.DialAndSend(m); err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func Handler(ctx context.Context, e *mail.MailEvent) error {
	log.Println("Function Received message: ", e)

	return SendMail(e)
}
