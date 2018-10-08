package main

import (
	"errors"
	"fmt"
	"go/ast"
	"strings"
)

func getFieldType(field *ast.Field) string {
	if tp, ok := field.Type.(*ast.Ident); ok {
		return tp.Name
	}

	if tp, ok := field.Type.(*ast.ArrayType); ok {
		if etp, ok := tp.Elt.(*ast.Ident); ok {
			if etp.Name == "byte" {
				return "binary"
			}
		}
		exitIfFatal(errors.New("Invalid slice element. Only []byte are allowed"))
	}
	exitIfFatal(fmt.Errorf("Unexpected field type at %s", field.Names[0].Name))
	return ""
}

func parsFieldTags(jsonTag, walhallaTag string, rule *fieldRule) (bParse bool) {
	if jsonTag == "-" {
		return false
	}

	// remove everything excluding the field's name
	if comaIdx := strings.Index(jsonTag, ","); comaIdx != -1 {
		jsonTag = jsonTag[:comaIdx]
	}
	if jsonTag == "" {
		rule.FieldAlias = rule.FieldName
	} else {
		rule.FieldAlias = jsonTag
	}

	var (
		walhallaJSON = upgardeToJSON(walhallaTag)
		bytes        = []byte(walhallaJSON)
	)
	exitIfFatal(rule.UnmarshalJSON(bytes))

	return true
}

func parseStruct(outGen *outBuilder, g *ast.GenDecl, currStruct *ast.StructType, tp *ast.TypeSpec) {
	rules := structRule{
		StructName: tp.Name.Name,
	}

	for _, field := range currStruct.Fields.List {
		rule := fieldRule{
			FieldName: field.Names[0].Name,
		}

		if field.Tag != nil {
			var (
				tag         = makeTag(field.Tag)
				jsonTag     = cleanupSpaces(tag.Get("json"))
				walhallaTag = cleanupSpaces(tag.Get("walhalla"))
			)
			if !parsFieldTags(jsonTag, walhallaTag, &rule) {
				continue
			}
		}

		rule.Type = getFieldType(field)
		rules.Fields = append(rules.Fields, rule)
	}

	outGen.Structs = append(outGen.Structs, rules)
}
