package bootstrap

import "github.com/go-telegram/bot"

func NewTelegram(token string) (*bot.Bot, error) {
	return bot.New(token)
}
