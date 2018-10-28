package main

import (
	"fmt"
	"go/ast"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

//--------------------| package

func extractDir(name string) *os.File {
	dir, err := os.Open(name)
	check(err)

	stat, err := dir.Stat()
	check(err)
	if !stat.IsDir() {
		check(fmt.Errorf("%s isn't a directory", name))
	}
	return dir
}

func validateSubcategories(subcategories []string) []string {
	for _, subcategory := range subcategories {
		if subcategory == "" {
			return subcategories
		}
	}
	return append(subcategories, "")
}

func extractFileNamesInternal(dir *os.File) (names []string) {
	infos, err := dir.Readdir(-1)
	check(err)

	rootName := dir.Name() + `/`
	for _, info := range infos {
		name := info.Name()
		if info.IsDir() {
			continue
		}
		if len(name) >= 3 && name[len(name)-3:] != ".go" ||
			len(name) >= 7 && name[len(name)-7:] == ".gen.go" ||
			len(name) >= 7 && name[len(name)-7:] == ".tmp.go" {
			continue
		}
		names = append(names, rootName+info.Name())
	}
	return names
}

func extractFileNames(dir *os.File, subcategories []string) (names []string) {
	for _, subcategory := range subcategories {
		internal, err := os.Open(dir.Name() + `/` + subcategory)
		if err != nil {
			continue
		}

		st, err := dir.Stat()
		if check(err); !st.IsDir() {
			check(fmt.Errorf("%s isn't a directory", dir.Name()))
		}

		names = append(names, extractFileNamesInternal(internal)...)
	}
	return names
}

func extractProjectName(data []byte) string {
	str := string(data)
	bgn := strings.Index(str, " ")
	end := strings.Index(str, "\n")
	return str[bgn+1 : end]
}

func extractPrefixAndProjetcName(dir *os.File) (prefix, project string) {
	var (
		abs, err   = filepath.Abs(dir.Name())
		targetPath = abs
	)
	check(err)

	for ; ; abs = filepath.Clean(abs + `/..`) {
		var (
			names, err = ioutil.ReadDir(abs)
		)
		check(err)

		for _, name := range names {
			if name.Name() == "go.mod" {
				var (
					fullName   = abs + `\` + name.Name()
					data, err1 = ioutil.ReadFile(fullName)
					pref, err2 = filepath.Rel(abs, targetPath)
				)
				check(err1)
				check(err2)
				return pref + `/`, extractProjectName(data)
			}
		}
	}
}

//--------------------| tag parsing

func makeTag(tag *ast.BasicLit) reflect.StructTag {
	return reflect.StructTag(tag.Value[1 : len(tag.Value)-1])
}

func cleanupSpaces(rule string) string {
	for _, trg := range []string{" ", "\t", "\n"} {
		rule = strings.Replace(rule, trg, "", -1)
	}
	return rule
}

var (
	jsonSplitter, _      = regexp.Compile(`(\w+)` + `:?` + `(\w*)?` + `(\[[\w,]*\])?`)
	jsonArraySplitter, _ = regexp.Compile(`(\w+)`)
)

const (
	groupKey = 1
	groupVal = 2
	groupArr = 3
)

func upgardeToJSONArray(arr string) string {
	toks := []string{}
	for _, submache := range jsonArraySplitter.FindAllStringSubmatch(arr, -1) {
		toks = append(toks, `"`+submache[groupKey]+`"`)
	}
	return `[` + strings.Join(toks, `, `) + `]`
}

func upgardeToJSON(rule string) (json string) {
	rule = cleanupSpaces(rule)
	pairs := []string{}
	for _, submaches := range jsonSplitter.FindAllStringSubmatch(rule, -1) {
		var (
			key = submaches[groupKey]
			val = submaches[groupVal]
			arr = submaches[groupArr]
		)
		if val == "" && arr == "" {
			pairs = append(pairs, `"`+key+`": "true"`)
			continue
		}
		if val != "" {
			pairs = append(pairs, `"`+key+`": "`+val+`"`)
			continue
		}
		if arr != "" {
			arr = upgardeToJSONArray(arr)
			pairs = append(pairs, `"`+key+`": `+arr)
			continue
		}
		check(fmt.Errorf("unexpected rule: %s", submaches[0]))
	}
	return `{` + strings.Join(pairs, ", ") + `}`
}

//--------------------| doc parsing

const (
	anchor  = "walhalla:gen"
	fileTag = "walhalla:file"
	appTag  = "walhalla:app"
)

func isTarget(doc *ast.CommentGroup, tag string) bool {
	if doc == nil {
		return false
	}
	return strings.Contains(doc.Text(), tag)
}
func isGenerationTarget(doc *ast.CommentGroup) bool { return isTarget(doc, anchor) }
func isFileTarget(doc *ast.CommentGroup) bool       { return isTarget(doc, fileTag) }
func isAppTarget(doc *ast.CommentGroup) bool        { return isTarget(doc, appTag) }

func extractRules(doc *ast.CommentGroup, tag string) string {
	var (
		text  = doc.Text()
		index = strings.Index(text, tag)
		rules = text[index:]
	)
	if len(rules) > len(tag) {
		rules = rules[len(tag):]
	}
	// nested types are not supported
	index1 := strings.Index(rules, "{")
	index2 := strings.Index(rules, "}")
	if index1 == -1 || index2 == -1 {
		return ""
	}
	return rules[index1+1 : index2]
}
