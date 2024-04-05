# logx
logx is a simple logging manager based on slog.Logger, supports multi handler and log cutting.

## Install
```bash
go get github.com/ginx-contribs/logx@latest
```

## Usage
quick start
```go
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

```

with logx handler
```go
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

```

with multi others handler
```go
package main

import (
	"github.com/getsentry/sentry-go"
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

	sentryHandler := slogsentry.Option{Level: slog.LevelDebug}.NewSentryHandler()

	logger, err := logx.New(
		logx.WithHandlers(slog.NewJSONHandler(os.Stderr, nil)),
		logx.WithHandlers(sentryHandler),
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

```