package logx

import (
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"path/filepath"
)

// WriterOptions configures the logger writer
type WriterOptions struct {
	// output filename
	Filename string `mapstructure:"filename"`

	// whether to writer into os.Stderr
	DisableStderr bool `mapstructure:"disableStderr"`

	// whether to cut log
	Cut bool `mapstructure:"cut"`

	// MaxSize is the maximum size in megabytes of the log file before it gets
	// rotated. It defaults to 100 megabytes.
	MaxSize int `mapstructure:"maxSize"`

	// MaxAge is the maximum number of days to retain old log files based on the
	// timestamp encoded in their filename.  Note that a day is defined as 24
	// hours and may not exactly correspond to calendar days due to daylight
	// savings, leap seconds, etc. The default is not to remove old log files
	// based on age.
	MaxAge int `mapstructure:"maxage"`

	// MaxBackups is the maximum number of old log files to retain.  The default
	// is to retain all old log files (though MaxAge may still cause them to get
	// deleted.)
	MaxBackups int `mapstructure:"maxbackups"`

	// Compress determines if the rotated log files should be compressed
	// using gzip. The default is not to perform compression.
	Compress bool `mapstructure:"compress"`
}

// NewWriter returns a new append-only file writer, supporting log cutting.
func NewWriter(options *WriterOptions) (io.WriteCloser, error) {

	if options == nil {
		options = &WriterOptions{}
	}

	var wcs []io.WriteCloser

	if options.Filename != "" {
		var writerCloser io.WriteCloser
		// cut log
		if options.Cut {
			writerCloser = &lumberjack.Logger{
				Filename:   options.Filename,
				MaxSize:    options.MaxSize,
				MaxAge:     options.MaxAge,
				MaxBackups: options.MaxBackups,
				Compress:   options.Compress,
				LocalTime:  true,
			}
		} else {
			file, err := openFile(options.Filename, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
			if err != nil {
				return nil, err
			}
			writerCloser = file
		}
		wcs = append(wcs, writerCloser)
	}

	if !options.DisableStderr {
		wcs = append(wcs, os.Stderr)
	}

	return MultiWriteCloser(wcs...), nil
}

func openFile(name string, flag int, perm os.FileMode) (*os.File, error) {
	dir := filepath.Dir(name)
	if dir != "." && dir != "" {
		err := os.MkdirAll(dir, 0666)
		if err != nil {
			return nil, err
		}
	}
	return os.OpenFile(name, flag, perm)
}
