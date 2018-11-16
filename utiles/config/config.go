package config

import (
	"io/ioutil"
)

// easyjson:json
type Configuration struct {
	SC ServerConfiguration `json:"server"`
	CC CORSConfiguration   `json:"cors"`
}

// easyjson:json
type ServerConfiguration struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

// easyjson:json
type CORSConfiguration struct {
	Origins     []string `json:"origins"`
	Headers     []string `json:"headers"`
	Credentials string   `json:"credentials"`
	Methods     []string `json:"methods"`
}

func New(path string) Configuration {
	config := Configuration{}
	data, err := ioutil.ReadFile(path)

	if err != nil {
		//data = []byte("{}")
		return config
	}

	if err := config.UnmarshalJSON(data); err != nil {
		//fmt.Printf("Incorrect config json: %s\nError: %v", data, err)
		return config
	}

	return config
}
