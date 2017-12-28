package main

import (
	"fmt"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
	"io/ioutil"
	"os"
)

func main() {
	// Create application
	app := widgets.NewQApplication(len(os.Args), os.Args)
	client := Client{Name: "Test", City: "Munich", Contact: "Frau Muster", Street: "Street", file: "testClient.json"}
	err := client.EncodeClient()
	if err != nil {
		fmt.Println(err)
	}
	invoice := Invoice{Project: "TestP", Number: 1, DueDays: 14, Items: []Item{{
		Date:        "12.2017",
		Description: "test",
		SinglePrice: 8.8,
		Quantifier:  "Hour",
		Quantity:    1,
	}}, file: "testInvoice.json"}
	fmt.Println(invoice)
	err = invoice.EncodeInvoice()
	if err != nil {
		fmt.Println(err)
	}

	cont, err := ioutil.ReadFile("invoice.pylat")
	if err != nil {
		fmt.Println(err)
	}
	template := string(cont)
	template = client.ReplaceTemplate(template)
	template = invoice.ReplaceTemplate(template)
	ioutil.WriteFile("invoice.tex", []byte(template), 0644)

	core.QCoreApplication_SetApplicationName("Invoice Creator")

	mw := initMainWindow()

	// Show the window
	mw.Show()

	// Execute app
	app.Exec()

}
