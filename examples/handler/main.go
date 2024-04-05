package main

import (
	"github.com/ginx-contribs/logx"
	"log"
	"log/slog"
	"time"
)

func main() {
	// file writer
	writer, err := logx.NewWriter(&logx.WriterOptions{
		Filename:      "./logx.log",
		DisableStderr: true,
	})

	if err != nil {
		log.Fatal(err)
	}

	handler, err := logx.NewHandler(writer, &logx.HandlerOptions{
		Level:      slog.LevelInfo,
		TimeFormat: time.TimeOnly,
	})

	logger, err := logx.New(
		logx.WithHandlers(handler),
	)

	if err != nil {
		log.Fatal(err)
	}

	// register close hooks
	logger.OnClose(func() error {
		return writer.Close()
	})

	defer logger.Close()

	slog.SetDefault(logger.Slog())

	slog.Info("hello world")
}
