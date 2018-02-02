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
	logger := log.WithFields([]*log.Field{{Key: "userID", Value: 1}, {Key: "name", Value: "Tom"}})
	logger.Error("data not found")

	logger.WithFields([]*log.Field{{Key: "count", Value: 2}}).Infof("Try to post topic:%s", "Which is the best city")
}

func TestLogger_SetFlags(t *testing.T) {
	log.SetFlags(log.Lmicroseconds | log.Lfunction)
	log.Info("System started")
}
