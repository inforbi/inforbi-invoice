package main

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
	"github.com/nylser/inforbi-invoice/data"
)

type ItemsDialog struct {
	widgets.QDialog

	table   *widgets.QTableWidget
	invoice *data.Invoice
}

type CustomItem struct {
	widgets.QItemDelegate
}

func (window *ItemsDialog) initCustomItem() *CustomItem {
	item := NewCustomItem(nil)
	item.ConnectCreateEditor(window.createEditor)
	return item
}

func (window *ItemsDialog) createEditor(parent *widgets.QWidget, option *widgets.QStyleOptionViewItem, index *core.QModelIndex) *widgets.QWidget {
	editor := widgets.NewQDoubleSpinBox(parent)
	return editor.QWidget_PTR()
}

var firstRow = []string{"Description", "Date", "Quantifier", "Quantity", "Cost/Q", "Total"}

func initItemsWindow(invoice *data.Invoice, parent *widgets.QWidget) *ItemsDialog {
	this := NewItemsDialog(parent, core.Qt__Dialog)
	this.SetWindowTitle("Edit Invoice Items")
	this.SetMinimumHeight(500)
	grid := widgets.NewQGridLayout(this)
	this.invoice = invoice
	this.table = widgets.NewQTableWidget2(len(invoice.Items)+1, 6, nil)
	grid.AddWidget3(this.table, 1, 0, 5, 6, core.Qt__AlignLeft)
	this.table.SetHorizontalHeaderLabels(firstRow)
	this.table.SetSizeAdjustPolicy(widgets.QAbstractScrollArea__AdjustToContents)
	this.table.SetItemDelegateForColumn(3, this.initCustomItem())
	this.table.SetItemDelegateForColumn(4, this.initCustomItem())
	this.table.ConnectItemChanged(this.itemChanged)
	this.table.SizePolicy().SetHorizontalPolicy(widgets.QSizePolicy__Expanding)
	this.table.HorizontalHeader().SetStretchLastSection(true)
	this.AdjustSize()
	this.ConnectCloseEvent(this.onClose)
	this.fromItems()
	return this
}

func (window *ItemsDialog) onClose(event *gui.QCloseEvent) {
	window.toItems()
}

func (window *ItemsDialog) itemChanged(new *widgets.QTableWidgetItem) {
	if window.countFreeRows() == 0 {
		window.table.InsertRow(window.table.RowCount())
	} else if new.Row() != window.table.RowCount()-1 {
		if window.isRowEmptyExcept(new.Row(), new) {
			window.table.RemoveRow(new.Row())
		}
	}
	if new.Column() == 5 {
	}
	if new.Column() == 3 || new.Column() == 4 {
		var ok bool
		var col3, col4 float64
		if window.table.Item(new.Row(), 3) != nil && window.table.Item(new.Row(), 4) != nil {
			col3 = window.table.Item(new.Row(), 3).Data(0).ToDouble(ok)
			col4 = window.table.Item(new.Row(), 4).Data(0).ToDouble(ok)
		}
		new5 := widgets.NewQTableWidgetItem(0)
		new5.SetData(0, core.NewQVariant12(col3*col4))
		window.table.SetItem(new.Row(), 5, new5)
	}
}

func (window *ItemsDialog) countFreeRows() int {
	rows := 0
	for i := 0; i < window.table.RowCount(); i++ {
		if window.isRowEmpty(i) {
			rows++
		}
	}
	return rows

}

func (window *ItemsDialog) removeFreeBetween() {
	for i := 0; i < window.table.RowCount()-1; i++ {
		if window.isRowEmpty(i) {
			window.table.RemoveRow(i)
		}
	}
}

func (window *ItemsDialog) isRowEmpty(row int) bool {
	for i := 0; i < 6; i++ {
		if len(window.table.Item(row, i).Text()) != 0 {
			return false
		}
	}
	return true
}

func (window *ItemsDialog) isRowEmptyExcept(row int, exception *widgets.QTableWidgetItem) bool {
	for i := 0; i < 6; i++ {
		if window.table.Item(row, i) != exception && len(window.table.Item(row, i).Text()) != 0 {
			return false
		}
	}
	return true
}

func (window *ItemsDialog) fromItems() {
	items := window.invoice.Items
	for i := 0; i < len(items); i++ {
		item := items[i]
		for j := 0; j < 6; j++ {
			var wItem *widgets.QTableWidgetItem
			switch j {
			case 0:
				wItem = widgets.NewQTableWidgetItem2(item.Description, 0)
			case 1:
				wItem = widgets.NewQTableWidgetItem2(item.Date, 0)
			case 2:
				wItem = widgets.NewQTableWidgetItem2(item.Quantifier, 0)
			case 3:
				wItem = widgets.NewQTableWidgetItem(0)
				wItem.SetData(0, core.NewQVariant12(item.Quantity))
			case 4:
				wItem = widgets.NewQTableWidgetItem(0)
				wItem.SetData(0, core.NewQVariant12(item.SinglePrice))
			case 5:
				wItem = widgets.NewQTableWidgetItem(0)
				wItem.SetData(0, core.NewQVariant12(item.SinglePrice*item.Quantity))
			}
			window.table.SetItem(i, j, wItem)
		}
	}
}

func (window *ItemsDialog) toItems() {
	var items []data.Item
	for i := 0; i < window.table.RowCount(); i++ {
		if !window.isRowEmpty(i) {
			item := data.Item{}
			for j := 0; j < 6; j++ {
				wItem := window.table.Item(i, j)
				var ok bool
				switch j {
				case 0:
					item.Description = wItem.Text()
				case 1:
					item.Date = wItem.Text()
				case 2:
					item.Quantifier = wItem.Text()
				case 3:
					item.Quantity = wItem.Data(0).ToDouble(ok)
				case 4:
					item.SinglePrice = wItem.Data(0).ToDouble(ok)
				case 5:
					// Do nothing
				}
			}
			items = append(items, item)
		}
	}
	window.invoice.Items = items
}
