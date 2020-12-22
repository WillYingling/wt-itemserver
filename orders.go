package main

import (
	"bytes"
	"fmt"
	"html/template"
)

const (
	EmailTemplate = `
<html>
	<h3> {{.Board.Title}} </h3>
	<p> Board Size: {{.Board.Size}} </p>

{{range $i, $p := .Board.PeopleOptions}}
	<p> Snowperson #{{inc $i }}: {{$p.Name}} <br/>
	&nbsp;&nbsp; Cap Color: {{$p.CapColor}} <br />
	&nbsp;&nbsp; Brim Color: {{$p.BrimColor}} <br />
	&nbsp;&nbsp; Pom Color: {{$p.PomColor}} <br />
	</p>
{{end}}

{{range $i, $e := .Board.Extras}}
	<p> Extra #{{inc $i }}: {{$e.Type}} <br/>
	&nbsp;&nbsp; Name: {{$e.Name}} <br />
	&nbsp;&nbsp; Notes: {{$e.Notes}} <br />
	</p>
{{end}}

	<p>
		<u> Contact Info </u> <br />
		Name: {{.Info.Name}} <br />
		Preference: {{.Info.Preference}} <br />
		Email: {{.Info.Email}} <br />
		Phone Number: {{.Info.PhoneNum}} <br />
	</p>

	<a href="http://{{.NotifyUrl}}" > Notify this customer </a>

</html>
	`
)

var (
	orderCt = 0

	orderList []boardOrder
)

type contactInfo struct {
	Name       string
	Preference string
	Email      string
	PhoneNum   string
}

type boardOrder struct {
	Board     snowmanBoard
	NotifyUrl string
	Info      contactInfo
	id        int
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
	fmt.Printf("Received order: %d", orderCt)

	bo.id = orderCt
	// validate the order
	bo.NotifyUrl += fmt.Sprintf("?id=%d", bo.id)
	orderList = append(orderList, bo)

	status := "failed"

	err := mailNotification(bo)
	if err == nil {
		status = "Success"
		orderCt++
	} else {
		fmt.Printf("Error submitting mail: %s\n", err)
	}
	return OrderStatus{Msg: status, OrderNum: orderCt}
}

func completeOrder(id int) error {
	var bo boardOrder
	found := false
	for _, order := range orderList {
		if order.id == id {
			bo = order
			found = true
		}
	}

	if found == false {
		return fmt.Errorf("Order %d does not exist", id)
	}
	return notifyCustomer(bo)
}
