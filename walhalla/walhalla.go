package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
)

func walhalla(out *os.File, tmpl tmpls, filenames []string) {
	outGen := &outBuilder{}

	for _, filename := range filenames {
		var (
			fset      = token.NewFileSet()
			node, err = parser.ParseFile(fset, filename, nil, parser.ParseComments)
		)
		println(filename)
		exitIfFatal(err)

		outGen.PackageName = node.Name.Name

		for _, f := range node.Decls {
			if g, ok := f.(*ast.GenDecl); ok {
				if !isGenerationTarget(g.Doc) {
					continue
				}

				for _, spec := range g.Specs {
					if currType, ok := spec.(*ast.TypeSpec); !ok {
						continue
					} else if currStruct, ok := currType.Type.(*ast.StructType); ok {
						parseStruct(outGen, g, currStruct, currType)
					}
				}
			} else if g, ok := f.(*ast.FuncDecl); ok {
				if !isGenerationTarget(g.Doc) {
					continue
				}
				parseFunc(outGen, g)
			}
		}
	}

	tmpl.header.Execute(out, outGen)
	tmpl.structs.Execute(out, outGen)
	tmpl.handlers.Execute(out, outGen)
	tmpl.router.Execute(out, outGen)
}
