package bootstrap

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	AGazetaCesanUrl string
	CronInterval    string
	RedisAddress    string
	TelegramToken   string
	Environment     string
}

func NewEnv() *Env {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Env{
		AGazetaCesanUrl: os.Getenv("AGAZETA_CESAN_URL"),
		CronInterval:    os.Getenv("CRON_INTERVAL"),
		RedisAddress:    os.Getenv("REDIS_ADDRESS"),
		TelegramToken:   os.Getenv("TELEGRAM_TOKEN"),
		Environment:     os.Getenv("ENVIRONMENT"),
	}
}
