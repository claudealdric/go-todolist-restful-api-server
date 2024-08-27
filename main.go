package main

import (
	"fmt"
	"log"
	"net/http"
)

const port = 8080

func main() {
	server := NewServer()
	log.Printf("Starting server on port %d", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), server)
	log.Fatal(err)
}
