package log

import (
	"runtime"
	"strings"
	"time"
)

type entry struct {
	Name     string
	Level    Level
	Time     time.Time
	File     string
	Line     int
	Function string
	Fields   []*Field
	Message  string
	Flags    int
}

func newEntry(flags int, level Level, name string, fields []*Field, message string, callDepth int) *entry {
	e := &entry{}

	if flags&Lname != 0 {
		e.Name = name
	}

	if flags&(Ltime|Ldate|Lmillisecond|Lmicroseconds) != 0 {
		e.Time = time.Now()
		if flags&LUTC != 0 {
			e.Time = e.Time.UTC()
		}
	}

	if flags&(Llongfile|Lshortfile|Lfunction) != 0 {
		function, file, line, _ := runtime.Caller(callDepth)
		e.Line = line
		if flags&(Llongfile|Lshortfile) != 0 {
			if flags&Lshortfile != 0 {
				file = ShortPath(file)
			} else {
				file = RelativePath(file)
			}
		} else {
			file = ""
		}

		e.File = file
		if flags&Lfunction != 0 {
			e.Function = runtime.FuncForPC(function).Name()
			if len(file) > 0 {
				i := strings.LastIndex(e.Function, ".")
				if i >= 0 {
					e.Function = e.Function[i+1:]
				}
			}
		}
	}

	e.Flags = flags
	e.Message = message
	e.Fields = fields
	e.Level = level
	return e
}

func RelativePath(path string) string {
	if len(PackagePath) > 0 {
		return strings.TrimPrefix(path, PackagePath)
	}
	start := strings.Index(path, GoSrc)
	if start > 0 {
		start += len(GoSrc)
		path = path[start:]
	}
	return path
}

func ShortPath(path string) string {
	path = RelativePath(path)
	names := strings.Split(path, "/")
	for i := 0; i < len(names)-1; i++ {
		if len(names[i]) > 0 {
			names[i] = names[i][0:1]
		}
	}
	return strings.Join(names, "/")
}
