package main

import (
	"Wave/utiles"
	"html/template"
)

// easyjson:json
type functionRules struct {
	FunctionName string
	URI          string // /exp/:catch/:param
	Method       string // GET \ POST \ HEAD...
	Data         string // form \ json \ uri
	Target       string // IncomeObjType \ ""
	Validation   string // true \ false
	Auth         string // true \ false
}

// easyjson:json
type fieldRule struct {
	FieldName  string `json:"-"`
	FieldAlias string `json:"-"`   // name in incomming data
	Type       string `json:"-"`   //
	Min        string `json:"min"` // minimal value \ lenght
	Max        string `json:"max"` // maximal value \ lenght
}

type structRule struct {
	StructName string
	Fields     []fieldRule
}

type outBuilder struct {
	PackageName string
	Handlers    []functionRules
	Structs     []structRule
	Config      utiles.MainConfig
}

type templateBank struct {
	Header *template.Template
	Handle *template.Template
	Struct *template.Template
	Router *template.Template
}

//-----------------|

func (fr *functionRules) Validate() bool {
	// TODO:: add more correct validation
	return fr.FunctionName != "" &&
		fr.Method != "" &&
		fr.URI != ""
}
