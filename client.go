package main

import (
	"encoding/json"
	"io/ioutil"
	"errors"
)

type Client struct {
	Name    string `json:"name"`
	Contact string `json:"contact"`
	Street  string `json:"street"`
	City    string `json:"city"`
	file    string
}

func DecodeClient(file string) (Client, error) {
	cont, err := ioutil.ReadFile(file)
	client := Client{file: file}
	if err == nil {
		err = json.Unmarshal(cont, &client)
	}
	if err == nil && len(client.Name) == 0 {
		err = errors.New("invalid file")
	}
	return client, err
}

func (client Client) EncodeClient() error {
	cont, err := json.Marshal(client)
	if err == nil {
		err = ioutil.WriteFile(client.file, cont, 0644)
	}
	return err
}
