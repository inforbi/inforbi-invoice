package main

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

type ClientEdit struct {
	widgets.QDialog

	nameField    *widgets.QLineEdit
	contactField *widgets.QLineEdit
	streetField  *widgets.QLineEdit
	cityField    *widgets.QLineEdit

	client Client
}

func initClientEditDialog(client Client, parent *widgets.QWidget) *ClientEdit {
	this := NewClientEdit(parent, core.Qt__Dialog)
	return this
}
