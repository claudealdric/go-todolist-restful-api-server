package main

import (
	"fmt"
	"log"
	"net/http"
)

const port = 8080

func main() {
	router := NewRouter()
	log.Printf("Starting server on port %d\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), router)
	log.Fatal(err)
}
