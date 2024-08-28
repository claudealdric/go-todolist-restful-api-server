package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/claudealdric/go-todolist-restful-api-server/datastore"
	"github.com/claudealdric/go-todolist-restful-api-server/handlers"
)

const port = 8080
const dbFileName = "db.json"

func main() {
	dbFile, err := os.OpenFile(dbFileName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("problem opening %s %v", dbFileName, err)
	}
	store, err := datastore.NewFileSystemDataStore(dbFile)
	server := handlers.NewServer(store)
	log.Printf("Starting server on port %d", port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), server)
	log.Fatal(err)
}
