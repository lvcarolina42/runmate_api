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

func Production() bool {
	return Env() == Prod
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

func KafkaHost() string {
	return os.Getenv("KAFKA_HOST")
}

func KafkaPort() string {
	return os.Getenv("KAFKA_PORT")
}

func KafkaURL() string {
	return fmt.Sprintf("%s:%s", KafkaHost(), KafkaPort())
}

func KafkaUsername() string {
	return "$ConnectionString"
}

func KafkaAccessKeyName() string {
	return os.Getenv("KAFKA_ACCESS_KEY_NAME")
}

func KafkaAccessKey() string {
	return os.Getenv("KAFKA_ACCESS_KEY")
}

func KafkaPassword() string {
	return fmt.Sprintf("Endpoint=sb://%s/;SharedAccessKeyName=%s;SharedAccessKey=%s", KafkaHost(), KafkaAccessKeyName(), KafkaAccessKey())
}

func FirebaseCredentials() []byte {
	return []byte(os.Getenv("FIREBASE_CREDENTIALS"))
}
