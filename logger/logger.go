package logger

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const timeFormatStr = "2006-01-02"

var logger *zap.Logger

type LogParam struct {
	Level  zapcore.Level
	Msg    string
	Fields []zap.Field
}

var logQueue = make(chan LogParam, 100)

func newLoggerEncoder() zapcore.Encoder {
	encoderConf := zap.NewProductionEncoderConfig()
	encoderConf.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConf.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConf)
}

func newLogWriter(dir, name string) zapcore.WriteSyncer {
	fileName := fmt.Sprintf(`%s\%s-%s`,
		dir,
		time.Now().Format(timeFormatStr),
		name)
	var file *os.File
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		file, _ = os.Create(fileName)
	} else {
		file, _ = os.OpenFile(fileName, os.O_APPEND, os.ModePerm)
	}

	return zapcore.AddSync(file)
}

func Init(dir string) {
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel
	})

	failureLogger := newLogWriter(dir, "log_err")
	normalLogger := newLogWriter(dir, "log_out")

	encoder := newLoggerEncoder()

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, failureLogger, highPriority),
		zapcore.NewCore(encoder, normalLogger, lowPriority),
	)

	logger = zap.New(core)
	defer logger.Sync()

	logger.Info("logger.Init()")
}

func Serve(ctx context.Context) {
	defer func() {
		for p := range logQueue {
			record(p)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case p := <-logQueue:
			record(p)
		default:
		}
	}
}

func Loading(lvl zapcore.Level, msg string, fields ...zap.Field) {
	p := LogParam{
		Level:  lvl,
		Msg:    msg,
		Fields: fields,
	}
	logQueue <- p
}

func record(p LogParam) {
	switch p.Level {
	case zap.DebugLevel:
		logger.Debug(p.Msg, p.Fields...)
	case zap.InfoLevel:
		logger.Info(p.Msg, p.Fields...)
	case zap.WarnLevel:
		logger.Warn(p.Msg, p.Fields...)
	case zap.ErrorLevel:
		logger.Error(p.Msg, p.Fields...)
	case zap.DPanicLevel:
		logger.DPanic(p.Msg, p.Fields...)
	case zap.PanicLevel:
		logger.Panic(p.Msg, p.Fields...)
	case zap.FatalLevel:
		logger.Fatal(p.Msg, p.Fields...)
	default:
	}
}
