package main

import (
	"log"
	"os"

	"go_final_project/pkg/db"
	"go_final_project/pkg/server"
)

func main() {
	dbFile := "scheduler.db"
	if env := os.Getenv("TODO_DBFILE"); env != "" {
		dbFile = env
	}

	if err := db.Init(dbFile); err != nil {
		log.Fatal(err)
	}

	if err := server.Start("./web"); err != nil {
		log.Fatal(err)
	}
}
