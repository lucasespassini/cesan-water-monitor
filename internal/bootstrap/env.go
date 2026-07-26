package bootstrap

import (
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	AGazetaCesanUrl string
	CronInterval    string
	RedisAddress    string
	RedisUser       string
	RedisPassword   string
	TelegramToken   string
	Environment     string
}

func NewEnv() *Env {
	_ = godotenv.Load()

	return &Env{
		AGazetaCesanUrl: os.Getenv("AGAZETA_CESAN_URL"),
		CronInterval:    os.Getenv("CRON_INTERVAL"),
		RedisAddress:    os.Getenv("REDIS_ADDRESS"),
		RedisUser:       os.Getenv("REDIS_USER"),
		RedisPassword:   os.Getenv("REDIS_PASSWORD"),
		TelegramToken:   os.Getenv("TELEGRAM_TOKEN"),
		Environment:     os.Getenv("ENVIRONMENT"),
	}
}
