package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func createRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/tasks/", tasksHandler).
		Methods("POST", "GET", "PUT", "DELETE").
		Schemes("", "http")

	r.HandleFunc("/tasks/complete/", tasksCompleteHandler).
		Methods("GET").
		Schemes("", "http")

	r.HandleFunc("/tasks/restore/", tasksRestoreHandler).
		Methods("GET").
		Schemes("", "http")

	return r
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		// TODO

	case "GET":
		// TODO: get all or with a certain status if requested, return JSON

	case "PUT":
		// TODO

	case "DELETE":
		// TODO

	default:
		http.Error(w, "", http.StatusBadRequest)
		return
	}
}

func tasksCompleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	// TODO
}

func tasksRestoreHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	// TODO
}
