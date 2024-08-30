package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/claudealdric/go-todolist-restful-api-server/api"
	"github.com/claudealdric/go-todolist-restful-api-server/data"
)

const port = 8080
const dbDirName = "tmp"
const dbFileName = "db.json"

func main() {
	if err := os.MkdirAll(dbDirName, os.ModePerm); err != nil {
		log.Fatalf("failed to create directory: %v", err)
	}
	dbPath := filepath.Join(dbDirName, dbFileName)
	dbFile, err := os.OpenFile(dbPath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("problem opening %s %v", dbFileName, err)
	}
	store, err := data.NewFileSystemStore(dbFile)
	if err != nil {
		log.Fatalf("problem creating file system data store: %v", err)
	}
	server := api.NewServer(store)
	log.Printf("Starting server on port %d", port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), server)
	log.Fatal(err)
}
