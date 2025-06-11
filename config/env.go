package config

import (
	"fmt"
	"os"
)

func DatabaseHost() string {
	return os.Getenv("DB_HOST")
}

func DatabasePort() string {
	return os.Getenv("DB_PORT")
}

func DatabaseUser() string {
	return os.Getenv("DB_USER")
}

func DatabasePassword() string {
	return os.Getenv("DB_PASSWORD")
}

func DatabaseName() string {
	return os.Getenv("DB_NAME")
}

func DatabaseURL() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", DatabaseHost(), DatabaseUser(), DatabasePassword(), DatabaseName(), DatabasePort())
}

func APIPort() string {
	return os.Getenv("API_PORT")
}
