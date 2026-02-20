package main

import (
	"mecanica_xpto/internal/infrastructure/database"
)

func main() {
	db := database.ConnectDatabase()
	database.Migrate(db)
	database.Seed(db)
}
