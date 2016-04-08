package main

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/boltdb/bolt"
	"net/http"
	"time"
)

func CreatePAC(r *http.Request, db *bolt.DB) (md5Hash string, err error) {
	md5Hash = GetMD5Hash(string(time.Now().Unix()))

	r.ParseForm()
	urls := r.FormValue("urls")

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("pacs"))
		err := b.Put([]byte(md5Hash), []byte(urls))
		return err
	})
	return md5Hash, err
}

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
