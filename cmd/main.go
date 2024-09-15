package main

import (
	"bookservice/internal/routes"
	"log"
	"os"

	"github.com/joho/godotenv"
)


func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	app := routes.SetupRouter()

	log.Fatal(app.Run(os.Getenv("PORT")))
}
