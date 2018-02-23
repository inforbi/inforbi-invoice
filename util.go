package main

import (
	"github.com/therecipe/qt/widgets"
	"log"
	"os/user"
	"path"
)

func chooseFile(parent widgets.QWidget_ITF, currentFile string) string {
	dia := widgets.NewQFileDialog(parent, 0)
	if len(currentFile) > 0 {
		dia.SelectFile(currentFile)
	} else {
		usr, err := user.Current()
		if err != nil {
			log.Fatal(err)
		}
		currentFile = usr.HomeDir
	}
	name := dia.GetSaveFileName(parent, "Choose location", path.Dir(currentFile), "*.json (JSON Files)", "*.json (JSON Files)", 0)
	return name
}
