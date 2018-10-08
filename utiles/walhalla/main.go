package main

import (
	"flag"
	"html/template"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"Wave/utiles"
)

var (
	confPath = flag.String("conf", "./resources/configs/main.json", " -- path to configuration")
	bDebug   = flag.Bool("debug", false, " -- in case of an inplase debug")
)

func toRoot(path *string) {
	const prefix = "./../."
	*path = prefix + *path
}

func generateTemplates(config utiles.MainConfig) (tmpl templateBank) {
	if *bDebug {
		toRoot(&config.Walhalla.HandleTmpl)
		toRoot(&config.Walhalla.HeaderTmpl)
		toRoot(&config.Walhalla.StructTmpl)
		toRoot(&config.Walhalla.RouterTmpl)
	}
	tmpl.Handle = template.Must(template.ParseFiles(config.Walhalla.HandleTmpl))
	tmpl.Struct = template.Must(template.ParseFiles(config.Walhalla.StructTmpl))
	tmpl.Header = template.Must(template.ParseFiles(config.Walhalla.HeaderTmpl))
	tmpl.Router = template.Must(template.ParseFiles(config.Walhalla.RouterTmpl))
	return tmpl
}

func loadConfig() utiles.MainConfig {
	if *bDebug {
		toRoot(confPath)
	}

	data, err := ioutil.ReadFile(*confPath)
	if err != nil {
		println("Invalid config path: ", *confPath)
		os.Exit(1)
	}

	config := utiles.MainConfig{}
	if err := config.UnmarshalJSON(data); err != nil {
		println("Invalid config")
		os.Exit(1)
	}
	return config
}

// usage (from the project root): go run ./walhalla <package>
func main() {
	flag.Parse()
	var (
		config    = loadConfig()
		templates = generateTemplates(config)
	)
	for _, pack := range flag.Args() {
		// generate code
		var (
			root     = extractDir(pack)
			rootName = extractName(root)
			files    = extractFileNames(root)
			outName  = rootName + "walhalla.gen.go"
			_        = os.Remove(outName)
		)
		buffer := strings.Builder{}
		walhalla(&buffer, templates, rootName, config, files)

		// write to a file
		out, err := os.Create(outName)
		exitIfFatal(err)
		out.WriteString(buffer.String())

		// format the file
		exec.Command("go", "fmt", outName).Run()
	}
}
