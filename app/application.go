package application

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/quanht2k/golang_basic_training/app/server"
	"github.com/quanht2k/golang_basic_training/config"
)

func Start() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config := config.LoadConfig()
	server := server.NewServer(config)
	error := server.Start()
	if error != nil {
		log.Fatal("Error starting server: ", error)
	}
}