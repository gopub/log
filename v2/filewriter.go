package log

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/natefinch/lumberjack.v2"
)

var _ io.WriteCloser = (*FileWriter)(nil)

type FileWriter struct {
	ll *lumberjack.Logger
}

func NewFileWriter(filename string) *FileWriter {
	filename = strings.TrimSpace(filename)
	if filename == "" {
		filename = filepath.Join(os.Args[0], "log")
	}

	ll := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    512, // 512 megabytes
		MaxAge:     30,  // 30 days
		MaxBackups: 32,
		LocalTime:  true,
		Compress:   true,
	}

	return &FileWriter{
		ll: ll,
	}
}

func (f FileWriter) Close() error {
	return f.ll.Close()
}

func (f FileWriter) Write(p []byte) (n int, err error) {
	return f.ll.Write(p)
}

func (f *FileWriter) SetMaxSize(sizeInMB int) {
	f.ll.MaxSize = sizeInMB
}

func (f *FileWriter) SetMaxAge(days int) {
	f.ll.MaxAge = days
}

func (f *FileWriter) SetMaxBackups(n int) {
	f.ll.MaxBackups = n
}
