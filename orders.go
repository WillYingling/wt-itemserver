package main

import (
	"fmt"
	"net/smtp"
	"os"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
)

var (
	fromAddress = "woodnthings@example.net"
	toAddress   = []string{"willyingling@gmail.com"}
	mailServer  = "smtp.gmail.com"
	mailPort    = "587"
)

var (
	notificationEmail    string
	notificationPassword string
)

func readMailCredentials() error {
	fmt.Printf("Email: ")
	fmt.Fscanf(os.Stdin, "%s", &notificationEmail)
	fmt.Printf("Password: ")
	pass, err := terminal.ReadPassword(0)
	fmt.Printf("\n")

	if err != nil {
		return err
	}

	notificationEmail = strings.TrimSpace(notificationEmail)
	notificationPassword = strings.TrimSpace(string(pass))
	return nil
}

type personOptions struct {
	Name      string
	HatType   string
	CapColor  string
	BrimColor string
	PomColor  string
}

type boardOrder struct {
	Title         string
	Size          int
	PeopleOptions []personOptions
}

type OrderStatus struct {
	Msg string
}

func placeOrder() OrderStatus {
	return OrderStatus{Msg: "Succuess"}
}

func mailNotification() error {
	fmt.Printf("Attempting to send mail\n")
	auth := smtp.PlainAuth("",
		notificationEmail,
		notificationPassword,
		mailServer)

	msg := "Here is a message"

	err := smtp.SendMail(mailServer+":"+mailPort,
		auth,
		fromAddress,
		toAddress,
		[]byte(msg))

	if err != nil {
		fmt.Printf("Error sending message: %s\n", err)
	}

	fmt.Printf("Sent mail\n")
	return nil
}
