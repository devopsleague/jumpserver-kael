package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var GlobalLogger *zap.Logger

func Setup() {
	cores := Zap.GetZapCores()
	logger := zap.New(zapcore.NewTee(cores...))
	GlobalLogger = logger.WithOptions(zap.AddCaller())
	zap.ReplaceGlobals(GlobalLogger)
}
