package main

import (
	"fmt"
	"net/smtp"

	"github.com/spf13/viper"
)

func notifyCustomer(bo boardOrder) error {
	msg := fmt.Sprintf("Your order \"%s\" is ready!", bo.Board.Title)
	if bo.Info.Preference == "Email" {
		return sendEmail("Your order is ready", msg, []string{bo.Info.Email})
	}

	return fmt.Errorf("Contact preference: %s not supported", bo.Info.Preference)
}

func mailNotification(bo boardOrder) error {
	// Exit early if there's no address
	if debug || !emailConfig.IsSet("Email") {
		verbosePrint("%s\n", bo.String())
		return nil
	}

	return sendEmail(fmt.Sprintf("Order #%d", bo.id),
		bo.String(), viper.GetStringSlice("AdminAddresses"))
}

func sendEmail(subject, body string, rec []string) error {

	server := emailConfig.GetString("MailServer")
	port := emailConfig.GetString("MailServerPort")

	auth := smtp.PlainAuth("",
		emailConfig.GetString("Email"),
		emailConfig.GetString("Password"),
		server,
	)

	err := smtp.SendMail(server+":"+port,
		auth,
		viper.GetString("ServerMailAddress"),
		rec,
		[]byte(getMailHeaders(subject)+"\r\n"+body),
	)

	return err
}

func getMailHeaders(subject string) string {
	ret := ""
	endline := "\r\n"
	ret += "Content-Type: text/html" + endline
	ret += "From: " + viper.GetString("ServerMailAddress") + endline
	ret += fmt.Sprintf("Subject: %s"+endline, subject)
	return ret
}
