package main

import (
	"Wave/utiles"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func parsFile(outGen *outBuilder, filename string) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
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

func walhalla(out io.Writer, tmpl templateBank, rootName string, config utiles.MainConfig, filenames []string) {
	// parse files
	outGen := &outBuilder{
		Config: config,
	}
	for _, filename := range filenames {
		parsFile(outGen, filename)
	}

	// generte templates
	exitIfFatal(tmpl.Header.Execute(out, outGen))
	exitIfFatal(tmpl.Struct.Execute(out, outGen))
	exitIfFatal(tmpl.Handle.Execute(out, outGen))
	exitIfFatal(tmpl.Router.Execute(out, outGen))

	// insert easyjson output
	if len(outGen.Structs) > 0 {
		// call easyjson
		jsonFile := rootName + outGen.PackageName + ".tmp.go"
		exitIfFatal(exec.Command("easyjson", "-output_filename", jsonFile, "-pkg", rootName).Run())

		// move the result to the out
		if data, err := ioutil.ReadFile(jsonFile); err == nil {
			str := string(data)
			if endOfImport := strings.Index(str, ")"); endOfImport != -1 {
				str = str[endOfImport+1:]
			}
			out.Write([]byte(str))
			os.Remove(jsonFile)
		}
	}
}
