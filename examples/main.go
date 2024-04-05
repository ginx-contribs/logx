package main

import (
	"github.com/ginx-contribs/logx"
	"log"
	"log/slog"
)

func main() {
	logger, err := logx.New()
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	slog.SetDefault(logger.Slog())
	// 2024-04-05 12:26:13 INFO hell world!
	slog.Info("hell world!")
}
