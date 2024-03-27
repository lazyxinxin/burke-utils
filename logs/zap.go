package logs

import (
	"fmt"
	"sort"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Zaplogger *zap.Logger
var ZapSugarLogger *zap.SugaredLogger

// InitGlobalZapLogger 初始化zaplogger,并将日志接口放入到全局配置中
func InitGlobalZapLogger(level zapcore.Level) *zap.Logger {
	lg, err := CreateDefaultZapLogger(level)
	if err != nil {
		panic(fmt.Sprintf("init the default zap logger error:%s", err))
	}
	Zaplogger = lg
	return lg
}

// InitGlobalZapLogger 初始化zaplogger,并将日志接口放入到全局配置中
func InitGlobalZapSugarLogger(level zapcore.Level) *zap.SugaredLogger {
	lg, err := CreateDefaultZapLogger(level)
	if err != nil {
		panic(fmt.Sprintf("init the default zap logger error:%s", err))
	}
	ZapSugarLogger = lg.Sugar()
	return lg.Sugar()
}

// CreateDefaultZapLogger creates a logger with default zap configuration
func CreateDefaultZapLogger(level zapcore.Level) (*zap.Logger, error) {
	lcfg := DefaultZapLoggerConfig
	lcfg.Level = zap.NewAtomicLevelAt(level)
	c, err := lcfg.Build()
	if err != nil {
		return nil, err
	}
	return c, nil
}

// DefaultZapLoggerConfig defines default zap logger configuration.
var DefaultZapLoggerConfig = zap.Config{
	Level: zap.NewAtomicLevelAt(ConvertToZapLevel(DefaultLogLevel)),

	Development: false,
	Sampling: &zap.SamplingConfig{
		Initial:    100,
		Thereafter: 100,
	},

	Encoding: DefaultLogFormat,

	// copied from "zap.NewProductionEncoderConfig" with some updates
	EncoderConfig: zapcore.EncoderConfig{
		TimeKey:       "ts",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.LowercaseLevelEncoder,

		// Custom EncodeTime function to ensure we match format and precision of historic capnslog timestamps
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02T15:04:05.999999Z0700"))
		},

		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	},

	// Use "/dev/null" to discard all
	OutputPaths:      []string{"stderr"},
	ErrorOutputPaths: []string{"stderr"},
}

// MergeOutputPaths merges logging output paths, resolving conflicts.
func MergeOutputPaths(cfg zap.Config) zap.Config {
	outputs := make(map[string]struct{})
	for _, v := range cfg.OutputPaths {
		outputs[v] = struct{}{}
	}
	outputSlice := make([]string, 0)
	if _, ok := outputs["/dev/null"]; ok {
		// "/dev/null" to discard all
		outputSlice = []string{"/dev/null"}
	} else {
		for k := range outputs {
			outputSlice = append(outputSlice, k)
		}
	}
	cfg.OutputPaths = outputSlice
	sort.Strings(cfg.OutputPaths)

	errOutputs := make(map[string]struct{})
	for _, v := range cfg.ErrorOutputPaths {
		errOutputs[v] = struct{}{}
	}
	errOutputSlice := make([]string, 0)
	if _, ok := errOutputs["/dev/null"]; ok {
		// "/dev/null" to discard all
		errOutputSlice = []string{"/dev/null"}
	} else {
		for k := range errOutputs {
			errOutputSlice = append(errOutputSlice, k)
		}
	}
	cfg.ErrorOutputPaths = errOutputSlice
	sort.Strings(cfg.ErrorOutputPaths)

	return cfg
}
