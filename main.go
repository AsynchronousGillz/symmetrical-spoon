package main

import (
	"log"
	"net/http"
)

var (
	// Version software version
	Version string
	// Build software build
	Build string
)

func main() {

	router := NewRouter()
	log.Fatal(http.ListenAndServe(":8080", router))
}
