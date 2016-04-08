package main

import (
	"html/template"
	"log"
	"net/http"
)

func PFWebAbout(w http.ResponseWriter, r *http.Request, templateData TemplateConf) {
	tmpl, err := template.New("about").ParseFiles("assets/templates/about.html")

	err = tmpl.Execute(w, templateData)

	if err != nil {
		log.Fatal(err)
	}

}
