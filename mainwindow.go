package main

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
	"os"
	"strconv"
)

type MainWindow struct {
	widgets.QMainWindow

	clientName   *widgets.QLabel
	clientChoose *widgets.QPushButton
	clientEdit   *widgets.QPushButton
	clientCreate *widgets.QPushButton

	invoiceName   *widgets.QLabel
	invoiceChoose *widgets.QPushButton
	invoiceEdit   *widgets.QPushButton
	invoiceCreate *widgets.QPushButton

	previewBtn *widgets.QPushButton
	saveTexBtn *widgets.QPushButton
	savePdfBtn *widgets.QPushButton
}

var (
	clientSelected  = false
	invoiceSelected = false
	selectedClient  Client
	selectedInvoice Invoice
)

func initMainWindow() *MainWindow {
	this := NewMainWindow(nil, 0)
	this.SetWindowTitle(core.QCoreApplication_ApplicationName())
	upperWidget := widgets.NewQWidget(this, 0)
	lowerWidget := widgets.NewQWidget(nil, 0)
	upperGrid := widgets.NewQGridLayout(upperWidget)
	lowerGrid := widgets.NewQGridLayout(lowerWidget)
	upperGrid.SetSpacing(10)
	lowerGrid.SetSpacing(50)

	this.SetCentralWidget(upperWidget)

	upperGrid.AddWidget3(lowerWidget, 2, 0, 1, 4, core.Qt__AlignCenter)

	this.clientName = widgets.NewQLabel2("<No client selected>", nil, 0)
	this.clientChoose = widgets.NewQPushButton2("Choose...", nil)
	this.clientEdit = widgets.NewQPushButton2("Edit", nil)
	this.clientCreate = widgets.NewQPushButton2("Create", nil)

	this.clientChoose.ConnectPressed(this.chooseClient)

	upperGrid.AddWidget(this.clientName, 0, 0, core.Qt__AlignLeft)
	upperGrid.AddWidget(this.clientChoose, 0, 1, core.Qt__AlignLeft)
	upperGrid.AddWidget(this.clientEdit, 0, 2, core.Qt__AlignLeft)
	upperGrid.AddWidget(this.clientCreate, 0, 3, core.Qt__AlignLeft)

	this.invoiceName = widgets.NewQLabel2("<No invoice selected>", nil, 0)
	this.invoiceChoose = widgets.NewQPushButton2("Choose...", nil)
	this.invoiceEdit = widgets.NewQPushButton2("Edit", nil)
	this.invoiceCreate = widgets.NewQPushButton2("Create", nil)

	upperGrid.AddWidget(this.invoiceName, 1, 0, core.Qt__AlignLeft)
	upperGrid.AddWidget(this.invoiceChoose, 1, 1, core.Qt__AlignLeft)
	upperGrid.AddWidget(this.invoiceEdit, 1, 2, core.Qt__AlignLeft)
	upperGrid.AddWidget(this.invoiceCreate, 1, 3, core.Qt__AlignLeft)

	this.invoiceChoose.ConnectPressed(this.chooseInvoice)
	this.invoiceEdit.ConnectPressed(this.editInvoice)

	this.previewBtn = widgets.NewQPushButton2("Preview", nil)
	this.saveTexBtn = widgets.NewQPushButton2("Save .tex", nil)
	this.savePdfBtn = widgets.NewQPushButton2("Save .pdf", nil)
	lowerGrid.AddWidget(this.previewBtn, 0, 0, core.Qt__AlignCenter)
	lowerGrid.AddWidget(this.saveTexBtn, 0, 1, core.Qt__AlignCenter)
	lowerGrid.AddWidget(this.savePdfBtn, 0, 2, core.Qt__AlignCenter)

	//this.updateInvoice()
	this.updateClient()

	return this
}

func (window *MainWindow) chooseClient() {
	wd, err := os.Getwd()
	if err != nil {
		widgets.NewQErrorMessage(window).ShowMessage("Can't get directory!")
	}
	dialog := widgets.NewQFileDialog(window, 0)
	path := dialog.GetOpenFileName(window, "Choose client", wd,
		"*.json", "", 0)
	if len(path) > 0 {
		client, err := DecodeClient(path)
		if err != nil {
			widgets.NewQErrorMessage(window).ShowMessage("Your selected file doesn't seem to be valid!")
		} else {
			selectedClient = client
			clientSelected = true
			window.updateClient()
			window.updateInvoice()
		}
	}
}

func (window *MainWindow) editInvoice() {
	iw := initInvoiceEditDialog(selectedInvoice, window.ParentWidget())
	//iw.SetWindowModality(core.Qt__ApplicationModal)
	iw.Exec()
	selectedInvoice = iw.invoice
	window.updateInvoice()
}

func (window *MainWindow) chooseInvoice() {
	wd, err := os.Getwd()
	if err != nil {
		widgets.NewQErrorMessage(window).ShowMessage("Can't get directory!")
	}
	dialog := widgets.NewQFileDialog(window, 0)
	path := dialog.GetOpenFileName(window, "Choose invoice", wd,
		"*.json", "", 0)
	if len(path) > 0 {
		invoice, err := DecodeInvoice(path)
		if err != nil {
			widgets.NewQErrorMessage(window).ShowMessage("Your selected file doesn't seem to be valid!")
		} else {
			selectedInvoice = invoice
			invoiceSelected = true
			window.updateInvoice()
		}
	}
}

func (window *MainWindow) updateInvoice() {
	if clientSelected && invoiceSelected {
		window.invoiceName.SetText("<" + strconv.Itoa(selectedInvoice.Number) + "> " + selectedInvoice.Project)
	}
	window.invoiceEdit.SetEnabled(clientSelected && invoiceSelected)
	window.invoiceChoose.SetEnabled(clientSelected)
	window.invoiceCreate.SetEnabled(invoiceSelected)
}

func (window *MainWindow) updateClient() {
	if clientSelected {
		window.clientName.SetText(selectedClient.Name)
	}
	window.clientEdit.SetEnabled(clientSelected)

}
