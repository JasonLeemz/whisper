package log

import (
	"go.uber.org/zap"
	"sync"
)

type MongoLogger struct {
	logger *zap.SugaredLogger
	mu     sync.Mutex
}

// newMongoLogger MongoLogger 初始化
func newMongoLogger() CreateLoggerFunc {
	return func(logger *zap.SugaredLogger) interface{} {
		return &MongoLogger{
			logger: logger,
			mu:     sync.Mutex{},
		}
	}
}

func (lg *MongoLogger) Info(level int, message string, keysAndValues ...interface{}) {
	lg.logger.Infof("level=%d msg=%s keysAndValues:%v", level, message, keysAndValues)
}

func (lg *MongoLogger) Error(err error, message string, keysAndValues ...interface{}) {
	lg.logger.Errorf("err=%d msg=%s keysAndValues:%v", err, message, keysAndValues)
}
