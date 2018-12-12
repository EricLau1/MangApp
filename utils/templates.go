package utils

import (
	"html/template"
	"net/http"
)

var templates *template.Template

func LoadTemplates(dir string) {

	templates = template.Must(template.ParseGlob(dir))

}

func ExecuteTemplate(w http.ResponseWriter, tmpl string, dados interface{} ) {

	templates.ExecuteTemplate(w, tmpl, dados)

}