package log

import (
	"io"
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

type EntryWriter interface {
	Write(entry *Entry, w io.Writer) error
}
