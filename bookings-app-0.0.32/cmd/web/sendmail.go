package main

import (
	"log"

	"github.com/tsawler/bookings-app/internal/models"

	mail "github.com/xhit/go-simple-mail/v2"
)

func ListenForMail() {
	go func() {
		for {
			msg := <-app.MailChan
			SendMessage(msg)

		}
	}()

}

func SendMessage(m models.MailData) {
	server := mail.NewSMTPClient()
	server.Host = "localhost"
	server.Port = 1234
	server.KeepAlive = false
	// server.ConnectTimeOut = 10
	// server.SendTimeOut = 10

	client, err := server.Connect()
	if err != nil {
		errorLog.Println(err)
	}
	email := mail.NewMSG()
	email.SetFrom(m.From).AddTo(m.To).SetSubject(m.Subject)
	email.SetBody(mail.TextHTML, string(m.Content))
	err = email.Send(client)
	if err != nil {
		log.Println(err)
	}
	log.Println("email sent la")
}
