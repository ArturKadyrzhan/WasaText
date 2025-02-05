package main

import (
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

type WebApiConfiguration struct {
	ServerPort string
	ApiSecret  string
	TtlHour    int
}

func LoadConfiguration() (WebApiConfiguration, error) {
	err := godotenv.Load()
	if err != nil {
		return WebApiConfiguration{}, err
	}

	ttlHour, _ := strconv.Atoi(os.Getenv("TTL_HOUR"))
	cfg := WebApiConfiguration{
		ServerPort: os.Getenv("SERVER_PORT"),
		ApiSecret:  os.Getenv("API_SECRET"),
		TtlHour:    ttlHour}
	return cfg, nil
}
