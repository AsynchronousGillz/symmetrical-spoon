package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"html/template"
	"github.com/gorilla/mux"
)

type WebSite struct {
	string
	
}

// Index the root
func Index(w http.ResponseWriter, r *http.Request) {
	if err := templates.ExecuteTemplate(w, "template.html", welcome); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Index the root
func Version(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	a := generateApplicationJSON()
	if err := json.NewEncoder(w).Encode(a); err != nil {
		panic(err)
	}
}

// TransactionIndex index
func TransactionIndex(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	templates := template.Must(template.ParseFiles("templates/welcome-template.html"))
	if err := json.NewEncoder(w).Encode(transactions); err != nil {
		panic(err)
	}
}

// TransactionShow show a transaction
func TransactionShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var transactionID int
	var err error
	if transactionID, err = strconv.Atoi(vars["transactionID"]); err != nil {
		panic(err)
	}
	transaction := RepoFindTransaction(transactionID)
	if transaction.ID > 0 {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(transaction); err != nil {
			panic(err)
		}
		return
	}

	// If we didn't find it, 404
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotFound)
	if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Not Found"}); err != nil {
		panic(err)
	}

}

// TransactionCreate create a transaction
// curl -H "Content-Type: application/json" -d '{"name":"New Transaction"}' http://localhost:8080/transaction
func TransactionCreate(w http.ResponseWriter, r *http.Request) {
	var transaction Transaction
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &transaction); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	t := RepoCreateTransaction(transaction)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(t); err != nil {
		panic(err)
	}
}
