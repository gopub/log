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
	defaultRotateKeep = 30       // 30 days
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

func NewFileWriter(dir string) (*FileWriter, error) {
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return nil, fmt.Errorf("make dir: %w", err)
	}

	if d, err := os.Open(dir); err != nil {
		return nil, fmt.Errorf("open dir: %w", err)
	} else {
		dir = d.Name()
	}

	fw := &FileWriter{
		dir:        dir,
		rotateSize: defaultRotateSize,
		rotateKeep: defaultRotateKeep,
	}

	if err = fw.rotate(); err != nil {
		return nil, fmt.Errorf("rotate: %w", err)
	}
	return fw, nil
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
	if keep < 0 {
		keep = 0
	}
	w.rotateKeep = keep
	names, err := w.listRotateFileNames()
	if err != nil {
		log.Printf("List rotate filenames: %v\n", err)
		return
	}
	go w.keepFilesByDate(names, keep)
}

func (w *FileWriter) Write(p []byte) (int, error) {
	if w.file == nil {
		return 0, errors.New("no open file")
	}
	if err := w.rotate(); err != nil {
		return 0, fmt.Errorf("rotate: %w", err)
	}
	n, err := w.file.Write(p)
	if err != nil {
		return n, fmt.Errorf("write: %w", err)
	}
	w.size += n
	return n, nil
}

func (w *FileWriter) rotate() error {
	day := time.Now().Day()
	if w.size <= w.rotateSize && w.date != nil && day == w.date.Day() {
		return nil
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.size <= w.rotateSize && w.date != nil && day == w.date.Day() {
		return nil
	}

	names, err := w.listRotateFileNames()
	if err != nil {
		return fmt.Errorf("list rotate file names: %w", err)
	}
	date := time.Now()
	dateStr := date.Format(rotateDateFormat)
	num := w.nextFileNumber(dateStr, names)
	filePath := path.Join(w.dir, fmt.Sprintf("%s.%d.%s", dateStr, num, rotateSuffix))
	newFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("open file %s: %w", filePath, err)
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
	go w.keepFilesByDate(names, w.rotateKeep)
	//latestFile := path.Join(w.dir, "latest.log")
	//go func() {
	//	w.keepFilesByDate(names, w.rotateKeep)
	//	err := exec.Command("ln", "-sf", filePath, latestFile).Run()
	//		if err != nil {
	//		log.Printf("Link: %v", err)
	//	}
	//}()
	return nil
}

func (w *FileWriter) Close() error {
	err := w.file.Close()
	w.file = nil
	w.date = nil
	return err
}

func (w *FileWriter) listRotateFileNames() ([]string, error) {
	d, err := os.Open(w.dir)
	if err != nil {
		return nil, fmt.Errorf("open dir %s: %w", w.dir, err)
	}
	l, err := d.Readdir(0)
	if err != nil {
		return nil, fmt.Errorf("read dir %s: %w", w.dir, err)
	}
	var names []string
	for _, fi := range l {
		if !rotateNameRegex.MatchString(fi.Name()) {
			continue
		}
		names = append(names, fi.Name())
	}
	sort.Slice(names, func(i, j int) bool {
		return compareFileName(names[i], names[j])
	})
	return names, nil
}

func (w *FileWriter) nextFileNumber(date string, sortedNames []string) int {
	for _, name := range sortedNames {
		if !strings.HasPrefix(name, date) {
			return 1
		}
		s := name[len(date)+1 : len(name)-len(rotateSuffix)-1]
		n, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			log.Printf("Parse number %s: %v\n", s, err)
			continue
		}
		return int(n + 1)
	}
	return 1
}

func (w *FileWriter) keepFilesByDate(names []string, days int) {
	dateStr := time.Now().AddDate(0, 0, -days).Format(rotateDateFormat)
	for _, name := range names {
		if strings.Compare(dateStr, name) < 0 {
			continue
		}
		w.deleteFile(name)
	}
}

// Deprecated: use keepFilesByDate strategy
func (w *FileWriter) keepFilesByNum(sortedNames []string, num int) {
	if len(sortedNames) <= num {
		return
	}
	for _, name := range sortedNames[num:] {
		w.deleteFile(name)
	}
}

func (w *FileWriter) deleteFile(name string) {
	fullPath := path.Join(w.dir, name)
	err := os.Remove(fullPath)
	if err != nil {
		log.Printf("Remove %s: %v\n", fullPath, err)
	}
}

func compareFileName(a, b string) bool {
	dateLen := len(rotateDateFormat)
	// Compare date first
	if v := strings.Compare(a[:dateLen], b[:dateLen]); v != 0 {
		return v > 0
	}

	// Same date, compare digit length first
	if v := len(a) - len(b); v != 0 {
		return v > 0
	}

	// Same digit length, compare value
	return strings.Compare(a[dateLen:], b[dateLen:]) > 0
}
