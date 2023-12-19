package log

import (
	"fmt"
	"go.uber.org/zap"
)

type ESLogger struct {
	*zap.SugaredLogger
}

// newESLogger ESLogger 初始化
func newESLogger() CreateLoggerFunc {
	return func(logger *zap.SugaredLogger) interface{} {
		return &ESLogger{
			logger,
		}
	}
}

func (lg ESLogger) Printf(format string, v ...interface{}) {
	lg.Info(fmt.Sprintf(format, v...))
}
