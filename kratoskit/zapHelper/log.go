package zapHelper

import (
	"fmt"
	"os"
	
	"github.com/go-kratos/kratos/v2/log"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var _ log.Logger = (*ZapLogger)(nil)

type ZapLogger struct {
	log  *zap.Logger
	Sync func() error
}

// Logger 配置zap日志,将zap日志库引入
func Logger(alert bool, webhook string) log.Logger {
	// 配置zap日志库的编码器
	encoder := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  "stack",
		EncodeTime:     zapcore.RFC3339TimeEncoder,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
	return NewZapLogger(
		encoder,
		zap.NewAtomicLevelAt(zapcore.InfoLevel),
		alert,
		webhook,
		zap.AddStacktrace(
			zap.NewAtomicLevelAt(zapcore.ErrorLevel)),
		zap.AddCaller(),
		zap.AddCallerSkip(2),
		zap.Development(),
	)
}

// getLogWriter 日志自动切割，采用 lumberjack 实现的
func getLogWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "./log/zap.log", // 指定日志存储位置
		MaxSize:    10,              // 日志的最大大小（M）
		MaxBackups: 5,               // 日志的最大保存数量
		MaxAge:     30,              // 日志文件存储最大天数
		Compress:   false,           // 是否执行压缩
	}
	return zapcore.AddSync(lumberJackLogger)
}

// NewZapLogger return a zaplog logger.
func NewZapLogger(encConf zapcore.EncoderConfig, level zap.AtomicLevel, alert bool, webhook string, opts ...zap.Option) *ZapLogger {
	// 日志切割
	writeSyncer := getLogWriter()
	coreList := []zapcore.Core{
		// stdout and file output core
		zapcore.NewCore(
			zapcore.NewJSONEncoder(encConf),                                                       // 编码器配置
			zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(writeSyncer)), // 打印到控制台和文件
			level,                                                                                 // 日志打印级别（大于该级别即打印到指定终端或者文件）
		),
	}
	if alert {
		coreList = append(coreList,
			// alert core
			zapcore.NewCore(
				NewAlertEncoder(encConf, webhook), // 推送编码器配置
				zapcore.NewMultiWriteSyncer(),
				zap.NewAtomicLevelAt(zapcore.ErrorLevel), // 日志推送级别（大于该级别即发送推送通知）
			),
		)
	}
	zapLogger := zap.New(zapcore.NewTee(coreList[:]...), opts...)
	return &ZapLogger{log: zapLogger, Sync: zapLogger.Sync}
}

// Log zap log
func (l *ZapLogger) Log(level log.Level, keyvals ...interface{}) error {
	if len(keyvals) == 0 || len(keyvals)%2 != 0 {
		l.log.Warn(fmt.Sprint("keyvalues must appear in pairs: ", keyvals))
		return nil
	}
	var data []zap.Field
	for i := 0; i < len(keyvals); i += 2 {
		data = append(data, zap.Any(fmt.Sprint(keyvals[i]), keyvals[i+1]))
	}
	switch level {
	case log.LevelDebug:
		l.log.Debug("", data...)
	case log.LevelInfo:
		l.log.Info("", data...)
	case log.LevelWarn:
		l.log.Warn("", data...)
	case log.LevelError:
		l.log.Error("", data...)
	case log.LevelFatal:
		l.log.Fatal("", data...)
	}
	return nil
}
