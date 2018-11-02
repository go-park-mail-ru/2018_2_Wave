package main

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"

	"Wave/utiles/walhalla/swagger"
)

const (
	swaggerFile = "swagger.yml"
	generated   = "generated/"
	restip      = generated + "restapi/"
)

func getOperations(root string) swagger.ParsedData {
	var (
		confPath  = path.Join(root, swaggerFile)
		data, err = ioutil.ReadFile(confPath)
	)
	check(err)
	return swagger.ParceSwaggerYaml(data)
}

func makeConfigureFileName(title string) string {
	lower := strings.ToLower(title)
	tokens := strings.Split(lower, ` `)
	for i, tk := range tokens {
		if tk != "" {
			continue
		}
		tokens = append(tokens[:i], tokens[i+1:]...)
	}
	return restip + `configure_` + strings.Join(tokens, `_`) + `.go`
}

func walhalla(tmpl *template.Template, root string) {
	var (
		// version
		versionFileName = path.Join(root, generated, "walhalla.version")
		versionBytes, _ = ioutil.ReadFile(versionFileName)
		version         = string(versionBytes)

		// enviroment
		swagger         = getOperations(root)
		outName         = makeConfigureFileName(swagger.Info.Title)
		prefix, project = extractPrefixAndProjectName(root)

		// misc
		buffer = &bytes.Buffer{}
		stat   = &statistics{
			Operations:    makeOperations(swagger.Operations),
			API:           swagger.API,
			Project:       project,
			Application:   path.Join(project, prefix),
			Subcategories: swagger.Subcategories,
		}
	)
	{ // api version
		println(versionFileName, `|`, swagger.Info.Version, `|`)
	}
	{ // pars packages
		for _, sb := range append(swagger.Subcategories, "") {
			var (
				path       = path.Join(root, sb)
				operations = swagger.Sub2Operation[sb]
			)
			parsePackage(path, sb, stat, operations)
		}
		stat.build()
	}
	{ // generate template
		check(tmpl.Execute(buffer, stat))
	}
	// call swagger
	if version != swagger.Info.Version {
		println(" -- swagger generator started")
		defer println(" -- swagger generator stoped")

		os.Mkdir(generated, 0755)
		var (
			flagTarget = prefix + generated
			flagSpec   = prefix + swaggerFile
			appDepth   = strings.Count(prefix, `/`)
			dir        = strings.Repeat(`../`, appDepth) + `.`
			cmd        = exec.Command("swagger", "generate", "server", "--target", flagTarget, "--spec", flagSpec)
		)
		cmd.Dir = dir
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		check(cmd.Run())
	} else {
		println(" -- swagger files are up to date")
	}
	{ // change configuration file
		check(ioutil.WriteFile(outName, buffer.Bytes(), os.ModeExclusive))
	}
	{ // write version file
		file, err := os.Create(versionFileName)
		check(err)
		file.WriteString(swagger.Info.Version)
	}
	{ // format generated file
		cmd := exec.Command("go", "fmt", `./`+restip)
		cmd.Stderr = os.Stderr
		cmd.Run()
	}
}
