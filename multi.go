package logx

import (
	"context"
	"errors"
	"fmt"
	"io"

	"log/slog"

	"github.com/samber/lo"
)

var _ slog.Handler = (*MultiHandler)(nil)

// MultiHandler implements slog.Handler interface, holds multiple handlers and call them sequentially.
type MultiHandler struct {
	handlers []slog.Handler
}

// Multi distributes records to multiple slog.Handler in parallel
func Multi(handlers ...slog.Handler) slog.Handler {
	return &MultiHandler{
		handlers: handlers,
	}
}

func (h *MultiHandler) Enabled(ctx context.Context, l slog.Level) bool {
	for i := range h.handlers {
		if h.handlers[i].Enabled(ctx, l) {
			return true
		}
	}

	return false
}

func try(callback func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = e
			} else {
				err = fmt.Errorf("unexpected error: %+v", r)
			}
		}
	}()

	err = callback()

	return
}

func (h *MultiHandler) Handle(ctx context.Context, r slog.Record) error {
	for i := range h.handlers {
		if h.handlers[i].Enabled(ctx, r.Level) {
			err := try(func() error {
				return h.handlers[i].Handle(ctx, r.Clone())
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (h *MultiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	handers := lo.Map(h.handlers, func(h slog.Handler, _ int) slog.Handler {
		return h.WithAttrs(attrs)
	})
	return Multi(handers...)
}

func (h *MultiHandler) WithGroup(name string) slog.Handler {
	handers := lo.Map(h.handlers, func(h slog.Handler, _ int) slog.Handler {
		return h.WithGroup(name)
	})
	return Multi(handers...)
}

var _ io.WriteCloser = (*multiWriteCloser)(nil)

func MultiWriteCloser(wrs ...io.WriteCloser) io.WriteCloser {
	return &multiWriteCloser{wrs: wrs}
}

// multiWriteCloser implements io.WriteCloser
type multiWriteCloser struct {
	wrs []io.WriteCloser
}

func (m *multiWriteCloser) Close() error {
	var joinErr error
	for _, wr := range m.wrs {
		if err := wr.Close(); err != nil {
			joinErr = errors.Join(joinErr, err)
		}
	}
	return joinErr
}

func (m *multiWriteCloser) Write(p []byte) (n int, err error) {
	for _, wr := range m.wrs {
		n, err := wr.Write(p)
		if err != nil {
			return n, err
		}
		if n != len(p) {
			return n, io.ErrShortWrite
		}
	}
	return len(p), nil
}
