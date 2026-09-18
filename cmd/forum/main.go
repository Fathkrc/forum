package main

import (
	"cmd/forum/internal/api"
	"cmd/forum/internal/database"
	"log"
)

func main() {
	db, err := database.NewDataBase("forum")
	if err != nil {
		log.Fatal("Error: ", err)
	}
	defer db.Close()

	if err := api.Run("8080", db); err != nil {
		log.Fatal(err)
	}
}
