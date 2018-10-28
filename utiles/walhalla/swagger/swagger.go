package swagger

import (
	"strings"

	"gopkg.in/yaml.v2"
)

func getSubcategory(tags []string) string {
	if len(tags) == 1 {
		return tags[0]
	}
	return ""
}

func getAPI(title string) string {
	tokens := strings.Split(title, " ")
	for i, token := range tokens {
		tokens[i] = UpLieader(token)
	}
	return strings.Join(tokens, ``) + `API`
}

func ParceSwaggerYaml(data []byte) (res ParsedData) {
	var (
		subcategories = map[string]bool{}
		doc           = document{}
	)
	res.Sub2Operation = map[string][]string{}

	yaml.UnmarshalStrict(data, &doc)
	for _, path := range doc.Paths {
		for _, route := range path {
			var (
				subcategory = getSubcategory(route.Tags)
				operation   = UpLieader(route.OperationID)
			)
			subcategories[subcategory] = true

			res.Sub2Operation[subcategory] = append(res.Sub2Operation[subcategory], operation)
			res.Operations = append(res.Operations, Operation{
				Subcategory: subcategory,
				OperationID: operation,
				Parametr:    operation + "Params",
				Function:    operation + "HandlerFunc",
				Handler:     UpLieader(subcategory) + operation + "Handler",
			})
		}
	}
	res.Info = doc.Info
	res.API = getAPI(res.Info.Title)

	for subcategory := range subcategories {
		res.Subcategories = append(res.Subcategories, subcategory)
	}

	return res
}

// ----------------| misc

func UpLieader(str string) string {
	if str != "" {
		str = strings.ToUpper(str[:1]) + str[1:]
	}
	return str
}
