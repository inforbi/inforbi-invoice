package data

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"strings"
)

type Client struct {
	Name    string `json:"name"`
	Contact string `json:"contact"`
	Street  string `json:"street"`
	City    string `json:"city"`
	file    string
}

func (client *Client) GetFile() (string) {
	return client.file
}

func (client *Client) SetFile(file string) {
	client.file = file
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

func (client Client) ReplaceTemplate(s string) string {
	s = strings.Replace(s, "$clientContact", client.Contact, 1)
	s = strings.Replace(s, "$clientStreet", client.Street, 1)
	s = strings.Replace(s, "$clientCity", client.City, 1)
	s = strings.Replace(s, "$clientName", client.Name, 1)
	return s
}
