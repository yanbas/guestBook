package service

import (
	"crypto/tls"

	"gopkg.in/gomail.v2"
)

func send() {
	d := gomail.NewDialer("smtp.example.com", 587, "user", "123456")
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

}
