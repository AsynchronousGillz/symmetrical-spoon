package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"github.com/gorilla/mux"
)

// https://golang.org/src/net/http/status.go

// Index the root
func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	a := generateApplicationJSON()
	if err := json.NewEncoder(w).Encode(a); err != nil {
		panic(err)
	}
}

// Application status
func ApplicationStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	a := applicationStatusJSON{"OK"}
	if err := json.NewEncoder(w).Encode(a); err != nil {
		panic(err)
	}
}

// TransactionList shows a list of transactions
// curl http://localhost:8080/transactions
func TransactionList(w http.ResponseWriter, r *http.Request) {
	transactions, err := Transactions()
	if err != nil {
		// We didn't find it, 404
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusInternalServerError) // Internal Server Error
		textErr := fmt.Sprintf("error: %s", err)
		if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusInternalServerError, Text: textErr}); err != nil {
			panic(err)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(transactions); err != nil {
		panic(err)
	}
}

// TransactionShow show a transaction
// curl http://localhost:8080/transaction/uuid
func TransactionShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var transactionID string = vars["transactionID"]
	transaction, err := findTransaction(transactionID)
	if err != nil {
		// We didn't find it, 404
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotFound)
		textErr := fmt.Sprintf("id %s not found", transactionID)
		if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: textErr}); err != nil {
			panic(err)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(transaction); err != nil {
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
		w.WriteHeader(http.StatusUnprocessableEntity) // unprocessable entity
		textErr := fmt.Sprintf("error: %s", err)
		if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusUnprocessableEntity, Text: textErr}); err != nil {
			panic(err)
		}
	}
	t, err := CreateTransaction(transaction)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusInternalServerError) // Internal Server Error
		textErr := fmt.Sprintf("error: %s", err)
		if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusInternalServerError, Text: textErr}); err != nil {
			panic(err)
		}
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(t); err != nil {
		panic(err)
	}
}

// TransactionDelete delete a transaction
// curl -X DELETE http://localhost:8080/transaction/uuid
func TransactionDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var transactionID string = vars["transactionID"]
	err := DeleteTransaction(transactionID)
	if err != nil {
		// We didn't find it, 404
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotFound)
		textErr := fmt.Sprintf("id %s not found", transactionID)
		if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: textErr}); err != nil {
			panic(err)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}
