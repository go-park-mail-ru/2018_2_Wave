package main

import (
	"bytes"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"Wave/utiles/walhalla/swagger"
)

const (
	swaggerFile = "swagger.yml"
	generated   = "generated/"
	restip      = generated + "restapi/"
)

func getOperations(root *os.File) swagger.ParsedData {
	data, err := ioutil.ReadFile(root.Name() + `/` + swaggerFile)
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

func walhalla(tmpl *template.Template, root *os.File) {
	var (
		// version
		versionFileName = `./` + generated + "walhalla.version"
		versionBytes, _ = ioutil.ReadFile(versionFileName)
		version         = string(versionBytes)

		// files
		swagger         = getOperations(root)
		files           = extractFileNames(root, validateSubcategories(swagger.Subcategories))
		outName         = makeConfigureFileName(swagger.Info.Title)
		prefix, project = extractPrefixAndProjetcName(root)

		// misc
		buffer = &bytes.Buffer{}
		stat   = &statistics{
			Operations:    swagger.Operations,
			Info:          swagger.Info,
			API:           swagger.API,
			Subcategories: swagger.Subcategories,
			Project:       project,
		}
	)
	{ // api version
		println(versionFileName, `|`, stat.Info.Version, `|`)
	}
	{ // parse files
		for _, file := range files {
			parseFile(file, stat)
		}
		stat.build()
	}
	{ // generate template
		check(tmpl.Execute(buffer, stat))
	}
	// call swagger
	if version != stat.Info.Version {
		println(" -- swagger generator started")
		defer println(" -- swagger generator stoped")

		os.Mkdir(generated, os.ModeDir)
		var (
			flagTarget = prefix + generated
			flagSpec   = prefix + swaggerFile
			appDepth   = strings.Count(prefix, `/`)
			dir        = strings.Repeat(`../`, appDepth) + `.`
			cmd        = exec.Command("swagger", "generate", "server", "--target", flagTarget, "--spec", flagSpec)
			stdout, _  = cmd.StdoutPipe()
			stderr, _  = cmd.StderrPipe()
		)
		cmd.Dir = dir

		go func() { io.Copy(os.Stdout, stdout) }()
		go func() { io.Copy(os.Stderr, stderr) }()
		check(cmd.Run())
	} else {
		println(" -- swagger files are up to data")
	}
	{ // change configuration file
		check(ioutil.WriteFile(outName, buffer.Bytes(), os.ModeExclusive))
	}
	{ // write version file
		file, err := os.Create(versionFileName)
		check(err)
		_, err = file.WriteString(stat.Info.Version)
		check(err)
	}
	{ // format
		check(exec.Command("go", "fmt", `./`+restip).Run())
	}
}
