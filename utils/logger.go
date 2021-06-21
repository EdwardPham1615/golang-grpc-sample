package utils

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func NewLogger() *zap.Logger {
	env := viper.GetString("workspace.env")
	var logger *zap.Logger
	switch env {
	case "DEVELOPMENT":
		logger, _ = zap.NewDevelopment()
	default:
		logger, _ = zap.NewProduction()
	}
	return logger
}
