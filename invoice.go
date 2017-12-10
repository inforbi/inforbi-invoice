package main

import (
	"encoding/json"
	"io/ioutil"
	"errors"
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
	SinglePrice float32 `json:"singlePrice"`
	Quantity    int     `json:"quantity"`
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

func (invoice Invoice) EncodeInvoice(file string) error {
	cont, err := json.Marshal(invoice)
	if err == nil {
		err = ioutil.WriteFile(invoice.file, cont, 0644)
	}
	return err
}
