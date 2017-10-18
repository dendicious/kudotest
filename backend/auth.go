package main 

import (
	"encoding/hex"
	"crypto/md5"
	"log"
	"encoding/base64"
	"net/http"
	"strings"
)

type handler func(w http.ResponseWriter, r *http.Request)

func basicAuth(pass handler) handler {

	return func(w http.ResponseWriter, r *http.Request) {

		auth := strings.SplitN(r.Header.Get("Authorization"), " ", 2)

		if len(auth) != 2 || auth[0] != "Basic" {
			http.Error(w, "authorization failed", http.StatusUnauthorized)
			return
		}

		payload, _ := base64.StdEncoding.DecodeString(auth[1])
		pair := strings.SplitN(string(payload), ":", 2)

		valid,_ :=validate(pair[0], pair[1])

		if len(pair) != 2 || !valid {
			http.Error(w, "authorization failed", http.StatusUnauthorized)
			return
		}

		pass(w, r)
	}
}

func validate(username, password string) (bool, Pengguna) {

	db := initDB()
	var pengguna Pengguna
	
	hasher := md5.New()
	passHash := hex.EncodeToString(hasher.Sum([]byte(password)))


	row := db.QueryRow("select id, username, nama, email from pengguna where username = ? and password = ?;",username, passHash)
	err := row.Scan(&pengguna.ID, &pengguna.Username, &pengguna.FullName, &pengguna.Email)

	if err != nil {
		log.Println("Pengguna tidak ditemukan")
		return false, Pengguna{}
	}

	return true, pengguna
}

func isUsernameAvail(username string) (bool) {
	
		db := initDB()
	
		row := db.QueryRow("select id from pengguna where username = ?;",username)
		err := row.Scan()
	
		if err != nil {
			log.Println("Pengguna bisa diregister")
			return true
		}
	
		return false
	}