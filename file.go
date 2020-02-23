package log

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	logSuffix  = "log"
	maxLogSize = 64 * 1024 * 1024
)

type FileWriter struct {
	dir string

	file   *os.File
	format string
	size   int
	mu     sync.Mutex
}

func NewFileWriter(dir string) *FileWriter {
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		log.Fatalf("Make dir: %s, %v", dir, err)
	}
	format := "2006-01-02"
	f, err := createLogFile(dir, format)
	if err != nil {
		log.Fatalf("Create log file: %v", err)
	}
	fw := &FileWriter{
		dir:    dir,
		file:   f,
		format: format,
	}
	return fw
}

func (w *FileWriter) Write(p []byte) (int, error) {
	if w.file == nil {
		return 0, errors.New("no open file")
	}
	n, err := w.file.Write(p)
	w.size += n
	w.checkFileSize()
	return n, err
}

func (w *FileWriter) checkFileSize() {
	if w.size <= maxLogSize {
		return
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.size <= maxLogSize {
		return
	}
	newFile, err := createLogFile(w.dir, w.format)
	if err != nil {
		log.Printf("Create log file: %v\n", err)
		return
	}
	old := w.file
	w.size = 0
	w.file = newFile
	go func() {
		time.Sleep(time.Second)
		err = old.Close()
		if err != nil {
			log.Printf("Close file: %v\n", err)
		}
	}()
}

func (w *FileWriter) Close() error {
	err := w.file.Close()
	w.file = nil
	return err
}

func createLogFile(dir, format string) (*os.File, error) {
	d, err := os.Open(dir)
	if err != nil {
		return nil, fmt.Errorf("open dir %s: %w", dir, err)
	}
	l, err := d.Readdir(0)
	if err != nil {
		return nil, fmt.Errorf("read dir %s: %w", dir, err)
	}
	name := time.Now().Format(format)
	num := 1
	for _, fi := range l {
		s := fi.Name()
		if !strings.HasPrefix(s, name) {
			continue
		}
		s = s[len(name):]
		if !strings.HasSuffix(s, logSuffix) {
			continue
		}
		s = s[:len(s)-len(logSuffix)]
		if len(s) < 3 || s[0] != '.' || s[len(s)-1] != '.' { // expect .[0-9]+.
			continue
		}
		s = s[1 : len(s)-1]
		n, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			continue
		}
		if int(n) >= num {
			num = int(n + 1)
		}
	}
	name += fmt.Sprintf(".%d.%s", num, logSuffix)
	fullPath := path.Join(dir, name)
	f, err := os.OpenFile(fullPath, os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("open file %s: %file", fullPath, err)
	}
	return f, nil
}
