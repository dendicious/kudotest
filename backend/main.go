package main

import (
	"database/sql"
	"log"
	"net/http"
)

func initDB() *sql.DB {

	// database initialization
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/kudotest")
	if err != nil {
		log.Panic(err.Error())
	}

	err = db.Ping()

	if err != nil {
		log.Panic(err.Error())
	}

	return db
}

func main() {

	// public endpoints
	http.HandleFunc("/", HandleIndex)
	http.HandleFunc("/registration", PostHandler(HandleRegistration))

	// protected endpoints
	http.HandleFunc("/profile", GetHandler(basicAuth(HandleSelfProfile)))
	http.HandleFunc("/edit-profile", PutHandler(basicAuth(HandleEditProfile)))
	
	log.Fatal(http.ListenAndServe(":8080", nil))
}
