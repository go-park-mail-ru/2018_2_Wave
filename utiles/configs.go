package utiles

import "io/ioutil"

// CORSConfig configuration
// easyjson:json
type CORSConfig struct {
	Origin      string `json:"origin"`
	Credentials string `json:"credentials"`
	Methods     string `json:"methods"`
	Headers     string `json:"headers"`
}

// ServerConfig configuration
// easyjson:json
type ServerConfig struct {
	Port string `json:"port"`
}

// DatabaseConfig configuration
// easyjson:json
type DatabaseConfig struct {
	User 		string `json:"user"`
	Password	string `json:"password"`
	Name		string `json:"name"`
}

// WalhallaConfig configuration
// easyjson:json
type WalhallaConfig struct {
	HeaderTmpl string `json:"headerTmpl"`
	HandleTmpl string `json:"handleTmpl"`
	StructTmpl string `json:"structTmpl"`
	RouterTmpl string `json:"routerTmpl"`
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
