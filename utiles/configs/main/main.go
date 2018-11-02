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
	confPath := ""
	if flag.Parse(); flag.NArg() != 0 {
		confPath = flag.Arg(0)
	}

	mainConfig := configs.MainConfig{}
	{ // read config
		data, err := ioutil.ReadFile(confPath)
		if err != nil {
			data = []byte("{}")
		}
		if err := mainConfig.UnmarshalJSON(data); err != nil {
			fmt.Printf("Incorrect config json: %s\nError: %v", data, err)
		}
	}
	out := bytes.Buffer{}
	{ // write config
		if data, err := mainConfig.MarshalJSON(); err != nil {
			fmt.Printf("Error during config generation: %v\n", err)
		} else {
			if err := json.Indent(&out, data, "", "  "); err != nil {
				fmt.Printf("Error during formating: %v\n", err)
			}
			out.WriteString("\n")
		}
	}
	ioutil.WriteFile(confPath, out.Bytes(), os.ModeExclusive)
}
