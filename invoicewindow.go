package main

import (
	"github.com/therecipe/qt/widgets"
	"github.com/therecipe/qt/core"
)

type InvoiceEdit struct {
	widgets.QMainWindow

	numberSpinner *widgets.QSpinBox
	projectField *widgets.QLineEdit
	dueDays *widgets.QSpinBox

	invoice Invoice
}

func initInvoiceEditWindow(invoice Invoice, parent *widgets.QWidget) *InvoiceEdit {
	this := NewInvoiceEdit(parent, core.Qt__Dialog)
	centralWidget := widgets.NewQWidget(this, 0)
	this.SetCentralWidget(centralWidget)
	return this
}