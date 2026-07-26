package usecase

import (
	"cesan-scraping/internal/bootstrap"
	"cesan-scraping/internal/domain"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-telegram/bot"
	"github.com/gocolly/colly"
	"github.com/redis/go-redis/v9"
)

type MonitorWaterOutageUseCase struct {
	Env      *bootstrap.Env
	Redis    *redis.Client
	Telegram *bot.Bot
	Notifier *domain.Notifier
}

var waterOutageRegex = regexp.MustCompile(
	`(?i)(sem água|falta d['’]?água|falta de água|desabastecimento|abastecimento.*(interrompido|suspenso)|interrupção.*abastecimento)`,
)

func (u MonitorWaterOutageUseCase) Name() string {
	return "NotifyWaterOutageJob"
}

func (u MonitorWaterOutageUseCase) Duration() time.Duration {
	cronInterval, err := strconv.ParseInt(u.Env.CronInterval, 10, 64)
	if err != nil {
		return 4 * time.Hour
	}

	return time.Duration(cronInterval) * time.Hour
}

func (u MonitorWaterOutageUseCase) Run(ctx context.Context) error {
	c := colly.NewCollector()

	c.OnHTML("a", func(e *colly.HTMLElement) {
		title := strings.TrimSpace(e.ChildText("h2"))
		href := e.Attr("href")

		if title == "" {
			return
		}

		if !strings.Contains(href, "/es/") &&
			!strings.Contains(href, "/tema/") {
			return
		}

		normalizedTitle := strings.Join(strings.Fields(title), " ")
		normalizedTitle = strings.ToLower(normalizedTitle)

		if !waterOutageRegex.MatchString(normalizedTitle) {
			return
		}

		key := newsKey(href)

		ok, err := u.Redis.SetNX(ctx, key, 1, 30*24*time.Hour).Result()
		if err != nil {
			log.Fatal(err)
			return
		}

		if !ok {
			return
		}

		err = NotifyWaterOutage(
			ctx,
			u.Telegram,
			-1004491522361,
			normalizedTitle,
			e.Request.AbsoluteURL(href),
		)
		if err != nil {
			log.Printf("failed to send Telegram notification: %v", err)
			return
		}
		fmt.Println("Message send")
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	targetURL, cleanup, err := resolveURL(u.Env.AGazetaCesanUrl)
	if err != nil {
		return err
	}
	defer cleanup()

	return c.Visit(targetURL)
}

func newsKey(url string) string {
	sum := sha256.Sum256([]byte(url))
	return "news:agazeta:" + hex.EncodeToString(sum[:])
}

// resolveURL converts a file:// URL into an ephemeral local HTTP server.
// For http/https URLs, it returns the original URL unchanged.
// The cleanup function must be called via defer to release resources.
func resolveURL(rawURL string) (targetURL string, cleanup func(), err error) {
	noop := func() {}

	if !strings.HasPrefix(rawURL, "file://") {
		return rawURL, noop, nil
	}

	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", noop, fmt.Errorf("invalid file URL: %w", err)
	}

	filePath := parsed.Path
	if len(filePath) > 2 && filePath[0] == '/' && filePath[2] == ':' {
		filePath = filePath[1:]
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", noop, fmt.Errorf("failed to read local file: %w", err)
	}

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return "", noop, fmt.Errorf("failed to start local server: %w", err)
	}

	body := string(data)
	go http.Serve(listener, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(body))
	}))

	return "http://" + listener.Addr().String() + "/", func() { listener.Close() }, nil
}
