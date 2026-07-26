package usecase

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func NotifyWaterOutage(
	ctx context.Context,
	b *bot.Bot,
	chatID int64,
	title string,
	url string,
) error {
	text := fmt.Sprintf(
		`🚰 <b>Falta de água detectada</b>

📰 <b>%s</b>

🔗 %s`, title, url)

	_, err := b.SendMessage(
		ctx,
		&bot.SendMessageParams{
			ChatID:    chatID,
			Text:      text,
			ParseMode: models.ParseModeHTML,
		},
	)

	return err
}
