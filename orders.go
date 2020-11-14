package main

import (
	"bytes"
	//	"crypto/tls"
	"encoding/json"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
)

const (
	EmailTemplate = `
<html>
	<h3> {{.Title}} </h3>
	<p> Board Size: {{.Size}} </p>

	{{range $i, $p := .PeopleOptions}}
		<p> Snowperson #{{inc $i }}: {{$p.Name}} <br/>
		&nbsp;&nbsp; Cap Color: {{$p.CapColor}} <br />
		&nbsp;&nbsp; Brim Color: {{$p.BrimColor}} <br />
		&nbsp;&nbsp; Pom Color: {{$p.PomColor}} <br />
		</p>
	{{end}}
</html>
	`
)

var (
	fromAddress = "woodnthings@example.net"
	toAddress   = []string{"willyingling@gmail.com"}
	mailServer  = "smtp.gmail.com"
	mailPort    = "587"

	orderCt      = 0
	mailCredPath = "$HOME/.config/emailAuth.json"
)

var (
	credentials mailCredentials
)

type mailCredentials struct {
	Email    string
	Password string
}

func (m *mailCredentials) readFromStdIn() error {
	fmt.Printf("Email: ")
	fmt.Fscanf(os.Stdin, "%s", &m.Email)
	fmt.Printf("Password: ")
	pass, err := terminal.ReadPassword(0)
	fmt.Printf("\n")

	if err != nil {
		return err
	}

	m.Email = strings.TrimSpace(m.Email)
	m.Password = strings.TrimSpace(string(pass))
	return nil
}

func (m *mailCredentials) readFromFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	dec := json.NewDecoder(file)
	return dec.Decode(m)
}

func readMailCredentials(interactive bool) error {
	if interactive {
		return credentials.readFromStdIn()
	}
	return credentials.readFromFile(os.ExpandEnv(mailCredPath))
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

func (bo boardOrder) String() string {
	funcMap := template.FuncMap{
		"inc": func(i int) int {
			return i + 1
		},
	}
	htmlOut := template.Must(template.New("Email").Funcs(funcMap).Parse(EmailTemplate))

	var b []byte
	buffer := bytes.NewBuffer(b)
	htmlOut.Execute(buffer, bo)

	return buffer.String()
}

type OrderStatus struct {
	Msg      string
	OrderNum int
}

func placeOrder(bo boardOrder) OrderStatus {
	orderCt++
	fmt.Printf("Received order: %d", orderCt)

	status := "failed"
	err := mailNotification(bo)
	if err == nil {
		status = "Success"
	} else {
		fmt.Printf("Error submitting mail: %s\n", err)
	}
	return OrderStatus{Msg: status, OrderNum: orderCt}
}

func mailNotification(bo boardOrder) error {
	// Exit early if there's no address
	if credentials.Email == "" {
		return nil
	}

	auth := smtp.PlainAuth("",
		credentials.Email,
		credentials.Password,
		mailServer)

	err := smtp.SendMail(mailServer+":"+mailPort,
		auth,
		fromAddress,
		toAddress,
		[]byte(getMailHeaders()+"\r\n"+bo.String()),
	)

	return err
}

func getMailHeaders() string {
	ret := ""
	endline := "\r\n"
	ret += "Content-Type: text/html" + endline
	ret += "From: " + fromAddress + endline
	ret += fmt.Sprintf("Subject: Order #%d"+endline, orderCt)
	return ret
}
