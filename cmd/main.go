package main

import (
	"cesan-scraping/internal/bootstrap"
	"cesan-scraping/internal/domain"
	"cesan-scraping/internal/usecase"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-co-op/gocron-ui/server"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	app, err := bootstrap.App(ctx)
	if err != nil {
		log.Fatal(err)
		return
	}

	var jobs []domain.Job

	jobs = append(
		jobs,
		usecase.MonitorWaterOutageUseCase{
			Env:      app.Env,
			Redis:    app.Redis,
			Telegram: app.Telegram,
		},
	)

	app.RegisterJobs(jobs)

	srv := server.NewServer(app.JobManager.Scheduler(), 8081)
	log.Println("GoCron UI available at http://localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", srv.Router))
}
