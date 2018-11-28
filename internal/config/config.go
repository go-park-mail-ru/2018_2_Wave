package config

import (
	"io/ioutil"
)

// easyjson:json
type Configuration struct {
	SC   ServerConfiguration `json:"server"`
	CC   CORSConfiguration   `json:"cors"`
	Auth GRPCConfiguration   `json:"auth"`
	Game GameConfiguration   `json:"game"`
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

// easyjson:json
type GRPCConfiguration struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

// easyjson:json
type GameConfiguration struct {
	WsPort string `json:"ws_port"`
	GRPCConfiguration
}

func Configure(path string) Configuration {
	config := Configuration{}
	data, err := ioutil.ReadFile(path)

	if err != nil {
		return Configuration{}
	}

	if err := config.UnmarshalJSON(data); err != nil {
		return config
	}

	return config
}
