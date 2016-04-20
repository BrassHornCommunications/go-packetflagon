package main

import (
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"net/http"
	"strings"
)

func PFWebServ(w http.ResponseWriter, r *http.Request, db *bolt.DB) {
	w.Header().Set("Content-Type", "application/json")
	//fmt.Fprintf(w, "PAC file\n")
	url := strings.Split(r.URL.String(), "/")
	//fmt.Fprintf(w, url[2])
	//var urls string
	var pacJSON []byte

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("pacs"))
		pacJSON = b.Get([]byte(url[2]))

		return nil
	})

	if err != nil {
		//fmt.Fprintf(w, "function FindProxyForURL(url, host)\n{\nreturn \"DIRECT\";\n}\n")
		http.Error(w, err.Error(), 500)
	} else {
		var pac PAC
		err := json.Unmarshal(pacJSON, &pac)

		if err == nil {
			fmt.Fprintf(w, "//%s\n", pac.Name)
			fmt.Fprintf(w, "//%s\n\n", pac.Description)
			fmt.Fprintf(w, "function FindProxyForURL(url, host)\n{\n")
			fmt.Fprintf(w, " var list = new Array(")

			/*urlSlice := strings.Split(urls, ",")

			for _, url := range urlSlice {
				fmt.Fprintf(w, "\"%s\", ", url)
			}*/
			for _, url := range pac.URLs {
				fmt.Fprintf(w, "\"%s\", ", url)
			}
			fmt.Fprintf(w, "\"torproject.org\");\n")

			fmt.Fprintf(w, "  for(var i=0; i < list.length; i++)\n  {\n  if (shExpMatch(host, list[i]))\n  {\n   return \"SOCKS 127.0.0.1:9050; 127.0.0.1:9051\";\n  }\n}\n return \"DIRECT\";\n}")
		} else {
			http.Error(w, err.Error(), 500)
		}
	}
}
