package swagger

type route struct {
	OperationID string   `yaml:"operationId"`
	Tags        []string `yaml:"tags"`
}

type Info struct {
	Version string `yaml:"version"`
	Title   string `yaml:"title"`
}

type document struct {
	Info  Info                        `yaml:"info"`
	Paths map[string]map[string]route `yaml:"paths"`
}

// Operation data
type Operation struct {
	OperationID string
	Subcategory string
	Handler     string
	Function    string
	Parametr    string
}

type ParsedData struct {
	Info          Info
	API           string
	Subcategories []string
	Operations    []Operation
	Sub2Operation map[string][]string
}
