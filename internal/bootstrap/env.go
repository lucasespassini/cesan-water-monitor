package bootstrap

import (
	"log"

	"github.com/spf13/viper"
)

type Env struct {
	AGazetaCesanUrl string `mapstructure:"AGAZETA_CESAN_URL"`
	CronInterval    string `mapstructure:"CRON_INTERVAL"`
	RedisAddress    string `mapstructure:"REDIS_ADDRESS"`
	TelegramToken   string `mapstructure:"TELEGRAM_TOKEN"`
	Environment     string `mapstructure:"ENVIRONMENT"`
}

func NewEnv() *Env {
	env := Env{}
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Can't find the file .env : ", err)
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatal("Environment can't be loaded: ", err)
	}

	return &env
}
