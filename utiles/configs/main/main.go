package main

import (
	"Wave/utiles/configs"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	confPath := "."

	flag.Parse()
	if flag.NArg() != 0 {
		confPath = flag.Arg(0)
	}

	data, err := ioutil.ReadFile(confPath)
	if err != nil {
		data = []byte("{}")
	}

	mainConfig := configs.MainConfig{}

	if err := mainConfig.UnmarshalJSON(data); err != nil {
		fmt.Printf("Incorrect config json: %s\nError: %v", data, err)
	}

	if data, err = mainConfig.MarshalJSON(); err != nil {
		fmt.Printf("Error during config generation: %v\n", err)
	}

	out := bytes.Buffer{}
	if err := json.Indent(&out, data, "", "  "); err != nil {
		fmt.Printf("Error during formating: %v\n", err)
	}
	out.WriteString("\n")

	ioutil.WriteFile(confPath, out.Bytes(), os.ModeExclusive)
}
