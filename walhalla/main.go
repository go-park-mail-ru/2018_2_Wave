package main

import (
	"html/template"
	"os"
)

// usage (from the project root): go run ./walhalla <package>
func main() {
	var (
		tmpl     = tmpls{}
		root     = extractPackageDir()
		rootName = extractRootName(root)
		out, err = os.Create(rootName + "walhalla.gen.go")
	)
	exitIfFatal(err)

	const debug = false
	if debug {
		tmpl.header = template.Must(template.ParseFiles("./tmpl_header.txt"))
		tmpl.handlers = template.Must(template.ParseFiles("./tmpl_handlers.txt"))
		tmpl.structs = template.Must(template.ParseFiles("./tmpl_structs.txt"))
		tmpl.router = template.Must(template.ParseFiles("./tmpl_router.txt"))
	} else {
		tmpl.header = template.Must(template.ParseFiles("./walhalla/tmpl_header.txt"))
		tmpl.handlers = template.Must(template.ParseFiles("./walhalla/tmpl_handlers.txt"))
		tmpl.structs = template.Must(template.ParseFiles("./walhalla/tmpl_structs.txt"))
		tmpl.router = template.Must(template.ParseFiles("./walhalla/tmpl_router.txt"))
	}
	walhalla(out, tmpl, extractPackageFileNames(root))
}
