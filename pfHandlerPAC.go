package main

import (
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
	var urls string

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("pacs"))
		urls = string(b.Get([]byte(url[2])))

		//fmt.Printf("The answer is: %s\n", v)
		return nil
	})

	if err != nil || len(urls) == 0 {
		fmt.Fprintf(w, "function FindProxyForURL(url, host)\n{\nreturn \"DIRECT\";\n}\n")
	} else {
		fmt.Fprintf(w, "function FindProxyForURL(url, host)\n{\n")
		fmt.Fprintf(w, " var list = new Array(")
		urlSlice := strings.Split(urls, ",")

		for _, url := range urlSlice {
			fmt.Fprintf(w, "\"%s\", ", url)
		}
		fmt.Fprintf(w, "\"torproject.org\");\n")

		fmt.Fprintf(w, "  for(var i=0; i < list.length; i++)\n  {\n  if (shExpMatch(host, list[i]))\n  {\n   return \"SOCKS 127.0.0.1:9050; 127.0.0.1:9051\";\n  }\n}\n return \"DIRECT\";\n}")
	}
}
