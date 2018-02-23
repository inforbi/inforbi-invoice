package main

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
	"strconv"
)

type InvoiceEdit struct {
	widgets.QDialog

	numberSpinner  *widgets.QSpinBox
	projectField   *widgets.QLineEdit
	dueDays        *widgets.QSpinBox
	itemsButton    *widgets.QPushButton
	totalLabel     *widgets.QLabel
	locationButton *widgets.QPushButton

	invoice Invoice
}

func initInvoiceEditDialog(invoice Invoice, parent *widgets.QWidget) *InvoiceEdit {
	this := NewInvoiceEdit(parent, core.Qt__Dialog)
	this.invoice = invoice
	this.SetWindowTitle("Edit Invoice")
	formLayout := widgets.NewQFormLayout(this)
	this.numberSpinner = widgets.NewQSpinBox(nil)
	this.projectField = widgets.NewQLineEdit(nil)
	this.dueDays = widgets.NewQSpinBox(nil)
	this.itemsButton = widgets.NewQPushButton2("Change Items (0)", nil)
	this.itemsButton.ConnectPressed(this.openItems)
	this.totalLabel = widgets.NewQLabel2("0€", nil, 0)
	this.locationButton = widgets.NewQPushButton2("Change...", nil)
	formLayout.AddRow3("Number", this.numberSpinner)
	formLayout.AddRow3("Project", this.projectField)
	formLayout.AddRow3("Due Days", this.dueDays)
	formLayout.AddRow3("Items", this.itemsButton)
	formLayout.AddRow3("Total", this.totalLabel)
	formLayout.AddRow3("Save Location", this.locationButton)
	this.locationButton.ConnectPressed(this.chooseFile)
	this.ConnectCloseEvent(this.onClose)
	this.fromInvoice()
	return this
}

func (window *InvoiceEdit) fromInvoice() {
	window.numberSpinner.SetValue(window.invoice.Number)
	window.projectField.SetText(window.invoice.Project)
	window.dueDays.SetValue(window.invoice.DueDays)
	window.itemsButton.SetText("Change Items (" + strconv.Itoa(len(window.invoice.Items)) + ")")
	window.totalLabel.SetText(strconv.FormatFloat(window.invoice.GetTotal(), 'f', 2, 64) + "€")
}

func (window *InvoiceEdit) toInvoice() {
	window.invoice.Number = window.numberSpinner.Value()
	window.invoice.DueDays = window.dueDays.Value()
	window.invoice.Project = window.projectField.Text()
}

func (window *InvoiceEdit) openItems() {
	window.toInvoice()
	iw := initItemsWindow(&window.invoice, window.Window())
	iw.SetWindowModality(core.Qt__ApplicationModal)
	iw.Exec()
	window.fromInvoice()
}

func (window *InvoiceEdit) chooseFile() {
	name := chooseFile(window.ParentWidget(), window.invoice.file)
	if len(name) > 0 {
		window.invoice.file = name
	}
}

func (window *InvoiceEdit) onClose(event *gui.QCloseEvent) {
	window.toInvoice()
	for len(window.invoice.file) == 0 {
		window.chooseFile()
	}
	window.invoice.EncodeInvoice()
}
