package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"opsws/master/internal/api"
)

func main() {
	db, err := sql.Open("sqlite3", "opsws.db")
	if err != nil {
		log.Fatalf("failed to open sqlite db: %v", err)
	}
	defer db.Close()

	if err := api.Migrate(db); err != nil {
		log.Fatalf("db migration error: %v", err)
	}

	r := mux.NewRouter()
	api.RegisterPipelineRoutes(r, db)
	log.Println("OpsWS Master started at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
