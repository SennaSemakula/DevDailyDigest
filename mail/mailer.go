package main

import (
	mail "github.com/Pioneersltd/DevDailyDigest/v1/mail/email"
)

func main() {

	client := mail.NewClient()
	client.Authenticate()

	mail := &mail.Mail{Subject: "Test email", To: []string{"tst.@gmail.com"}, Body: nil}
	to := mail.To[0]

	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + mail.Subject + ":\r\n" +
		"\r\n" +
		string(mail.Body) + "\r\n")

	mail.Body = msg

	client.Send(mail)
}
