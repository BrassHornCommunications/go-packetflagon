package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"html/template"
	"log"
	"net/http"
)

func PFWebCreate(w http.ResponseWriter, r *http.Request, db *bolt.DB, templateData TemplateConf) {

	if r.Method == "POST" {
		//POST indicates someone submitted data to create a new PAC

		md5Hash, err := CreatePAC(r, db)

		if err == nil {
			w.Header().Set("pf-success", "true")
			fmt.Fprintf(w, `<!DOCTYPE html><html lang=\"en\"><head><meta charset=\"utf-8\"><title>PacketFlagon | PAC Created</title></head>
			<body><h1>PAC File Successfully Created</h1>`)
			fmt.Fprintf(w, "<p>Your PAC file is <a href=\"http://localhost:8080/pac/%s\">http://localhost:8080/pac/%s</a></p><br/><br/>\n", md5Hash, md5Hash)
			fmt.Fprintf(w, "<p>Visit the <strong>about</strong> page to learn how to configure your browser.</p>\n")
		} else {
			w.Header().Set("pf-success", "false")
			fmt.Fprintf(w, "<!DOCTYPE html><html lang=\"en\"><head><meta charset=\"utf-8\"><title>PacketFlagon | PAC Creation Failed</title></head>")
			fmt.Fprintf(w, "<body><h1>PAC File Creation Failed</h1></body>")
		}

	} else if r.Method == "GET" {

		tmpl, err := template.New("create").ParseFiles("assets/templates/create.html")
		err = tmpl.Execute(w, templateData)

		if err != nil {
			log.Fatal(err)
		}

	} else {
		//We don't accept PUT / DELETE etc on the /create/ URL
		w.Header().Set("pf-success", "false")
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Method %s is not allowed!", r.Method)
		fmt.Fprintf(w, "Method %s is not allowed!", r.Method)
	}
}
