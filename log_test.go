package log_test

import (
	"github.com/gopub/log"
	"testing"
)

func TestLogger_Debug(t *testing.T) {
	log.Debug("This is a debug message")
}

func TestFieldLogger_WithFields(t *testing.T) {
	logger := log.WithFields([]*log.Field{{Key: "a", Value: 1}, {Key: "hello", Value: "world"}})
	logger.Info("good")
	logger.Error("wow")

	logger.WithFields([]*log.Field{{Key: "c", Value: 2}}).Infof("hello:%s", "wowowwow")
}
