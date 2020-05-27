package service

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

type MailConfig struct {
	Recipient string
	Ctx       context.Context
}

func (m *MailConfig) Send(flag string) {

	fmt.Println(m.Ctx.Value("mailSender"))

	mail := gomail.NewMessage()
	mail.SetHeader("From", fmt.Sprintf("%v", m.Ctx.Value("asa")))
	mail.SetHeader("To", m.Recipient)
	mail.SetHeader("Subject", "Already Come In MyCompany")
	mail.SetBody("text/html", "Hello, <b>you have alrady in my company </b>")

	dial := gomail.NewDialer(
		fmt.Sprintf("%v", m.Ctx.Value("mailHost")),
		m.Ctx.Value("mailPort").(int),
		fmt.Sprintf("%v", m.Ctx.Value("mailSender")),
		fmt.Sprintf("%v", m.Ctx.Value("mailPassword")),
	)

	err := dial.DialAndSend(mail)
	if err != nil {
		log.Error("Error DialAndSend, ", err.Error())
	}

	log.Info("Message Sent")

}
