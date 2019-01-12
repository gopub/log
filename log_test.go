package log_test

import (
	"github.com/gopub/log"
	"testing"
)

func TestLogger_Debug(t *testing.T) {
	log.Debug("This is a debug message")
	log.Infof("count:%d", 3)
}

func TestFieldLogger_WithFields(t *testing.T) {
	logger := log.With("userID", 1, "name", "Tom")
	logger.Error("data not found")

	logger.WithFields([]*log.Field{{Key: "count", Value: 2}}).Infof("Try to post topic:%s", "Which is the best city")
}

func TestLogger_SetFlags(t *testing.T) {
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lfunction | log.LUTC)
	log.Info("Log.Ldate | Log.Lmicroseconds | Log.Lfunction | Log.LUTC")
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lfunction)
	log.Info("Log.Ldate | Log.Lmicroseconds | Log.Lfunction")
}

func TestLogger_Derive(t *testing.T) {
	l := log.Default().Derive("SuperAPP")
	l.SetFlags(log.LstdFlags)
	l.Info("Let's go!!!")
}

func BenchmarkDebugf(b *testing.B) {
	log.Debugf("ShortMessage:i:%d,f:%f,s:%s,v:%v", 10, 10.999, "hello", "world")
}
