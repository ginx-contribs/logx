package logx

import (
	"log/slog"
	"testing"
	"time"
)

func TestHandlerDefault(t *testing.T) {
	handler, err := NewHandler(nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	logger := slog.New(handler)
	logger.Info("default handler")
}

func TestHandlerOptions(t *testing.T) {
	handler, err := NewHandler(nil, &HandlerOptions{
		Level:       slog.LevelInfo,
		Format:      TextFormat,
		Source:      true,
		TimeFormat:  time.TimeOnly,
		ReplaceAttr: nil,
		Color:       true,
	})
	if err != nil {
		t.Fatal(err)
	}
	logger := slog.New(handler)
	logger.Info("default handler")
}
