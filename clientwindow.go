package main

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
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
	this.client = client
	this.SetWindowTitle("Edit Client")
	formLayout := widgets.NewQFormLayout(this)
	this.nameField = widgets.NewQLineEdit(nil)
	this.contactField = widgets.NewQLineEdit(nil)
	this.streetField = widgets.NewQLineEdit(nil)
	this.cityField = widgets.NewQLineEdit(nil)
	formLayout.AddRow3("Client Name", this.nameField)
	formLayout.AddRow3("Client Contact", this.contactField)
	formLayout.AddRow3("Client Street", this.streetField)
	formLayout.AddRow3("Client City", this.cityField)
	this.ConnectCloseEvent(this.onClose)
	return this
}

func (window *ClientEdit) fromClient() {
	window.nameField.SetText(window.client.Name)
	window.contactField.SetText(window.client.Contact)
	window.cityField.SetText(window.client.City)
	window.streetField.SetText(window.client.Street)
}

func (window *ClientEdit) toClient() {
	window.client.Name = window.nameField.Text()
	window.client.Contact = window.contactField.Text()
	window.client.City = window.cityField.Text()
	window.client.Street = window.streetField.Text()
}

func (window *ClientEdit) onClose(event *gui.QCloseEvent) {
	window.toClient()
	window.client.EncodeClient()
}
