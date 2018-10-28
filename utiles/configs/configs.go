package configs

import "io/ioutil"

//go:generate easyjson -output_filename configs.gen.go configs.go

// CORSConfig configuration
// easyjson:json
type CORSConfig struct {
	Credentials bool     `json:"credentials"`
	Origins     []string `json:"origins"`
	Methods     []string `json:"methods"`
	Headers     []string `json:"headers"`
}

// ServerConfig configuration
// easyjson:json
type ServerConfig struct {
	Port string `json:"port"`
	Log  string `json:"log"`
}

// DatabaseConfig configuration
// easyjson:json
type DatabaseConfig struct {
}

// WalhallaConfig configuration
// easyjson:json
type WalhallaConfig struct {
	MainTmpl   string   `json:"mainTmpl"`
	OtherTmpls []string `json:"otherTmpls"`
}

// MainConfig ...
// easyjson:json
type MainConfig struct {
	CORS     CORSConfig     `json:"cors"`
	Server   ServerConfig   `json:"server"`
	Walhalla WalhallaConfig `json:"walhalla"`
	Database DatabaseConfig `json:"database"`
}

//-----------------|

func (mc *MainConfig) ReadFromFile(filePath string) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	return mc.UnmarshalJSON(data)
}
