package main

import "net/http"

// Route data struct
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes an array of Route type
type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	Route{
		"TransactionIndex",
		"GET",
		"/transaction",
		TransactionIndex,
	},
	Route{
		"TransactionCreate",
		"POST",
		"/transaction",
		TransactionCreate,
	},
	Route{
		"TransactionShow",
		"GET",
		"/transaction/{transactionId}",
		TransactionShow,
	},
}
