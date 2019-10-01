package main

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
	"github.com/inforbi/inforbi-invoice/data"
)

type ClientEdit struct {
	widgets.QDialog

	nameField      *widgets.QLineEdit
	contactField   *widgets.QLineEdit
	streetField    *widgets.QLineEdit
	cityField      *widgets.QLineEdit
	locationButton *widgets.QPushButton

	client data.Client
}

func initClientEditDialog(client data.Client, parent *widgets.QWidget) *ClientEdit {
	this := NewClientEdit(parent, core.Qt__Dialog)
	this.client = client
	this.SetWindowTitle("Edit Client")
	formLayout := widgets.NewQFormLayout(this)
	this.nameField = widgets.NewQLineEdit(nil)
	this.contactField = widgets.NewQLineEdit(nil)
	this.streetField = widgets.NewQLineEdit(nil)
	this.cityField = widgets.NewQLineEdit(nil)
	this.locationButton = widgets.NewQPushButton2("Change...", nil)
	formLayout.AddRow3("Client Name", this.nameField)
	formLayout.AddRow3("Client Contact", this.contactField)
	formLayout.AddRow3("Client Street", this.streetField)
	formLayout.AddRow3("Client City", this.cityField)
	formLayout.AddRow3("Save Location", this.locationButton)
	this.locationButton.ConnectPressed(this.chooseFile)
	this.ConnectCloseEvent(this.onClose)
	this.fromClient()
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

func (window *ClientEdit) chooseFile() {
	name := chooseFile(window.ParentWidget(), window.client.GetFile())
	if len(name) > 0 {
		window.client.SetFile(name)
	}
}

func (window *ClientEdit) onClose(event *gui.QCloseEvent) {
	window.toClient()
	if len(window.client.GetFile()) == 0 {
		window.chooseFile()
	}
	window.client.EncodeClient()
}
