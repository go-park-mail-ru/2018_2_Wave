package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path"
)

func parseFile(filename string, stat *statistics, unparsed map[string]bool) {
	var (
		fileSet   = token.NewFileSet()
		node, err = parser.ParseFile(fileSet, filename, nil, parser.ParseComments)
	)
	check(err)

	for _, f := range node.Comments {
		if isAppTarget(f) {
			parseAppComment(f, stat)
		}
		if isPackTarget(f) {
			parsePackComment(f, stat)
		}
		if isFileTarget(f) {
			parseFileComment(f, stat)
		}
	}

	for _, n := range node.Decls {
		if f, ok := n.(*ast.FuncDecl); ok {
			if !isGenerationTarget(f.Doc) {
				continue
			}
			name := f.Name.Name
			if !unparsed[name] {
				dir := path.Dir(name)
				check(fmt.Errorf("Unexpected function name %s in package %s", name, dir))
			}
			parseFunction(f, stat)
			unparsed[name] = false
		}
	}
}

func parsePackage(packagePath, packageName string, stat *statistics, operations []string) {
	files, err := extractFiles(packagePath)
	if err != nil {
		fmt.Printf(" -- Note: missing package: %s \n", packagePath)
		return
	}
	fmt.Printf("Package: %s\n", packagePath)
	stat.setPackageActive(packageName)

	unparsed := map[string]bool{}
	{ // parse files
		for _, operation := range operations {
			unparsed[operation] = true
		}
		for _, file := range files {
			parseFile(file, stat, unparsed)
		}
	}
	{ // print missing operations
		missing := []string{}
		for op, miss := range unparsed {
			if miss {
				missing = append(missing, op)
			}
		}
		if len(missing) > 0 {
			fmt.Printf(" -- Note: missing operations(%s): %v\n", packageName, missing)
		}
	}
}
