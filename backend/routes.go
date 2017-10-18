package main

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"log"
	"net/http"
	"strings"
)

var (
	db = initDB()
)

type RegistrationResult struct {
	Success bool `json:"success,omiempty"`
}

type Pengguna struct {
	ID       uint16 `json:"id,omiempty"`
	Username string `json:"username,omiempty"`
	FullName string `json:"nama,omiempty"`
	Email    string `json:"email,omiempty"`
}

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "hello, kudo. namaku dendi. hobiku ngoding \n")
}

func HandleRegistration(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	nama := r.PostFormValue("nama")
	email := r.PostFormValue("email")
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")

	// md5 hash password
	hasher := md5.New()
	md5pass := hex.EncodeToString(hasher.Sum([]byte(password)))

	// cek apakah username sudah terdaftar
	reg := RegistrationResult{}

	if !isUsernameAvail(username) {
		reg.Success = false
	}else{
		// query tambah data ke dalam tabel pengguna
		stmt, err := db.Prepare("insert into pengguna (username, password, nama, email) values(?,?,?,?);")

		if err != nil {
			log.Println(err.Error())
		}

		_, err = stmt.Exec(username, md5pass, nama, email)

		if err != nil {
			log.Println(err.Error())
			reg.Success = false
		} else {
			log.Println("Pengguna baru berhasil teregistrasi")
			reg.Success = true
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reg)
}

func HandleEditProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	r.ParseForm()
	id := r.PostFormValue("id")
	nama := r.PostFormValue("nama")
	email := r.PostFormValue("email")
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")

	// md5 hash password
	hasher := md5.New()
	md5pass := hex.EncodeToString(hasher.Sum([]byte(password)))

	// query edit data pengguna
	stmt, err := db.Prepare("update pengguna set username=?, password=?, nama=?, email=? where id=?;")

	if err != nil {
		log.Println(err.Error())
	}

	_, err = stmt.Exec(username, md5pass, nama, email, id)

	reg := RegistrationResult{}

	if err != nil {
		log.Println(err.Error())
		reg.Success = false
	} else {
		log.Println("Data pengguna berhasil diedit")
		reg.Success = true
	}

	log.Println()

	json.NewEncoder(w).Encode(reg)
}

func HandleSelfProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	auth := strings.SplitN(r.Header.Get("Authorization"), " ", 2)

	if len(auth) != 2 || auth[0] != "Basic" {
		http.Error(w, "authorization failed", http.StatusUnauthorized)
		return
	}

	payload, _ := base64.StdEncoding.DecodeString(auth[1])
	pair := strings.SplitN(string(payload), ":", 2)

	_, pengguna := validate(pair[0], pair[1])

	json.NewEncoder(w).Encode(pengguna)
}
