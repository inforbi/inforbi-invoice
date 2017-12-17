package main

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

type InvoiceEdit struct {
	widgets.QMainWindow

	numberSpinner *widgets.QSpinBox
	projectField  *widgets.QLineEdit
	dueDays       *widgets.QSpinBox

	invoice Invoice
}

func initInvoiceEditWindow(invoice Invoice, parent *widgets.QWidget) *InvoiceEdit {
	this := NewInvoiceEdit(parent, core.Qt__Dialog)
	centralWidget := widgets.NewQWidget(this, 0)
	this.SetCentralWidget(centralWidget)
	formLayout := widgets.NewQFormLayout(centralWidget)
	this.numberSpinner = widgets.NewQSpinBox(nil)
	this.projectField = widgets.NewQLineEdit(nil)
	this.dueDays = widgets.NewQSpinBox(nil)
	formLayout.AddRow3("Number", this.numberSpinner)
	formLayout.AddRow3("Project", this.projectField)
	formLayout.AddRow3("Due Days", this.dueDays)
	return this
}

func (window *InvoiceEdit) fromInvoice() {
}
