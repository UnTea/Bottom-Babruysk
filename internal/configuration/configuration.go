package configuration

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Configuration struct {
	DatabaseConnectionURL string
	HTTPAddress           string
}

func Load() (configuration *Configuration, err error) {
	if err = godotenv.Load(); err != nil {
		log.Fatalf("error loading .env file")
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbSSLMode := os.Getenv("DB_SSLMODE")

	databaseConnectionURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		dbUser, dbPassword, dbHost, dbPort, dbName, dbSSLMode,
	)

	return &Configuration{
		DatabaseConnectionURL: databaseConnectionURL,
		HTTPAddress:           os.Getenv("HTTP_ADDRESS"),
	}, nil
}
