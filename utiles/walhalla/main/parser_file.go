package main

import (
	"go/ast"
	"go/parser"
	"go/token"
)

func parseFile(filename string, stat *statistics) {
	println(filename)
	node, err := parser.ParseFile(token.NewFileSet(), filename, nil, parser.ParseComments)
	check(err)

	stat.beginFileParsing()
	stat.PackageName = node.Name.Name
	for _, f := range node.Comments {
		if isFileTarget(f) {
			parseFileComment(f, stat)
		}
		if isAppTarget(f) {
			parseAppComment(f, stat)
		}
	}

	for _, n := range node.Decls {
		if f, ok := n.(*ast.FuncDecl); ok {
			if !isGenerationTarget(f.Doc) {
				continue
			}
			parseFunction(f, stat)
		}
	}
}
