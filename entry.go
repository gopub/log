package log

import "time"

type Field struct {
	Key   string
	Value interface{}
}

type Entry struct {
	Level    Level
	Time     time.Time
	File     string
	Line     string
	Function string
	Fields   []*Field
	Message  string

	Flags int
}
