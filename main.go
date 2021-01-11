package main

import (
	"fmt"
	genshinapi "github.com/darylhjd/genshin-api"
	"html/template"
	"log"
	"net/http"
	"strings"
)

var dataTypes = []string{
	genshinapi.ArtifactsDType, genshinapi.CharactersDType, genshinapi.DomainsDType,
	genshinapi.ElementsDType, genshinapi.NationsDType, genshinapi.WeaponsDType,
}

// indexTemplate : Template for index page.
var indexTemplate = template.Must(template.ParseFiles(
	"static/templates/boilerplate.html",
	"static/templates/index.html"))

// dataTypeTemplate : Template for each data type
var dataTypeTemplate *template.Template

// dataEntryTemplate : Template for each data entry
var dataEntryTemplate *template.Template

func main() {
	mux := http.NewServeMux()

	// Register templates
	// Template for index page is already registered
	RegisterDataTypeTemplate(&dataTypeTemplate) // Register template for data type
	RegisterDataEntryTemplate(&dataEntryTemplate) // Register template for data entry

	// Register handlers
	/* This is for serving static files
	mux.Handle("/static/styling/",
		http.StripPrefix("/static/styling/", http.FileServer(http.Dir("static/styling"))))
	 */
	mux.HandleFunc("/", IndexHandler) // index page handler
	RegisterDataTypeHandler(mux)

	log.Fatal(http.ListenAndServe(":8080", mux))
}

// RegisterDataTypeTemplate : Register HTML template for data type name list pages.
func RegisterDataTypeTemplate(t **template.Template) {
	files := []string{
		"static/templates/boilerplate.html",
	}
	for _, dtype := range dataTypes {
		files = append(files, fmt.Sprintf("static/templates/datatypes/%s.html", dtype))
	}

	*t = template.Must(template.ParseFiles(files...))
}

// RegisterDataEntryTemplate : Register HTML template for data entry for each data type.
func RegisterDataEntryTemplate(t **template.Template) {
	files := []string{
		"static/templates/boilerplate.html",
	}
	for _, dtype := range dataTypes {
		files = append(files,
			fmt.Sprintf("static/templates/datatype/%s.html",
				strings.TrimSuffix(dtype, "s")))
	}

	*t = template.Must(template.ParseFiles(files...))
}

// RegisterDataTypeHandler : Register handlers for each data type.
func RegisterDataTypeHandler(m *http.ServeMux) {
	for _, dtype := range dataTypes {
		m.HandleFunc(MakeDataTypeHandler(dtype))
	}
}
