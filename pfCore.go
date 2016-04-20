package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/boltdb/bolt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func CreatePAC(r *http.Request, db *bolt.DB) (pac PAC, err error) {
	md5Hash := GetMD5Hash(strconv.FormatInt(time.Now().Unix(), 10))

	r.ParseForm()
	urls := r.FormValue("urls")
	pac = PAC{ID: md5Hash, Name: r.FormValue("name"), Description: r.FormValue("desc"), Password: GetMD5Hash(r.FormValue("password"))}

	pac.URLs = strings.Split(urls, ",")

	for index, url := range pac.URLs {
		pac.URLs[index] = strings.TrimSpace(url)
	}

	buf, err := json.Marshal(pac)

	log.Print(string(buf))

	if err == nil {

		err = db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("pacs"))
			err := b.Put([]byte(md5Hash), buf)
			return err
		})
		return pac, err
	} else {
		return pac, err
	}
}

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
