package config

import (
	"os"
)

type Config struct {
	Port string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	RabbitMQURI     string
	RabbitQueueName string

	PrometheusMetricsPort string
}

func LoadEnv() *Config {

	return &Config{
		Port: os.Getenv("PORT"),

		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),

		RabbitMQURI:           os.Getenv("RABBITMQ_URI"),
		RabbitQueueName:       os.Getenv("RABBITMQ_QUEUE"),
		PrometheusMetricsPort: os.Getenv("PROMETHEUS_METRICS_PORT"),
	}
}
