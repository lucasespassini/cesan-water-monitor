package bootstrap

import (
	"cesan-scraping/internal/domain"
	"context"

	"github.com/go-telegram/bot"
	"github.com/redis/go-redis/v9"
)

type Application struct {
	Env        *Env
	JobManager *domain.JobManager
	Redis      *redis.Client
	Telegram   *bot.Bot
	Context    context.Context
}

func App(ctx context.Context) (*Application, error) {
	jm, err := domain.NewJobManager(ctx)
	if err != nil {
		return nil, err
	}

	env := NewEnv()
	redis := NewRedisClient(env.RedisAddress)
	err = redis.Ping(ctx).Err()
	if err != nil {
		panic(err)
	}

	telegram, err := NewTelegram(env.TelegramToken)
	if err != nil {
		return nil, err
	}

	app := &Application{
		Env:        env,
		JobManager: jm,
		Redis:      redis,
		Telegram:   telegram,
		Context:    ctx,
	}

	return app, nil
}

func (a *Application) RegisterJobs(jobs []domain.Job) {
	for _, job := range jobs {
		a.JobManager.RegisterJob(job)
	}

	a.JobManager.StartScheduler()
}
