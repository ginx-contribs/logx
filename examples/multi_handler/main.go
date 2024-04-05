package main

import (
	"github.com/ginx-contribs/logx"
	slogsentry "github.com/samber/slog-sentry/v2"
	"log"
	"log/slog"
	"os"
	"time"
)

func main() {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:           "https://xxxxxxx@yyyyyyy.ingest.sentry.io/zzzzzzz",
		EnableTracing: false,
	})
	if err != nil {
		log.Fatal(err)
	}

	logger, err := logx.New(
		logx.WithHandlers(
			slog.NewJSONHandler(os.Stderr, nil),
			slogsentry.Option{Level: slog.LevelDebug}.NewSentryHandler(),
		),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	logger.OnClose(func() error {
		return sentry.Flush(2 * time.Second)
	})

	slog.SetDefault(logger.Slog())

	slog.Info("hello world")
}
