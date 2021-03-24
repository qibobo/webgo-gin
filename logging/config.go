package logging

import (
	"encoding/json"
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	loggerConfig = `{
        "level": "info",
        "development": false,
        "outputPaths": ["stdout"],
        "errorOutputPaths": ["stderr"],
        "encoding": "json",
        "encoderConfig": {
          "timeKey": "dt",
          "levelKey": "level",
          "nameKey": "logger",
          "callerKey": "caller",
          "messageKey": "message",
          "stacktraceKey": "stacktrace",
          "lineEnding": "",
          "levelEncoder": "",
          "timeEncoder": "rfc3339",
          "durationEncoder": "",
          "callerEncoder": ""
        }
      }`
)

type Config struct {
	LoggingConfig string
	LoggingLevel  *zapcore.Level
}

func NewLogger(logLevel string, opts ...zap.Option) (*zap.Logger, error) {
	logger, _, err := newLoggerFromConfig(logLevel, opts)
	if err != nil {
		return nil, err
	}
	return logger, nil
}
func NewSugaredLogger(logLevel string, opts ...zap.Option) (*zap.SugaredLogger, error) {
	logger, _, err := newLoggerFromConfig(logLevel, opts)
	if err != nil {
		return nil, err
	}
	return logger.Sugar(), nil
}

func newLoggerFromConfig(logLevel string, opts []zap.Option) (*zap.Logger, zap.AtomicLevel, error) {
	loggingCfg, err := zapConfigFromJSON(loggerConfig)
	if err != nil {
		return nil, zap.AtomicLevel{}, err
	}

	if logLevel != "" {
		if level, err := levelFromString(logLevel); err == nil {
			loggingCfg.Level = zap.NewAtomicLevelAt(*level)
		}
	}

	logger, err := loggingCfg.Build(opts...)
	if err != nil {
		return nil, zap.AtomicLevel{}, err
	}

	logger.Info("Successfully created the logger.")
	logger.Info("Logging level set to: " + loggingCfg.Level.String())
	return logger, loggingCfg.Level, nil
}

func zapConfigFromJSON(configJSON string) (*zap.Config, error) {
	loggingCfg := zap.Config{}
	if configJSON != "" {
		if err := json.Unmarshal([]byte(configJSON), &loggingCfg); err != nil {
			return nil, err
		}
	}
	return &loggingCfg, nil
}

func levelFromString(level string) (*zapcore.Level, error) {
	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(level)); err != nil {
		return nil, fmt.Errorf("invalid logging level: %v", level)
	}
	return &zapLevel, nil
}
