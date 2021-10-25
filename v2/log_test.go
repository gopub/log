package log_test

import (
	"github.com/gopub/log/v2"
	"io"
	"os"
	"testing"
)

func TestLog(t *testing.T) {
	log.Debug("This is a debug message")
	log.Debugf("5+9=%d", 5+9)
	l := log.With("user", "tom", "count", 3)
	l.Info("Hahaha")

	fw := log.NewFileWriter("testdata/a.log")
	mw := io.MultiWriter(os.Stderr, fw)
	zl := log.NewZapLogger(mw)
	l = log.NewLogger(zl)
	l.Warn("Test")
	err := l.Sync()
	if err != nil {
		t.Error(err)
	}

	l.Named("KKKK").Infof("LLLL")
	l.WithLevel(log.ErrorLevel).Errorf("zzz")
	l.WithLevel(log.ErrorLevel).Infof("zzz")
}
