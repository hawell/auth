package logger

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(config ZapConfig) (*zap.Logger, error) {
	level := zapcore.InfoLevel
	_ = level.UnmarshalText([]byte(config.GetLevel()))

	zapConfig := zap.Config{
		Level:       zap.NewAtomicLevelAt(level),
		Development: false,
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "message",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.EpochTimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{config.GetDestination()},
		ErrorOutputPaths: []string{config.GetDestination()},
	}
	return zapConfig.Build()
}

func MiddlewareFunc(logger *zap.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		logger.Info("",
			zap.String("method", ctx.Request.Method),
			zap.String("uri", ctx.Request.RequestURI),
			zap.String("query", ctx.Request.URL.RawQuery),
		)
		ctx.Next()
	}
}
