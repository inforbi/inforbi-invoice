package main

import (
	"net"
	"fmt"
	"io/ioutil"
	"strconv"
	"bufio"
	"strings"
	"github.com/nylser/inforbi-invoice/data"
)

func main() {

	conn, err := net.Dial("tcp", "localhost:7714")
	if err != nil {
		println(err)
	}
	cont, err := ioutil.ReadFile("invoice.tex")
	if err != nil {

	}
	fmt.Fprintf(conn, "begin_send"+strconv.Itoa(len(cont))+"\n")
	writer := bufio.NewWriter(conn)
	reader := bufio.NewReader(conn)
	writer.Write(cont)
	writer.Flush()

	response, err := reader.ReadString('\n')
	if strings.Trim(response, "\n") == "success" {
		println("receiving pdf?")
		response, err = reader.ReadString('\n')
		if strings.HasPrefix(response, "begin_send") {
			println("receiving pdf!")
			response = strings.TrimPrefix(response, "begin_send")
			response = strings.Trim(response, "\n")
			ulen, err := strconv.ParseInt(response, 10, 64)
			if err != nil {

			}
			ilen := int(ulen)
			println(ilen)
			file := data.ReceiveBlob(*reader, ilen)
			//n, err := reader.Read(file)
			println(len(file))
			//println(n)
			if len(file) == ilen {
				writer.WriteString("success\n")
				err = ioutil.WriteFile("received.pdf", file, 0644)
			}
		}
	} else {
		return
	}

	conn.Close()
}
