package zap

import (
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestLogger(t *testing.T) {
	logger := NewLogger(*zap.New(zapcore.NewNopCore()))
	logger.Infof("Hi %s", "there")
	logger.Error("Bad wolf")
}
