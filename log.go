package husky

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Log interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Panic(args ...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Panicf(format string, args ...interface{})
	Skip(skip int) Log
	WithField(key, value string) Log
}

type _Log struct {
	ins *zap.Logger
}

type LogConfig struct {
	Level string `toml:"level"`
	Ansi  bool   `toml:"ansi"`
	Json  bool   `toml:"json"`
}

func (l *_Log) Debug(args ...interface{}) {
	l.ins.Sugar().Debugln(args...)
}

func (l *_Log) Info(args ...interface{}) {
	l.ins.Sugar().Infoln(args...)
}

func (l *_Log) Warn(args ...interface{}) {
	l.ins.Sugar().Warnln(args...)
}

func (l *_Log) Error(args ...interface{}) {
	l.ins.Sugar().Errorln(args...)
}

func (l *_Log) Panic(args ...interface{}) {
	l.ins.Sugar().Panicln(args...)
}

func (l *_Log) Debugf(format string, args ...interface{}) {
	l.ins.Sugar().Debugf(format, args...)
}

func (l *_Log) Infof(format string, args ...interface{}) {
	l.ins.Sugar().Infof(format, args...)
}

func (l *_Log) Warnf(format string, args ...interface{}) {
	l.ins.Sugar().Warnf(format, args...)
}

func (l *_Log) Errorf(format string, args ...interface{}) {
	l.ins.Sugar().Errorf(format, args...)
}

func (l *_Log) Panicf(format string, args ...interface{}) {
	l.ins.Sugar().Panicf(format, args...)
}

func (l *_Log) Skip(skip int) Log {
	return &_Log{l.ins.WithOptions(zap.AddCallerSkip(skip))}
}

func (l *_Log) WithField(key, value string) Log {
	return &_Log{l.ins.With(zap.String(key, value))}
}

var logIns Log

func InitLog(config *LogConfig) {
	level, err := zapcore.ParseLevel(config.Level)
	if err != nil {
		level = zapcore.DebugLevel
	}
	atom := zap.NewAtomicLevelAt(level)
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.MessageKey = "message"
	encoderConfig.LevelKey = "level"
	encoderConfig.TimeKey = "time"
	encoderConfig.CallerKey = "location"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeCaller = zapcore.FullCallerEncoder
	if config.Ansi && !config.Json {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}
	var encoder zapcore.Encoder
	if config.Json {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}
	l := zap.New(
		zapcore.NewCore(
			encoder,
			zapcore.Lock(os.Stdout),
			atom,
		),
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	)
	logIns = &_Log{ins: l}
}

func LogDebug(args ...interface{}) {
	logIns.Skip(1).Debug(args...)
}

func LogInfo(args ...interface{}) {
	logIns.Skip(1).Info(args...)
}

func LogWarn(args ...interface{}) {
	logIns.Skip(1).Warn(args...)
}

func LogError(args ...interface{}) {
	logIns.Skip(1).Error(args...)
}

func LogPanic(args ...interface{}) {
	logIns.Skip(1).Panic(args...)
}

func LogDebugf(format string, args ...interface{}) {
	logIns.Skip(1).Debugf(format, args...)
}

func LogInfof(format string, args ...interface{}) {
	logIns.Skip(1).Infof(format, args...)
}

func LogWarnf(format string, args ...interface{}) {
	logIns.Skip(1).Warnf(format, args...)
}

func LogErrorf(format string, args ...interface{}) {
	logIns.Skip(1).Errorf(format, args...)
}

func LogPanicf(format string, args ...interface{}) {
	logIns.Skip(1).Panicf(format, args...)
}

func LogSkip(skip int) Log {
	return logIns.Skip(skip)
}

func LogWithField(key, value string) Log {
	return logIns.Skip(1).WithField(key, value)
}
