package main

import (
	"fmt"

	"Wave/utiles/walhalla/swagger"
	"github.com/asaskevich/govalidator"
)

// ----------------| statistics

// easyjson:json
type statistics struct {
	PackageName string

	Project    string
	API        string
	Info       swagger.Info
	Operations []swagger.Operation

	Subcategories     []string
	FuncSettings      []funcSettings
	GlobalMiddlewares []string
	BMiddlewares      bool

	appSettings  appSettings
	fileSettings *fileSettings
}

// ------|

func (st *statistics) setFileSettings(fileSettings *fileSettings) {
	st.fileSettings = fileSettings
}

func (st *statistics) setAppSettings(appSettings *appSettings) {
	st.appSettings = *appSettings
}

func (st *statistics) beginFileParsing() {
	st.fileSettings = nil
}

func (st *statistics) addFunction(fs funcSettings) {
	if fs.Model == `` && st.fileSettings != nil {
		fs.Model = st.fileSettings.Model
	}
	if fs.Model == `-` {
		fs.Model = ``
	}
	st.FuncSettings = append(st.FuncSettings, fs)
}

func (st *statistics) build() {
	{ // remove an empty subcategory
		for i, sub := range st.Subcategories {
			if sub == "" {
				st.Subcategories = append(st.Subcategories[:i], st.Subcategories[i+1:]...)
			}
		}
	}
	{ // set default package for operations with empty subcategories
		for i, op := range st.Operations {
			if op.Subcategory == "" {
				op.Subcategory = "operations"
			}
			st.Operations[i] = op
		}
	}
	{ // set global middlewares
		st.GlobalMiddlewares = st.appSettings.GlobalMiddlewares
	}
	{ // function settings
		ID2Cat := map[string]string{}   // operationID -> subcategory
		ID2Func := map[string]string{}  // operationID -> generated function name
		ID2Param := map[string]string{} // operationID -> generated parametr name
		for _, op := range st.Operations {
			ID2Cat[op.OperationID] = op.Subcategory
			ID2Func[op.OperationID] = op.Function
			ID2Param[op.OperationID] = op.Parametr
		}

		for i, fs := range st.FuncSettings {
			{ // set Subcategory && Function && Parametr
				if _, ok := ID2Cat[fs.Name]; !ok {
					check(fmt.Errorf("Unexpected handler %s", fs.Name))
				} else {
					fs.Subcategory = ID2Cat[fs.Name]
					fs.Function = ID2Func[fs.Name]
					fs.Parametr = ID2Param[fs.Name]
					fs.Package = fs.Subcategory

					if fs.Package == "operations" {
						fs.Package = ""
					}
				}
			}
			{ // middlewares
				fs.Middlewares = append(fs.Middlewares, st.appSettings.OperationMiddlewars...)
				if fs.Auth == "true" {
					fs.Middlewares = append(fs.Middlewares, "auth")
				}

				// invert
				max := len(fs.Middlewares) - 1
				for i := 0; i < len(fs.Middlewares)/2; i++ {
					fs.Middlewares[i], fs.Middlewares[max-i] = fs.Middlewares[max-i], fs.Middlewares[i]
				}
				st.BMiddlewares = st.BMiddlewares || max >= 0
			}
			st.FuncSettings[i] = fs
		}
	}
}

// ----------------| settings

// easyjson:json
type appSettings struct {
	GlobalMiddlewares   []string `json:"globalMiddlewares"`   // list of middleware names in a direct order
	OperationMiddlewars []string `json:"operationMiddlewars"` // list of middleware names in a direct order; will be applyed to each operation
}

// ------|

func (fs *appSettings) Validate() error {
	_, err := govalidator.ValidateStruct(fs)
	return err
}

// ----------------| settings

// easyjson:json
type fileSettings struct {
	Model string `json:"model"`
}

// ------|

func (fs *fileSettings) Validate() error {
	_, err := govalidator.ValidateStruct(fs)
	return err
}

// ----------------| settings

// easyjson:json
type funcSettings struct {
	Name        string `json:"name"`
	Model       string `json:"model"`
	Auth        string `json:"auth"`
	Subcategory string `json:"-"`
	Package     string `json:"-"`
	Parametr    string `json:"-"`
	Function    string `json:"-"`

	Middlewares []string `json:"mdw"`
}

// ------|

func (fs *funcSettings) Validate() error {
	_, err := govalidator.ValidateStruct(fs)
	return err
}
