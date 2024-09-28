package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"

	"github.com/claudealdric/go-todolist-restful-api-server/api"
	"github.com/claudealdric/go-todolist-restful-api-server/data"
)

const port = 8080

func main() {
	os.Remove("./data/data.db")
	db, err := sql.Open("sqlite3", "./data/data.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	data.InitDb(db)

	store := data.NewSqliteStore(db)
	if err != nil {
		log.Fatalf("problem creating file system data store: %v", err)
	}
	server := api.NewServer(store)
	log.Printf("Starting server on port %d", port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), server)
	log.Fatal(err)
}
