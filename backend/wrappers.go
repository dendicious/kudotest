package main

import (
	"net/http"
)

func GetHandler(h handler) handler {

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			h(w, r)
			return
		}
		http.Error(w, "get only", http.StatusMethodNotAllowed)
	}
}

func PostHandler(h handler) handler {

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			h(w, r)
			return
		}
		http.Error(w, "post only", http.StatusMethodNotAllowed)
	}
}

func PutHandler(h handler) handler {
    
        return func(w http.ResponseWriter, r *http.Request) {
            if r.Method == "PUT" {
                h(w, r)
                return
            }
            http.Error(w, "put only", http.StatusMethodNotAllowed)
        }
    }
