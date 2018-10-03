package main

import (
	"encoding/json"
	"errors"
	"go/ast"
)

func parseFunc(outGen *outBuilder, g *ast.FuncDecl) {
	var (
		rulesText  = extractRules(g.Doc)
		rulesJSON  = upgardeToJSON(rulesText)
		rulesBytes = []byte(rulesJSON)
		rule       = functionRules{
			FunctionName: g.Name.Name,
		}
	)
	err := json.Unmarshal(rulesBytes, &rule)
	exitIfFatal(err)

	if !rule.Validate() {
		exitIfFatal(errors.New("Invalid api params"))
	}

	outGen.Handlers = append(outGen.Handlers, rule)
}
