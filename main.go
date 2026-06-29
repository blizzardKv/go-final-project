package main

import (
	"log"

	"go_final_project/pkg/server"
)

func main() {
	if err := server.Start("./web"); err != nil {
		log.Fatal(err)
	}
}
