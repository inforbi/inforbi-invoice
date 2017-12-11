package main

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
	"strconv"
)

type ItemsEdit struct {
	widgets.QMainWindow

	table *widgets.QTableWidget
	items []Item
}

type CustomItem struct {
	widgets.QItemDelegate
}

func (window *ItemsEdit) initCustomItem() *CustomItem {
	item := NewCustomItem(nil)
	item.ConnectCreateEditor(window.createEditor)
	return item
}

func (window *ItemsEdit) createEditor(parent *widgets.QWidget, option *widgets.QStyleOptionViewItem, index *core.QModelIndex) *widgets.QWidget {
	editor := widgets.NewQDoubleSpinBox(parent)
	return editor.QWidget_PTR()
}

var firstRow = []string{"Description", "Date", "Quantifier", "Quantity", "Cost/Q", "Total"}

func initItemsWindow(items []Item, parent *widgets.QWidget) *ItemsEdit {
	this := NewItemsEdit(parent, core.Qt__Dialog)
	this.SetWindowTitle("Edit Invoice Items")
	this.SetMinimumHeight(500)
	widget := widgets.NewQWidget(this, 0)
	this.SetCentralWidget(widget)
	grid := widgets.NewQGridLayout(widget)
	this.items = items
	this.table = widgets.NewQTableWidget2(len(items)+1, 6, nil)
	grid.AddWidget3(this.table, 1, 0, 5, 6, core.Qt__AlignLeft)
	this.table.SetHorizontalHeaderLabels(firstRow)
	this.table.SetSizeAdjustPolicy(widgets.QAbstractScrollArea__AdjustToContentsOnFirstShow)
	this.table.SetItemDelegateForColumn(2, this.initCustomItem())
	this.table.SetItemDelegateForColumn(3, this.initCustomItem())

	this.fromItems()
	this.table.ConnectItemChanged(this.itemChanged)
	this.AdjustSize()
	return this
}

func (window *ItemsEdit) itemChanged(new *widgets.QTableWidgetItem) {
	if window.countFreeRows() == 0 {
		window.table.InsertRow(window.table.RowCount())
	} else if new.Row() != window.table.RowCount()-1 {
		if window.isRowEmptyExcept(new.Row(), new) {
			window.table.RemoveRow(new.Row())
		}
	}

	if new.Column() == 2 || new.Column() == 3 {
		var ok bool
		col2 := window.table.Item(new.Row(), 2).Data(0).ToDouble(ok)
		col3 := window.table.Item(new.Row(), 3).Data(0).ToDouble(ok)
		new5 := widgets.NewQTableWidgetItem(0)
		new5.SetData(0, core.NewQVariant12(col2*col3))
		window.table.SetItem(new.Row(), 5, new5)
	}
}

func (window *ItemsEdit) countFreeRows() int {
	rows := 0
	for i := 0; i < window.table.RowCount(); i++ {
		if window.isRowEmpty(i) {
			rows++
		}
	}
	return rows

}

func (window *ItemsEdit) removeFreeBetween() {
	for i := 0; i < window.table.RowCount()-1; i++ {
		if window.isRowEmpty(i) {
			window.table.RemoveRow(i)
		}
	}
}

func (window *ItemsEdit) isRowEmpty(row int) bool {
	for i := 0; i < 6; i++ {
		if len(window.table.Item(row, i).Text()) != 0 {
			return false
		}
	}
	return true
}

func (window *ItemsEdit) isRowEmptyExcept(row int, exception *widgets.QTableWidgetItem) bool {
	for i := 0; i < 6; i++ {
		if window.table.Item(row, i) != exception && len(window.table.Item(row, i).Text()) != 0 {
			return false
		}
	}
	return true
}

func (window *ItemsEdit) fromItems() {
	for i := 0; i < len(window.items); i++ {
		item := window.items[i]
		for j := 0; j < 6; j++ {
			var wItem *widgets.QTableWidgetItem
			switch j {
			case 0:
				wItem = widgets.NewQTableWidgetItem2(item.Description, 0)
			case 1:
				wItem = widgets.NewQTableWidgetItem2(item.Date, 0)
			case 2:
				wItem = widgets.NewQTableWidgetItem(0)
				wItem.SetData(0, core.NewQVariant12(item.SinglePrice))
			case 3:
				wItem = widgets.NewQTableWidgetItem(0)
				wItem.SetData(0, core.NewQVariant12(item.Quantity))
			case 4:
				wItem = widgets.NewQTableWidgetItem2(item.Quantifier, 0)
			case 5:
				wItem = widgets.NewQTableWidgetItem(0)
				wItem.SetData(0, core.NewQVariant12(item.SinglePrice*item.Quantity))
			}
			window.table.SetItem(i, j, wItem)
		}
	}
}

func (window *ItemsEdit) toItems() {
	window.items = []Item{}
	for i := 0; i < window.table.RowCount(); i++ {
		if !window.isRowEmpty(i) {
			item := Item{}
			for j := 0; j < 6; j++ {
				wItem := window.table.Item(i, j)
				switch j {
				case 0:
					item.Description = wItem.Text()
				case 1:
					item.Date = wItem.Text()
				case 2:
					item.SinglePrice, _ = strconv.ParseFloat(wItem.Text(), 64)
				case 3:
					item.Quantity, _ = strconv.ParseFloat(wItem.Text(), 64)
				case 4:
					item.Quantifier = wItem.Text()
				case 5:
					// Do nothing
				}
			}
			window.items = append(window.items, item)
		}
	}

}
