package log

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	rotateSuffix      = "log"
	rotateDateFormat  = "20060102"
	minRotateSize     = 1 << 20  // 1M
	defaultRotateSize = 64 << 20 // 64M
	defaultRotateKeep = 64
)

var rotateNameRegex = regexp.MustCompile("[0-9]{8}\\.[0-9]+\\." + rotateSuffix)

// FileWriter writes logs into files
// Example:
// fw := log.NewFileWriter("/var/logs/myapp")
// log.SetDefault(log.NewLogger(fw))
type FileWriter struct {
	dir string

	file       *os.File
	date       *time.Time
	size       int
	mu         sync.Mutex
	rotateSize int
	rotateKeep int
}

func NewFileWriter(dir string) *FileWriter {
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		log.Fatalf("Make dir: %s, %v", dir, err)
	}
	if d, err := os.Open(dir); err != nil {
		log.Fatalf("Open dir: %s, %v", dir, err)
	} else {
		dir = d.Name()
	}

	fw := &FileWriter{
		dir:        dir,
		rotateSize: defaultRotateSize,
		rotateKeep: defaultRotateKeep,
	}
	fw.rotate()
	if fw.file == nil {
		log.Fatalf("Cannot write logs under dir %s", dir)
	}
	return fw
}

func (w *FileWriter) RotateSize() int {
	return w.rotateSize
}

func (w *FileWriter) SetRotateSize(size int) {
	if size < minRotateSize {
		w.rotateSize = minRotateSize
	} else {
		w.rotateSize = size
	}
}

func (w *FileWriter) RotateKeep() int {
	return w.rotateKeep
}

func (w *FileWriter) SetRotateKeep(keep int) {
	if keep <= 0 {
		keep = 1
	}
	w.rotateKeep = keep
	names, err := w.getSortedNames()
	if err != nil {
		log.Printf("GetSortedNames: %v", err)
		return
	}
	go w.keepFiles(names, keep)
}

func (w *FileWriter) Write(p []byte) (int, error) {
	if w.file == nil {
		return 0, errors.New("no open file")
	}
	w.rotate()
	n, err := w.file.Write(p)
	if err != nil {
		return 0, fmt.Errorf("write: %w", err)
	}
	w.size += n
	return n, nil
}

func (w *FileWriter) rotate() {
	day := time.Now().Day()
	if w.size <= w.rotateSize && w.date != nil && day == w.date.Day() {
		return
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.size <= w.rotateSize && w.date != nil && day == w.date.Day() {
		return
	}

	names, err := w.getSortedNames()
	if err != nil {
		log.Printf("GetSortedNames: %v\n", err)
		return
	}
	date := time.Now()
	dateStr := date.Format(rotateDateFormat)
	num, err := w.parseNextNum(dateStr, names)
	if err != nil {
		log.Printf("ParseNextNum: %v\n", err)
		return
	}

	name := fmt.Sprintf("%s.%d.%s", dateStr, num, rotateSuffix)
	fullPath := path.Join(w.dir, name)
	newFile, err := os.OpenFile(fullPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("OpenFile %s: %v", fullPath, err)
	}
	old := w.file
	w.size = 0
	w.file = newFile
	w.date = &date
	if old != nil {
		err = old.Close()
		if err != nil {
			log.Printf("Close file: %v\n", err)
		}
	}
	go w.keepFiles(names, w.rotateKeep-1)
}

func (w *FileWriter) Close() error {
	err := w.file.Close()
	w.file = nil
	w.date = nil
	return err
}

func (w *FileWriter) getSortedNames() (rotateNameList, error) {
	d, err := os.Open(w.dir)
	if err != nil {
		return nil, fmt.Errorf("open dir %s: %w", w.dir, err)
	}
	l, err := d.Readdir(0)
	if err != nil {
		return nil, fmt.Errorf("read dir %s: %w", w.dir, err)
	}
	var names rotateNameList
	for _, fi := range l {
		if !rotateNameRegex.MatchString(fi.Name()) {
			continue
		}
		names = append(names, fi.Name())
	}
	sort.Sort(names)
	return names, nil
}

func (w *FileWriter) parseNextNum(date string, sortedNames []string) (int, error) {
	if len(sortedNames) == 0 {
		return 1, nil
	}

	latest := sortedNames[0]
	if !strings.HasPrefix(sortedNames[0], date) {
		return 1, nil
	}
	numPart := latest[len(date)+1 : len(latest)-len(rotateSuffix)-1]
	n, err := strconv.ParseInt(numPart, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("parseInt %s: %w", numPart, err)
	}
	return int(n + 1), nil
}

func (w *FileWriter) keepFiles(sortedNames []string, size int) {
	if len(sortedNames) <= size {
		return
	}
	for _, name := range sortedNames[size:] {
		fullPath := path.Join(w.dir, name)
		err := os.Remove(fullPath)
		if err != nil {
			log.Printf("Remove %s: %v\n", fullPath, err)
		}
	}
}

type rotateNameList []string

func (l rotateNameList) Len() int {
	return len(l)
}

func (l rotateNameList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l rotateNameList) Less(i, j int) bool {
	dateLen := len(rotateDateFormat)
	// Compare date first
	if v := strings.Compare(l[i][:dateLen], l[j][:dateLen]); v != 0 {
		return v > 0
	}

	// Same date, compare digit length first
	if v := len(l[i]) - len(l[j]); v != 0 {
		return v > 0
	}

	// Same digit length, compare value
	return strings.Compare(l[i][dateLen:], l[j][dateLen:]) > 0
}
