package config

import (
	"fmt"
	"os"
)

const (
	Dev  = "dev"
	Prod = "prod"
)

func Env() string {
	return os.Getenv("ENV")
}

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
	url := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", DatabaseHost(), DatabaseUser(), DatabasePassword(), DatabaseName(), DatabasePort())
	if Env() == Dev {
		url += " sslmode=disable"
	}
	return url
}

func APIPort() string {
	return os.Getenv("API_PORT")
}

func KafkaURL() string {
	return os.Getenv("KAFKA_URL")
}
