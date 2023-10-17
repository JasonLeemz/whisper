package log

import (
	"fmt"
	"go.uber.org/zap"
)

type XXXLogger struct {
	logger *zap.SugaredLogger
}

// newESLogger ESLogger 初始化
func newXXXLogger() CreateLoggerFunc {
	return func(logger *zap.SugaredLogger) interface{} {
		return &ESLogger{
			logger: logger,
		}
	}
}

func (lg XXXLogger) Printf(format string, v ...interface{}) {
	lg.logger.Info(fmt.Sprintf(format, v...))
}
