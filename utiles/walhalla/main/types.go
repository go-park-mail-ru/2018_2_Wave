package main

import (
	"Wave/utiles/walhalla/swagger"

	"github.com/asaskevich/govalidator"
)

// ----------------| statistics

// easyjson:json
type statistics struct {
	Project     string
	Application string
	API         string
	Operations  []operation

	Subcategories     []string
	ImplSubcategories []string
	FuncSettings      []funcSettings
	GlobalMiddlewares []string
	BMiddlewares      bool

	buildContext buildContext
}

type operation struct {
	swagger.Operation
	Implemented bool
}

func makeOperations(ops []swagger.Operation) (res []operation) {
	for _, op := range ops {
		res = append(res, operation{
			Operation:   op,
			Implemented: false,
		})
	}
	return res
}

// ------| settings

func (st *statistics) pushSettings(fs funcSettings) {
	st.buildContext.pushFuncRule(fs)
}

func (st *statistics) popSettings() {
	st.buildContext.popFuncRule()
}

func (st *statistics) addFunction(fs funcSettings) {
	{
		st.buildContext.pushFuncRule(fs)
		fs = st.buildContext.buildFuncRule()
		st.buildContext.popFuncRule()
	}
	st.FuncSettings = append(st.FuncSettings, fs)
	{
		for i, op := range st.Operations {
			if op.OperationID != fs.Name {
				continue
			}
			st.Operations[i].Implemented = true
		}
	}
}

func (st *statistics) setAppSettings(as appSettings) {
	st.GlobalMiddlewares = append(as.GlobalMiddlewares, st.GlobalMiddlewares...)
	st.pushSettings(as.ExtractFuncSettings())
}

// ------| packages

func (st *statistics) setPackageActive(pack string) {
	if pack == "" {
		return
	}
	st.ImplSubcategories = append(st.ImplSubcategories, pack)
}

// ------|

func (st *statistics) build() {
	{ // set default package for operations with empty subcategories
		for i, op := range st.Operations {
			if op.Subcategory == "" {
				op.Subcategory = "operations"
			}
			st.Operations[i] = op
		}
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
			{ // set Subcategory, Function, Parametr
				fs.Subcategory = ID2Cat[fs.Name]
				fs.Function = ID2Func[fs.Name]
				fs.Parametr = ID2Param[fs.Name]
				fs.Package = fs.Subcategory

				if fs.Package == "operations" {
					fs.Package = ""
				}
			}
			{ // middlewares
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
	{
		st.BMiddlewares = st.BMiddlewares || len(st.GlobalMiddlewares) > 0
	}
}

// ----------------| buildContext

type buildContext struct {
	funcRuleStack []funcSettings
}

// ------|

func (bc *buildContext) pushFuncRule(fs funcSettings) {
	bc.funcRuleStack = append(bc.funcRuleStack, fs)
}

func (bc *buildContext) popFuncRule() {
	bc.funcRuleStack = bc.funcRuleStack[:len(bc.funcRuleStack)-1]
}

func (bc *buildContext) buildFuncRule() funcSettings {
	res := funcSettings{}
	for _, fs := range bc.funcRuleStack {
		fs.Override(&res)
	}
	return res
}

// ----------------| appSettings

// easyjson:json
type appSettings struct {
	GlobalMiddlewares   []string `json:"globalMiddlewares"`   // list of middleware names in a direct order
	OperationMiddlewars []string `json:"operationMiddlewars"` // list of middleware names in a direct order; will be applyed to each operation
}

// ------|

func (as *appSettings) Validate() error {
	_, err := govalidator.ValidateStruct(as)
	return err
}

func (as *appSettings) ExtractFuncSettings() funcSettings {
	return funcSettings{
		Middlewares: as.OperationMiddlewars,
	}
}

// ----------------| appSettings

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

func override(target *string, src string) {
	if src != "" {
		*target = src
	}
}

// ------|

func (fs *funcSettings) Validate() error {
	_, err := govalidator.ValidateStruct(fs)
	return err
}

func (fs *funcSettings) Override(target *funcSettings) {
	override(&target.Name, fs.Name)
	override(&target.Auth, fs.Auth)
	override(&target.Model, fs.Model)
	override(&target.Package, fs.Package)
	override(&target.Parametr, fs.Parametr)
	override(&target.Function, fs.Function)
	override(&target.Subcategory, fs.Subcategory)
	target.Middlewares = append(target.Middlewares, fs.Middlewares...)
}
