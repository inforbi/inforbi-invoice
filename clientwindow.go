package main

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

type ClientEdit struct {
	widgets.QMainWindow

	nameField    *widgets.QLineEdit
	contactField *widgets.QLineEdit
	streetField  *widgets.QLineEdit
	cityField    *widgets.QLineEdit

	client Client
}

func initClientEditWindow(client Client, parent *widgets.QWidget) *ClientEdit {
	this := NewClientEdit(parent, core.Qt__Dialog)
	centralWidget := widgets.NewQWidget(this, 0)
	this.SetCentralWidget(centralWidget)
	return this
}
