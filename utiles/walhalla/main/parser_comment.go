package main

import (
	"go/ast"
)

func parseFileComment(doc *ast.CommentGroup, stat *statistics) {
	var (
		rules    = extractRules(doc, fileTag)
		json     = upgardeToJSON(rules)
		bytes    = []byte(json)
		settings = &fileSettings{}
	)
	check(settings.UnmarshalJSON(bytes))
	check(settings.Validate())
	stat.setFileSettings(settings)
}

func parseAppComment(doc *ast.CommentGroup, stat *statistics) {
	var (
		rules    = extractRules(doc, appTag)
		json     = upgardeToJSON(rules)
		bytes    = []byte(json)
		settings = &appSettings{}
	)
	check(settings.UnmarshalJSON(bytes))
	check(settings.Validate())
	stat.setAppSettings(settings)
}
