package main

import (
	"Wave/utiles"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

const confPath = "./resources/configs/main.json"

// from project root: go run /utiles/configs .
func main() {
	data, err := ioutil.ReadFile(confPath)
	if err != nil {
		data = []byte("{}")
	}

	mainConfig := utiles.MainConfig{}

	if err := mainConfig.UnmarshalJSON(data); err != nil {
		fmt.Printf("Incorrect config json: %s\nError: %v", data, err)
	}

	if data, err = mainConfig.MarshalJSON(); err != nil {
		fmt.Printf("Error during config generation: %v\n", err)
	}

	out := bytes.Buffer{}
	if err := json.Indent(&out, data, "", "  "); err != nil {
		fmt.Printf("Error during formatin: %v\n", err)
	}
	out.WriteString("\n")

	ioutil.WriteFile(confPath, out.Bytes(), os.ModeExclusive)
}
