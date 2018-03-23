package main

import (
	"github.com/skratchdot/open-golang/open"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
	"github.com/nylser/inforbi-invoice/data"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"net"
	"bufio"
)

type MainWindow struct {
	widgets.QMainWindow

	clientName   *widgets.QLabel
	clientChoose *widgets.QPushButton
	clientEdit   *widgets.QPushButton
	clientCreate *widgets.QPushButton

	invoiceName   *widgets.QLabel
	invoiceChoose *widgets.QPushButton
	invoiceEdit   *widgets.QPushButton
	invoiceCreate *widgets.QPushButton

	previewBtn   *widgets.QPushButton
	saveTexBtn   *widgets.QPushButton
	savePdfBtn   *widgets.QPushButton
	useRemoteBox *widgets.QCheckBox
}

var (
	clientSelected  = false
	invoiceSelected = false
	selectedClient  data.Client
	selectedInvoice data.Invoice
)

func initMainWindow() *MainWindow {
	this := NewMainWindow(nil, 0)
	this.SetWindowTitle(core.QCoreApplication_ApplicationName())
	upperWidget := widgets.NewQWidget(this, 0)
	lowerWidget := widgets.NewQWidget(nil, 0)
	upperGrid := widgets.NewQGridLayout(upperWidget)
	lowerGrid := widgets.NewQGridLayout(lowerWidget)
	upperGrid.SetSpacing(10)
	lowerGrid.SetSpacing(50)

	this.SetCentralWidget(upperWidget)

	upperGrid.AddWidget3(lowerWidget, 2, 0, 1, 4, core.Qt__AlignCenter)

	this.clientName = widgets.NewQLabel2("<No client selected>", nil, 0)
	this.clientChoose = widgets.NewQPushButton2("Choose...", nil)
	this.clientEdit = widgets.NewQPushButton2("Edit", nil)
	this.clientCreate = widgets.NewQPushButton2("Create", nil)

	this.clientChoose.ConnectPressed(this.chooseClient)
	this.clientEdit.ConnectPressed(this.editClient)
	this.clientCreate.ConnectPressed(this.createClient)

	upperGrid.AddWidget(this.clientName, 0, 0, core.Qt__AlignLeft)
	upperGrid.AddWidget(this.clientChoose, 0, 1, core.Qt__AlignLeft)
	upperGrid.AddWidget(this.clientEdit, 0, 2, core.Qt__AlignLeft)
	upperGrid.AddWidget(this.clientCreate, 0, 3, core.Qt__AlignLeft)

	this.invoiceName = widgets.NewQLabel2("<No invoice selected>", nil, 0)
	this.invoiceChoose = widgets.NewQPushButton2("Choose...", nil)
	this.invoiceEdit = widgets.NewQPushButton2("Edit", nil)
	this.invoiceCreate = widgets.NewQPushButton2("Create", nil)

	upperGrid.AddWidget(this.invoiceName, 1, 0, core.Qt__AlignLeft)
	upperGrid.AddWidget(this.invoiceChoose, 1, 1, core.Qt__AlignLeft)
	upperGrid.AddWidget(this.invoiceEdit, 1, 2, core.Qt__AlignLeft)
	upperGrid.AddWidget(this.invoiceCreate, 1, 3, core.Qt__AlignLeft)

	this.invoiceChoose.ConnectPressed(this.chooseInvoice)
	this.invoiceEdit.ConnectPressed(this.editInvoice)
	this.invoiceCreate.ConnectPressed(this.createInvoice)

	this.previewBtn = widgets.NewQPushButton2("Preview", nil)
	this.saveTexBtn = widgets.NewQPushButton2("Save .tex", nil)
	this.savePdfBtn = widgets.NewQPushButton2("Save .pdf", nil)
	this.useRemoteBox = widgets.NewQCheckBox2("Use remote-server for rendering", nil)
	this.useRemoteBox.SetChecked(true)

	this.previewBtn.ConnectPressed(this.preview)
	this.savePdfBtn.ConnectPressed(this.savePDF)
	this.saveTexBtn.ConnectPressed(this.saveTex)

	lowerGrid.SetSpacing(2)
	lowerGrid.AddWidget3(this.useRemoteBox, 0, 0, 1, 3, core.Qt__AlignCenter)
	lowerGrid.AddWidget(this.previewBtn, 1, 0, core.Qt__AlignCenter)
	lowerGrid.AddWidget(this.saveTexBtn, 1, 1, core.Qt__AlignCenter)
	lowerGrid.AddWidget(this.savePdfBtn, 1, 2, core.Qt__AlignCenter)

	this.updateInvoice()
	this.updateClient()
	this.updateBottomBtns()

	return this
}

func (window *MainWindow) createClient() {
	selectedClient = data.Client{}
	window.editClient()
}

func (window *MainWindow) createInvoice() {
	if clientSelected {
		selectedInvoice = data.Invoice{}
		window.editInvoice()
	}
}

func (window *MainWindow) chooseClient() {
	wd, err := os.Getwd()
	if err != nil {
		widgets.NewQErrorMessage(window).ShowMessage("Can't get directory!")
	}
	dialog := widgets.NewQFileDialog(window, 0)
	path := dialog.GetOpenFileName(window, "Choose client", wd,
		"*.json", "", 0)
	if len(path) > 0 {
		client, err := data.DecodeClient(path)
		if err != nil {
			widgets.NewQErrorMessage(window).ShowMessage("Your selected file doesn't seem to be valid!")
		} else {
			selectedClient = client
			clientSelected = true
			window.updateClient()
			window.updateInvoice()
		}
	}
}

func (window *MainWindow) editClient() {
	cw := initClientEditDialog(selectedClient, window.ParentWidget())
	cw.Exec()
	if cw.Result() == 0 {
		selectedClient = cw.client
		window.updateClient()
		if !clientSelected {
			clientSelected = true
		}
	} else {
		selectedClient.EncodeClient()
	}

}

func (window *MainWindow) editInvoice() {
	iw := initInvoiceEditDialog(selectedInvoice, window.ParentWidget())
	iw.Exec()
	if iw.Result() == 0 {
		selectedInvoice = iw.invoice
		window.updateInvoice()
		if !invoiceSelected {
			invoiceSelected = true
		}
	} else {
		selectedInvoice.EncodeInvoice()
	}
}

func (window *MainWindow) chooseInvoice() {
	wd, err := os.Getwd()
	if err != nil {
		widgets.NewQErrorMessage(window).ShowMessage("Can't get directory!")
	}
	dialog := widgets.NewQFileDialog(window, 0)
	path := dialog.GetOpenFileName(window, "Choose invoice", wd,
		"*.json", "", 0)
	if len(path) > 0 {
		invoice, err := data.DecodeInvoice(path)
		if err != nil {
			widgets.NewQErrorMessage(window).ShowMessage("Your selected file doesn't seem to be valid!")
		} else {
			selectedInvoice = invoice
			invoiceSelected = true
			window.updateInvoice()
		}
	}
}

func (window *MainWindow) updateInvoice() {
	if clientSelected && invoiceSelected {
		window.invoiceName.SetText("<" + strconv.Itoa(selectedInvoice.Number) + "> " + selectedInvoice.Project)
	}
	window.invoiceEdit.SetEnabled(clientSelected && invoiceSelected)
	window.invoiceChoose.SetEnabled(clientSelected)
	window.invoiceCreate.SetEnabled(invoiceSelected)
	window.updateBottomBtns()
}

func (window *MainWindow) updateClient() {
	if clientSelected {
		window.clientName.SetText(selectedClient.Name)
	}
	window.clientEdit.SetEnabled(clientSelected)
	window.updateBottomBtns()
}

func (window *MainWindow) updateBottomBtns() {
	condition := invoiceSelected && clientSelected
	window.previewBtn.SetEnabled(condition)
	window.savePdfBtn.SetEnabled(condition)
	window.saveTexBtn.SetEnabled(condition)
}

func (window *MainWindow) generateLatex() string {
	if clientSelected && invoiceSelected {
		bytes, err := ioutil.ReadFile("invoice.pylat")
		if err != nil {
			panic("Test")
		}
		template := string(bytes)
		template = selectedClient.ReplaceTemplate(template)
		template = selectedInvoice.ReplaceTemplate(template)
		return template
	} else {
		return ""
	}
}

func (window *MainWindow) preview() {
	dir, err := ioutil.TempDir("", "preview")
	tmpPDF := filepath.Join(dir, "preview.pdf")
	if err != nil {
		log.Fatal(err)
	}
	if window.useRemoteBox.IsChecked() {
		err := window.remoteRender(tmpPDF)
		if err != nil {
			println(err)
		}
		go func() {
			time.Sleep(1 * time.Second)
			os.RemoveAll(dir)
		}()

	} else {
		window.localRender(dir, tmpPDF)
		open.Run(filepath.Join(dir, "preview.pdf"))
		go func() {
			time.Sleep(1 * time.Second)
			os.RemoveAll(dir)
		}()
	}

}

func (window *MainWindow) savePDF() {
	wd, err := os.Getwd()
	if err != nil {
		widgets.NewQErrorMessage(window).ShowMessage("Can't get directory!")
	}
	dialog := widgets.NewQFileDialog(window, 0)
	path := dialog.GetSaveFileName(window, "Save invoice", wd,
		"*.pdf", "*.pdf", 0)
	if len(path) == 0 {
		widgets.NewQErrorMessage(window).ShowMessage("Can't save file without selected destination!")
		return
	}

	if window.useRemoteBox.IsChecked() {
		window.remoteRender(path)
	} else {
		dir, err := ioutil.TempDir("", "save")
		if err != nil {
			log.Fatal(err)
		}
		window.localRender(dir, path)
		os.RemoveAll(dir)
	}
}

func (window *MainWindow) saveTex() {
	wd, err := os.Getwd()
	if err != nil {
		widgets.NewQErrorMessage(window).ShowMessage("Can't get directory!")
	}
	dialog := widgets.NewQFileDialog(window, 0)
	path := dialog.GetSaveFileName(window, "Save latex", wd,
		"*.tex", "*.tex", 0)
	if len(path) == 0 {
		widgets.NewQErrorMessage(window).ShowMessage("Can't save file without selected destination!")
		return
	}
	latex := window.generateLatex()
	err = ioutil.WriteFile(path, []byte(latex), 0644)
	if err != nil {
		log.Fatal(err)
		widgets.NewQErrorMessage(window).ShowMessage("Couldn't save file! " + err.Error())
	}

}

func (window *MainWindow) localRender(target_dir string, target_file string) {
	latex := window.generateLatex()
	tmplat := filepath.Join(target_dir, "preview.tex")
	tmpcls := filepath.Join(target_dir, "dapper-invoice.cls")

	data.CopyDir("Fonts", filepath.Join(target_dir, "Fonts"))
	data.CopyFile("dapper-invoice.cls", tmpcls)

	err := ioutil.WriteFile(tmplat, []byte(latex), 0644)
	if err != nil {
		log.Fatal(err)
		widgets.NewQErrorMessage(window).ShowMessage("Couldn't save file! " + err.Error())
	} else {
		go func() {
			command := exec.Command("xelatex", "-synctex=1", "-interaction=nonstopmode", "render.tex")
			command.Dir = target_dir
			out, err := command.CombinedOutput()
			if err != nil {
				log.Fatal(err)
			}
			outstr := string(out)
			if strings.Contains(strings.ToLower(outstr), "rerun") {
				command := exec.Command("xelatex", "-synctex=1", "-interaction=nonstopmode", "render.tex")
				command.Dir = target_dir
				command.Run()
			}
			data.CopyFile(filepath.Join(target_dir, "render.pdf"), target_file)
		}()
	}
}

func (window *MainWindow) remoteRender(target_file string) (error) {
	latex := window.generateLatex()
	conn, err := net.Dial("tcp", "mineguild.net:7714")
	if err != nil {
		println(err)
		return err
	}
	writer := bufio.NewWriter(conn)
	reader := bufio.NewReader(conn)
	go func() {
		writer.WriteString("begin_send" + strconv.Itoa(len(latex)) + "\n")
		writer.Write([]byte(latex))
		writer.Flush()
		response, err := reader.ReadString('\n')
		if err != nil {
			widgets.NewQErrorMessage(window).ShowMessage("Error communicating!")
			return
		}
		if strings.Trim(response, "\n") == "success" {
			println("receiving pdf?")
			response, err = reader.ReadString('\n')
			if strings.HasPrefix(response, "begin_send") {
				println("receiving pdf!")
				response = strings.TrimPrefix(response, "begin_send")
				response = strings.Trim(response, "\n")
				ulen, err := strconv.ParseInt(response, 10, 64)
				if err != nil {
					return
				}
				ilen := int(ulen)
				println(ilen)
				file := data.ReceiveBlob(*reader, ilen)
				if len(file) == ilen {
					writer.WriteString("success\n")
					err = ioutil.WriteFile(target_file, file, 0644)
				}
			}
		}
	}()
	return nil
}
