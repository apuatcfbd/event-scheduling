package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"slices"
	"strings"
)

func init() {
	// load .env file
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file.", err.Error())
	}
}

// EnvDBDriver returns env key "DB_DRIVER"'s value after validating
func EnvDBDriver() string {
	d := os.Getenv("DB_DRIVER")
	if d == "" {
		log.Fatalln("DB_DRIVER environment variable not set")
	}

	if !slices.Contains(DbDrivers, d) {
		log.Fatalln("DB_DRIVER environment variable value is invalid, it can be one of:", strings.Join(DbDrivers, ", "))
	}

	return d
}

// EnvDBDsn returns env key "DB_DSN"'s value after validating
func EnvDBDsn() string {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatalln("DB_DSN environment variable not set")
	}

	return dsn
}
