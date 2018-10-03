package main

import (
	"encoding/json"
	"errors"
	"go/ast"
	"reflect"
)

func parseStruct(outGen *outBuilder, g *ast.GenDecl, currStruct *ast.StructType, tp *ast.TypeSpec) {
	rules := structRules{
		StructName: tp.Name.Name,
	}

	for _, field := range currStruct.Fields.List {
		rule := fieldRules{
			FieldName: field.Names[0].Name,
		}

		if field.Tag != nil {
			var (
				tag          = reflect.StructTag(field.Tag.Value[1 : len(field.Tag.Value)-1])
				jsonTag      = cleanupSpaces(tag.Get("json"))
				walhallaTag  = tag.Get("walhalla")
				walhallaJSON = upgardeToJSON(walhallaTag)
				bytes        = []byte(walhallaJSON)
			)
			if jsonTag == "-" {
				continue
			}
			if jsonTag == "" {
				rule.FieldAlias = rule.FieldName
			} else {
				rule.FieldAlias = jsonTag
			}

			err := json.Unmarshal(bytes, &rule)
			exitIfFatal(err)
		}

		if tp, ok := field.Type.(*ast.Ident); ok {
			rule.Type = tp.Name

		} else if tp, ok := field.Type.(*ast.ArrayType); ok {
			typeError := errors.New("Invalid slice element. Only []byte are allowed")

			if etp, ok := tp.Elt.(*ast.Ident); ok {
				if etp.Name != "byte" {
					exitIfFatal(typeError)
				}
				rule.Type = "binary"
			} else {
				exitIfFatal(typeError)
			}
		}
		rules.Fields = append(rules.Fields, rule)
	}
	outGen.Structs = append(outGen.Structs, rules)
}
