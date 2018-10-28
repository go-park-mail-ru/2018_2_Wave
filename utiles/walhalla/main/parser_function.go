package main

import (
	"go/ast"
)

func parseFunction(f *ast.FuncDecl, stat *statistics) {
	var (
		rules    = extractRules(f.Doc, anchor)
		json     = upgardeToJSON(rules)
		bytes    = []byte(json)
		settings = funcSettings{
			Name: f.Name.Name,
			Auth: "true",
		}
	)
	check(settings.UnmarshalJSON(bytes))
	check(settings.Validate())
	stat.addFunction(settings)
}
