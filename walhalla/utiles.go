package main

import (
	"go/ast"
	"html/template"
	"log"
	"os"
	"strings"
)

type functionRules struct {
	FunctionName string
	URI          string `json:"uri"`
	Method       string `json:"method"`
	Data         string `json:"data"`
	Target       string `json:"target"`
	Validation   string `json:"validation"`
	Auth         string `json:"auth"`
}

type fieldRules struct {
	FieldName  string `json:"-"`
	FieldAlias string `json:"-"`
	Type       string `json:"-"`
	Min        string `json:"min"`
	Max        string `json:"max"`
}

type structRules struct {
	StructName string       `json:"-"`
	Fields     []fieldRules `json:"-"`
}

type outBuilder struct {
	PackageName string
	Handlers    []functionRules
	Structs     []structRules
}

type tmpls struct {
	header   *template.Template
	handlers *template.Template
	structs  *template.Template
	router   *template.Template
}

func (fr *functionRules) Validate() bool {
	// TODO:: add more correct validation
	return fr.FunctionName != "" &&
		fr.Method != "" &&
		fr.URI != ""
}

//--------------------|

func exitIfFatal(err error) {
	if err != nil {
		log.Print(err)
		panic(err)
	}
}

func extractPackageDir() *os.File {
	name := os.Args[1]

	dir, err := os.Open(name)
	exitIfFatal(err)

	stat, err := dir.Stat()
	exitIfFatal(err)

	if !stat.IsDir() {
		log.Fatal("Package expected")
	}
	return dir
}

func extractRootName(root *os.File) string {
	return "./" + root.Name() + "/"
}

func extractPackageFileNames(dir *os.File) (names []string) {

	infos, err := dir.Readdir(-1)
	exitIfFatal(err)

	root := extractRootName(dir)
	for _, info := range infos {
		name := info.Name()
		if !info.IsDir() &&
			name[len(name)-3:] == ".go" &&
			name[len(name)-7:] != ".gen.go" {
			names = append(names, root+info.Name())
		}
	}
	return names
}

//--------------------|

func cleanupSpaces(args string) string {
	for _, trg := range []string{" ", "\t", "\n"} {
		args = strings.Replace(args, trg, "", -1)
	}
	return args
}

func upgardeToJSON(args string) (json string) {
	args = cleanupSpaces(args)
	if len(args) == 0 {
		return `{}`
	}

	pairs := []string{}
	for _, pair := range strings.Split(args, ",") {
		var (
			idx = strings.Index(pair, ":")
			key = ""
			val = ""
		)

		if idx == -1 {
			key, val = pair, "yes"
		} else {
			key, val = pair[:idx], pair[idx+1:]
		}
		pairs = append(pairs, `"`+key+`": "`+val+`"`)
	}
	return `{` + strings.Join(pairs, ", ") + `}`
}

//--------------------|

const anchor = "walhalla:"

func isGenerationTarget(doc *ast.CommentGroup) (bNeed bool) {
	if doc == nil {
		return false
	}

	for _, comment := range doc.List {
		if strings.Contains(comment.Text, anchor) {
			return true
		}
	}
	return false
}

func extractRules(doc *ast.CommentGroup) string {
	text := doc.Text()

	index := strings.Index(text, anchor)
	rules := text[index:]

	if len(rules) > len(anchor) {
		rules = rules[len(anchor):]
	}

	index1 := strings.Index(rules, "{")
	index2 := strings.Index(rules, "}")
	if index1 == -1 || index2 == -1 {
		return ""
	}
	return rules[index1+1 : index2]
}
