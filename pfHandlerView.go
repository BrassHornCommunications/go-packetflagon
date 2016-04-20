package main

import (
	"encoding/json"
	"github.com/boltdb/bolt"
	"html/template"
	"log"
	"net/http"
	"strings"
)

func PFWebView(w http.ResponseWriter, r *http.Request, db *bolt.DB, templateData TemplateConf) {
	url := strings.Split(r.URL.String(), "/")
	var pacJSON []byte

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("pacs"))
		pacJSON = b.Get([]byte(url[2]))

		return nil
	})

	if err != nil {
		http.Error(w, err.Error(), 500)
	} else {
		var pac PAC
		err := json.Unmarshal(pacJSON, &pac)

		if err == nil {
			tmpl, err := template.New("view").ParseFiles("assets/templates/view.html")
			templateData.PAC = pac
			err = tmpl.Execute(w, templateData)

			if err != nil {
				log.Fatal(err)
			}

		} else {
			http.Error(w, err.Error(), 500)
		}
	}
}
