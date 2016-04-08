package main

import (
	//"html/template"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/boltdb/bolt"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"strconv"
)

type CoreConf struct {
	DbPath     string `json:"dbpath"`
	ListenPort int64  `json:"listenport"`
	Debug      bool   `json:"debug"`
	TLS        bool   `json:"tls_enabled"`
	TLSKey     string `json:"tls_private_key"`
	TLSCert    string `json:"tls_certificate"`
}

type TemplateConf struct {
	FQDN       string
	ListenPort int64
	Success    bool
	URLs       []string
	PACName    string
	PACDesc    string
}

func main() {
	log.Println("---------------------------------------------")

	//Grab all our command line config
	configuration := flag.String("conf", "", "path to configuration file")
	flag.Parse()
	conf := readConfig(*configuration)

	//Populate our template conf - normally passed by value to other functions
	tmplConf := TemplateConf{
		FQDN:       "localhost",
		ListenPort: conf.ListenPort,
	}

	//We have our main web thread and the watcher threads
	runtime.GOMAXPROCS(3)

	//Tell people what we're up to
	log.Println("Using config: " + *configuration)
	log.Println("DB Path is: " + conf.DbPath)

	//We need a DB for holding interactions
	db, err := bolt.Open(conf.DbPath, 0600, nil)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("DB opened, all is OK")
	}
	defer db.Close()

	//Make sure our buckets exist
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("pacs"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	//Set webpage elements
	//http.HandleFunc("/", PFWebIndex)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { PFWebIndex(w, r, tmplConf) })
	//http.HandleFunc("/about/", PFWebAbout)
	http.HandleFunc("/about/", func(w http.ResponseWriter, r *http.Request) { PFWebAbout(w, r, tmplConf) })
	http.HandleFunc("/create/", func(w http.ResponseWriter, r *http.Request) { PFWebCreate(w, r, db, tmplConf) })
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("assets/css/"))))
	http.Handle("/font/", http.StripPrefix("/font/", http.FileServer(http.Dir("assets/font/"))))

	//Actually serve the PAC file
	http.HandleFunc("/pac/", func(w http.ResponseWriter, r *http.Request) { PFWebServ(w, r, db) })

	ListenPort := strconv.FormatInt(conf.ListenPort, 10)
	if conf.TLS {
		http.ListenAndServeTLS(":"+ListenPort, conf.TLSCert, conf.TLSKey, nil)
	} else {
		http.ListenAndServe(":"+ListenPort, nil)
	}
	//And we're done
}

// Reads our JSON formatted config file
// and returns a struct
func readConfig(filename string) CoreConf {
	var conf CoreConf

	if filename == "" {
		conf.DbPath = "./urls.bolt"
		conf.ListenPort = 8080
	} else {
		b, err := ioutil.ReadFile(filename)
		if err != nil {
			/*fmt.Println("Cannot read configuration file ", filename)
			  os.Exit(1)*/
			log.Fatal("Cannot read configuration file ", filename)
		}
		var conf CoreConf
		err = json.Unmarshal(b, &conf)
		if err != nil {
			/*fmt.Println("Cannot parse configuration file: ", err)
			  os.Exit(1)*/
			log.Fatal("Cannot read configuration file ", filename)
		}
	}
	return conf
}
