package main

import (
	"go/ast"
	"log"
	"os"
	"reflect"
	"strings"
)

func exitIfFatal(err error) {
	if err != nil {
		panic(err)
	}
}

//--------------------| package

func extractDir(name string) *os.File {
	dir, err := os.Open(name)
	exitIfFatal(err)

	stat, err := dir.Stat()
	exitIfFatal(err)
	if !stat.IsDir() {
		log.Fatalf("%s isn't a directory", name)
	}
	return dir
}

func extractName(root *os.File) string {
	return root.Name()
}

func extractFileNames(dir *os.File) (names []string) {
	infos, err := dir.Readdir(-1)
	exitIfFatal(err)

	rootName := extractName(dir)
	for _, info := range infos {
		name := info.Name()
		if !info.IsDir() &&
			name[len(name)-3:] == ".go" &&
			name[len(name)-7:] != ".gen.go" &&
			name[len(name)-7:] != ".tmp.go" {
			names = append(names, rootName+info.Name())
		}
	}
	return names
}

//--------------------| tag parsing

func makeTag(tag *ast.BasicLit) reflect.StructTag {
	return reflect.StructTag(tag.Value[1 : len(tag.Value)-1])
}

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
			key, val = pair, "true"
		} else {
			key, val = pair[:idx], pair[idx+1:]
		}
		pairs = append(pairs, `"`+key+`": "`+val+`"`)
	}
	return `{` + strings.Join(pairs, ", ") + `}`
}

//--------------------| doc parsing

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
