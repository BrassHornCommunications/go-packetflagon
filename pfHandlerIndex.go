package main

import (
	"html/template"
	"log"
	"net/http"
)

func PFWebIndex(w http.ResponseWriter, r *http.Request, templateData TemplateConf) {

	tmpl, err := template.New("index").ParseFiles("assets/templates/index.html")

	err = tmpl.Execute(w, templateData)

	if err != nil {
		log.Fatal(err)
	}
}
