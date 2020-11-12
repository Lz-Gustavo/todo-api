package main

import (
	"encoding/json"
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
	args := r.URL.Query()
	if _, ok := args["usr"]; !ok {
		http.Error(w, "no user informed", http.StatusBadRequest)
		return
	}
	usr := args["usr"][0]

	switch r.Method {
	case "POST":
		// TODO: parse description from POST body, JSON maybe?
		err := createTask(r.Context(), usr, "testandoae")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	case "GET":
		// TODO: get all or with a certain status if requested, return JSON
		tasks, err := listTasks(r.Context(), usr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		raw, err := json.Marshal(tasks)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = w.Write(raw)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	case "PUT":
		// TODO: parse id and description
		err := updateTask(r.Context(), usr, 0, "")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	case "DELETE":
		// TODO: parse id
		err := deleteTask(r.Context(), usr, 0)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

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
