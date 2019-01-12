package log

import "os"

func init() {
	defaultLogger = NewLogger(os.Stderr, AllLevel, LstdFlags)
}
