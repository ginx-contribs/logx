package logx

import (
	"fmt"
	"github.com/ginx-contribs/tint"
	"io"
	"log/slog"
	"os"
	"time"
)

const (
	TextFormat = "TEXT"
	JSONFormat = "JSON"

	PromptKey = tint.PromptKey
)

type HandlerOptions struct {
	// Log level, default INFO
	Level slog.Level `mapstructure:"level"`

	// log prompt
	Prompt string `mapstructure:"prompt"`

	// TEXT or JSON
	Format string `mapstructure:"format"`

	// whether to show source files
	Source bool `mapstructure:"source"`

	// custom time format
	TimeFormat string `mapstructure:"timeFormat"`

	ReplaceAttr func(groups []string, a slog.Attr) slog.Attr

	// color log only available in TEXT format
	Color bool `mapstructure:"color"`
}

// NewHandler returns a default log handler, support TEXT and JSON format.
func NewHandler(writer io.Writer, options *HandlerOptions) (slog.Handler, error) {

	if writer == nil {
		writer = os.Stdout
	}

	if options == nil {
		options = &HandlerOptions{}
	}

	if options.Format == "" {
		options.Format = TextFormat
	}

	if options.TimeFormat == "" {
		options.TimeFormat = time.DateTime
	}

	// wrap options.ReplaceAttr
	replaceAttr := func(groups []string, a slog.Attr) slog.Attr {
		switch a.Key {
		case slog.TimeKey:
			a.Value = slog.AnyValue(a.Value.Time().Format(options.TimeFormat))
		}
		if options.ReplaceAttr != nil {
			return options.ReplaceAttr(groups, a)
		}
		return a
	}

	switch options.Format {
	case TextFormat:
		return tint.NewHandler(writer, &tint.Options{
			AddSource:   options.Source,
			Level:       options.Level,
			Prompt:      options.Prompt,
			ReplaceAttr: replaceAttr,
			TimeFormat:  options.TimeFormat,
			NoColor:     !options.Color,
		}), nil
	case JSONFormat:
		return slog.NewJSONHandler(writer, &slog.HandlerOptions{
			AddSource:   options.Source,
			Level:       options.Level,
			ReplaceAttr: replaceAttr,
		}).WithAttrs([]slog.Attr{slog.String(PromptKey, options.Prompt)}), nil
	default:
		return nil, fmt.Errorf("unsupported format %s", options.Format)
	}
}
