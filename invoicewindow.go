package main

import (
	"math"
	"strconv"

	"github.com/inforbi/inforbi-invoice/data"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

type InvoiceEdit struct {
	widgets.QDialog

	numberSpinner  *widgets.QSpinBox
	projectField   *widgets.QLineEdit
	dateSelector   *widgets.QDateEdit
	dueDays        *widgets.QSpinBox
	itemsButton    *widgets.QPushButton
	totalLabel     *widgets.QLabel
	locationButton *widgets.QPushButton

	invoice data.Invoice
}

func initInvoiceEditDialog(invoice data.Invoice, parent *widgets.QWidget) *InvoiceEdit {
	this := NewInvoiceEdit(parent, core.Qt__Dialog)
	this.invoice = invoice
	this.SetWindowTitle("Edit Invoice")
	formLayout := widgets.NewQFormLayout(this)
	this.numberSpinner = widgets.NewQSpinBox(nil)
	this.numberSpinner.SetMaximum(math.MaxInt32)
	this.numberSpinner.SetMinimum(math.MinInt32)
	this.projectField = widgets.NewQLineEdit(nil)
	this.dueDays = widgets.NewQSpinBox(nil)
	this.itemsButton = widgets.NewQPushButton2("Change Items (0)", nil)
	this.itemsButton.ConnectPressed(this.openItems)
	this.totalLabel = widgets.NewQLabel2("0€", nil, 0)
	this.locationButton = widgets.NewQPushButton2("Change...", nil)
	this.dateSelector = widgets.NewQDateEdit(nil)
	this.dateSelector.SetDisplayFormat("dd-MM-yyyy")
	formLayout.AddRow3("Number", this.numberSpinner)
	formLayout.AddRow3("Project", this.projectField)
	formLayout.AddRow3("Due Days", this.dueDays)
	formLayout.AddRow3("Inv. Date", this.dateSelector)
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
	window.dateSelector.SetDate(core.QDate_FromString2(window.invoice.Date, "yyyy-MM-dd"))
	window.itemsButton.SetText("Change Items (" + strconv.Itoa(len(window.invoice.Items)) + ")")
	window.totalLabel.SetText(strconv.FormatFloat(window.invoice.GetTotal(), 'f', 2, 64) + "€")
}

func (window *InvoiceEdit) toInvoice() {
	window.invoice.Number = window.numberSpinner.Value()
	window.invoice.DueDays = window.dueDays.Value()
	window.invoice.Project = window.projectField.Text()
	window.invoice.Date = window.dateSelector.Date().ToString("yyyy-MM-dd")
}

func (window *InvoiceEdit) openItems() {
	window.toInvoice()
	iw := initItemsWindow(&window.invoice, window.Window())
	iw.SetWindowModality(core.Qt__ApplicationModal)
	iw.Exec()
	window.fromInvoice()
}

func (window *InvoiceEdit) chooseFile() {
	name := chooseFile(window.ParentWidget(), window.invoice.GetFile())
	if len(name) > 0 {
		window.invoice.SetFile(name)
	}
}

func (window *InvoiceEdit) onClose(event *gui.QCloseEvent) {
	window.toInvoice()
	for len(window.invoice.GetFile()) == 0 {
		window.chooseFile()
	}
	window.invoice.EncodeInvoice()
}
