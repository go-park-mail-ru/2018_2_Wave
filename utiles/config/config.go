package config

import (
	"fmt"
	"io/ioutil"
)

// easyjson:json
type Configuration struct {
	SC ServerConfiguration   `json:"server"`
	DC DatabaseConfiguration `json:"database"`
	CC CORSConfiguration     `json:"cors"`
}

// easyjson:json
type ServerConfiguration struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

// easyjson:json
type DatabaseConfiguration struct {
	User   string `json:"user"`
	DBName string `json:"dbname"`
}

// easyjson:json
type CORSConfiguration struct {
	Hosts   []string `json:"hosts"`
	Methods []string `json:"methods"`
}

func Configure(path string) Configuration {
	config := Configuration{}
	data, err := ioutil.ReadFile(path)

	if err != nil {
		data = []byte("{}")
		return config
	}

	if err := config.UnmarshalJSON(data); err != nil {
		fmt.Printf("Incorrect config json: %s\nError: %v", data, err)
	}

	return config
}
