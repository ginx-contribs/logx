package logx

import (
	"testing"
)

func TestDefault(t *testing.T) {
	logger, err := New()
	if err != nil {
		t.Fatal(err)
	}
	logger.Slog().Info("default logger")
}
