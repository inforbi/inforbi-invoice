package main

import (
	"net"
	"bufio"
	"strings"
	"strconv"
	"io/ioutil"
	"github.com/nylser/inforbi-invoice/data"
	"log"
	"path/filepath"
	"os/exec"
	"time"
	//"os"
	"os"
)

func main() {

	ln, err := net.Listen("tcp", ":7714")
	if err != nil {
		println(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			println(err)
		}
		go handleConnection(conn)
	}

}

func handleConnection(conn net.Conn) {
	println(conn.RemoteAddr().String() + " has connected!")
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			println(err)
			break
		}
		if strings.HasPrefix(line, "begin_send") {
			print(line)
			size := strings.TrimPrefix(line, "begin_send")
			size = strings.Trim(size, "\n")
			result, err := strconv.ParseInt(size, 10, 64)
			if err != nil {
				return
			}
			iresult := int(result)
			file := data.ReceiveBlob(*reader, iresult)

			println("Result")
			println("Given: " + strconv.Itoa(iresult))
			println("Received: " + strconv.Itoa(len(file)))
			if int(result) == len(file) {
				writer.WriteString("success\n")
			} else {
				writer.WriteString("fail\n")
				continue
			}

			// Do thing and send PDF
			f, err := createFromTemplate(file)
			ioutil.WriteFile("pre_send.pdf", f, 0644)
			if err != nil {
			}
			writer.WriteString("begin_send" + strconv.Itoa(len(f)) + "\n")
			writer.Write(f)
			writer.Flush()
			line, err := reader.ReadString('\n')
			if err != nil {
			}
			if line == "success\n" {
				println("Yay!")
				return
			}
		}
	}
}

func createFromTemplate(input []byte) ([]byte, error) {
	dir, err := ioutil.TempDir("", "preview")
	if err != nil {
		log.Fatal(err)
	}
	tmplat := filepath.Join(dir, "render.tex")
	tmpcls := filepath.Join(dir, "dapper-invoice.cls")

	data.CopyDir("Fonts", filepath.Join(dir, "Fonts"))
	data.CopyFile("dapper-invoice.cls", tmpcls)

	err = ioutil.WriteFile(tmplat, input, 0644)

	command := exec.Command("xelatex", "-synctex=1", "-interaction=nonstopmode", "render.tex")
	command.Dir = dir
	out, err := command.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	outstr := string(out)
	if strings.Contains(strings.ToLower(outstr), "rerun") {
		command := exec.Command("xelatex", "-synctex=1", "-interaction=nonstopmode", "render.tex")
		command.Dir = dir
		command.Run()
	}
	cont, err := ioutil.ReadFile(filepath.Join(dir, "render.pdf"))
	go func() {
		time.Sleep(1 * time.Second)
		os.RemoveAll(dir)
	}()

	return cont, err
}
