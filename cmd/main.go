package main

import (
	"DevelopsToday/config"
	"DevelopsToday/internal/app"
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Hello, DevelopsToday!")

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}
	app.Run(cfg)

}
