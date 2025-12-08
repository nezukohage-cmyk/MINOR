package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	//"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}
}

func GetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		fmt.Printf("Environment variable %s is not set\n", key)
	}
	return value
}
