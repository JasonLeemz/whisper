package log

import (
	"fmt"
	"go.uber.org/zap"
)

type ESLogger struct {
	logger *zap.SugaredLogger
}

// newESLogger ESLogger 初始化
func newESLogger() CreateLoggerFunc {
	return func(logger *zap.SugaredLogger) interface{} {
		return &ESLogger{
			logger: logger,
		}
	}
}

func (log ESLogger) Printf(format string, v ...interface{}) {
	log.logger.Info(fmt.Sprintf(format, v...))
}
