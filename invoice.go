package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"strconv"
	"strings"
)

type Invoice struct {
	Number  int    `json:"number"`
	Project string `json:"project"`
	DueDays int    `json:"dueDays"`
	Items   []Item `json:"items"`
	file    string
}

type Item struct {
	Description string  `json:"description"`
	Quantifier  string  `json:"quantifier"`
	SinglePrice float64 `json:"singlePrice"`
	Quantity    float64 `json:"quantity"`
	Date        string  `json:"date"`
}

func DecodeInvoice(file string) (Invoice, error) {
	cont, err := ioutil.ReadFile(file)
	invoice := Invoice{file: file}
	if err == nil {
		err = json.Unmarshal(cont, &invoice)
	}
	if err == nil && invoice.Number < 0 {
		err = errors.New("invalid file")
	}
	return invoice, err
}

func (invoice Invoice) EncodeInvoice() error {
	cont, err := json.Marshal(invoice)
	if err == nil {
		err = ioutil.WriteFile(invoice.file, cont, 0644)
	}
	return err
}

func (invoice Invoice) GetTotal() (total float64) {
	for _, element := range invoice.Items {
		total += element.Quantity * element.SinglePrice
	}
	return
}

func (item Item) GenerateLatex() string {
	latStr := "\\lineitemu{" + item.Date + "}{" +
		strconv.FormatFloat(item.Quantity, 'f', 2, 64) + "}{" +
		strconv.FormatFloat(item.SinglePrice, 'f', 2, 64) +
		"}{" + item.Quantifier + "}{\\item " + item.Description + "}"
	return latStr
}

func (invoice Invoice) ReplaceTemplate(s string) string {
	s = strings.Replace(s, "$number", strconv.Itoa(invoice.Number), -1)
	s = strings.Replace(s, "$balance", strconv.FormatFloat(invoice.GetTotal(), 'f', 2, 64), -1)
	s = strings.Replace(s, "$project", invoice.Project, -1)
	s = strings.Replace(s, "$duein", strconv.Itoa(invoice.DueDays), -1)
	items := ""
	for _, element := range invoice.Items {
		items += element.GenerateLatex() + "\n"
	}
	s = strings.Replace(s, "$items", items, -1)
	return s
}
