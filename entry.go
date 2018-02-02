package log

import (
	"io"
	"runtime"
	"strings"
	"time"
)

type Field struct {
	Key   string
	Value interface{}
}

type Entry struct {
	Level    Level
	Time     time.Time
	File     string
	Line     int
	Function string
	Fields   []*Field
	Message  string
	Flags    int
}

type EntryPrinter interface {
	Print(entry *Entry, w io.Writer) error
}

func MakeEntry(flags int, level Level, fields []*Field, message string, callDepth int) *Entry {
	entry := &Entry{}
	if flags&(Ltime|Ldate|Lmicroseconds) != 0 {
		entry.Time = time.Now()
		if flags&LUTC != 0 {
			entry.Time = entry.Time.UTC()
		}
	}

	if flags&(Llongfile|Lshortfile|Lfunction) != 0 {
		function, file, line, _ := runtime.Caller(callDepth)
		entry.Line = line
		if flags&(Llongfile|Lshortfile) != 0 {
			if len(PackagePath) > 0 {
				file = strings.TrimPrefix(file, PackagePath)
			} else {
				start := strings.Index(file, GoSrc)
				if start > 0 {
					start += len(GoSrc)
				}
				file = file[start:]
			}

			if flags&Lshortfile != 0 {
				names := strings.Split(file, "/")
				for i := 1; i < len(names)-1; i++ {
					names[i] = names[i][0:1]
				}
				file = strings.Join(names, "/")
			}
		} else {
			file = ""
		}

		entry.File = file
		if flags&Lfunction != 0 {
			entry.Function = runtime.FuncForPC(function).Name()
			if len(file) > 0 {
				i := strings.LastIndex(entry.Function, ".")
				if i >= 0 {
					entry.Function = entry.Function[i+1:]
				}
			}
		}
	}

	entry.Flags = flags
	entry.Message = message
	entry.Fields = fields
	entry.Level = level

	return entry
}
