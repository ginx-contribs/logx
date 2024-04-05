# logx
logx is a simple logging manager based on slog.Logger, supports multi handler and log cutting.

## Install
```bash
go get github.com/ginx-contribs/logx@latest
```

## Usage
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