package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"

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
		id, desc, status, err := parseUpdateArgs(args)
		if err != nil {
			http.Error(w, "error parsing args on update operation, err:"+err.Error(), http.StatusBadRequest)
			return
		}

		err = updateTask(r.Context(), usr, id, desc, status)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	case "DELETE":
		id, err := parseIDFromArgs(args)
		if err != nil {
			http.Error(w, "error parsing args on delete operation, err:"+err.Error(), http.StatusBadRequest)
			return
		}

		err = deleteTask(r.Context(), usr, id)
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

	args := r.URL.Query()
	if _, ok := args["usr"]; !ok {
		http.Error(w, "no user informed", http.StatusBadRequest)
		return
	}
	usr := args["usr"][0]

	id, err := parseIDFromArgs(args)
	if err != nil {
		http.Error(w, "error parsing args on complete task operation, err:"+err.Error(), http.StatusBadRequest)
		return
	}

	err = completeTask(r.Context(), usr, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func tasksRestoreHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	args := r.URL.Query()
	if _, ok := args["usr"]; !ok {
		http.Error(w, "no user informed", http.StatusBadRequest)
		return
	}
	usr := args["usr"][0]

	id, err := parseIDFromArgs(args)
	if err != nil {
		http.Error(w, "error parsing args on restore task operation, err:"+err.Error(), http.StatusBadRequest)
		return
	}

	err = restoreTask(r.Context(), usr, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func parseIDFromArgs(args url.Values) (int, error) {
	if _, ok := args["id"]; !ok {
		return 0, errors.New("no id informed for update operation")
	}

	id, err := strconv.Atoi(args["id"][0])
	if err != nil {
		return 0, errors.New("could not parse id value informed")
	}
	return id, nil
}

func parseUpdateArgs(args url.Values) (int, string, string, error) {
	id, err := parseIDFromArgs(args)
	if err != nil {
		return 0, "", "", err
	}

	if _, ok := args["desc"]; !ok {
		return 0, "", "", errors.New("no desc informed for update operation")
	}

	if _, ok := args["status"]; !ok {
		return 0, "", "", errors.New("no status informed for update operation")
	}
	return id, args["desc"][0], args["status"][0], nil
}
