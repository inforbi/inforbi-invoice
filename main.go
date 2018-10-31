package main

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
	"os"
)

func main() {
	// Create application
	app := widgets.NewQApplication(len(os.Args), os.Args)
	/*client := data.Client{Name: "Test", City: "Munich", Contact: "Frau Muster", Street: "Street"}
	client.SetFile("testClient.json")
	err := client.EncodeClient()
	if err != nil {
		fmt.Println(err)
	}
	invoice := data.Invoice{Project: "TestP", Number: 1, DueDays: 14, Items: []data.Item{{
		Date:        "12.2017",
		Description: "test",
		SinglePrice: 8.8,
		Quantifier:  "Hour",
		Quantity:    1,
	}}, Date: time.Now().Format("2006-01-02")}
	invoice.SetFile("testInvoice.json")
	fmt.Println(invoice)
	err = invoice.EncodeInvoice()
	if err != nil {
		fmt.Println(err)
	}*/

	core.QCoreApplication_SetApplicationName("Invoice Creator")

	mw := initMainWindow()

	// Show the window
	mw.Show()

	// Execute app
	app.Exec()

}
