package main

// Transaction the datastructure to hold a transaction
type Transaction struct {
	ID      int     `json:"id"`
	Date    string  `json:"date"`
	Name    string  `json:"name"`
	Info    string  `json:"info"`
	Amount  float64 `json:"amount"`
	Account string  `json:"account"`
}

// Transactions an array of Transaction
type Transactions []Transaction
