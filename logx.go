package logx

import (
	"errors"
	"log/slog"
)

type Options struct {
	// multi handler
	handlers []slog.Handler
}

type Option func(*Options)

func WithHandlers(handlers ...slog.Handler) Option {
	return func(options *Options) {
		options.handlers = handlers
	}
}

// Logger is a slog handler wrapper
type Logger struct {
	l *slog.Logger

	options Options

	handler slog.Handler

	onClose []func() error
}

func (l *Logger) OnClose(closeFns ...func() error) {
	l.onClose = append(l.onClose, closeFns...)
}

func (l *Logger) Handler() slog.Handler {
	return l.handler
}

func (l *Logger) Slog() *slog.Logger {
	return l.l
}

func (l *Logger) Close() error {
	var joinErr error
	for _, closeFn := range l.onClose {
		if err := closeFn(); err != nil {
			joinErr = errors.Join(joinErr, err)
		}
	}
	return joinErr
}

// New return a logger with options
func New(opts ...Option) (*Logger, error) {

	logger := new(Logger)

	var options Options
	for _, opt := range opts {
		opt(&options)
	}

	if len(options.handlers) == 0 {
		// default handlers
		writer, err := NewWriter(nil)
		if err != nil {
			return nil, err
		}

		logger.OnClose(func() error {
			return writer.Close()
		})

		handler, err := NewHandler(writer, nil)
		if err != nil {
			return nil, err
		}

		options.handlers = append(options.handlers, handler)
	}

	logger.options = options

	logger.handler = Multi(options.handlers...)

	logger.l = slog.New(logger.handler)

	return logger, nil
}

// FileLoggerOption combines with WriterOptions and HandlerOptions
type FileLoggerOption struct {
	*WriterOptions
	*HandlerOptions
}

// NewFileLogger is a helper function, returns a single handler logger.
func NewFileLogger(options *FileLoggerOption) (*Logger, error) {

	if options == nil {
		options = new(FileLoggerOption)
	}

	writer, err := NewWriter(options.WriterOptions)
	if err != nil {
		return nil, err
	}

	handler, err := NewHandler(writer, options.HandlerOptions)
	if err != nil {
		return nil, err
	}

	logger, err := New(WithHandlers(handler))
	if err != nil {
		return nil, err
	}

	logger.OnClose(func() error {
		return writer.Close()
	})

	return logger, nil
}
