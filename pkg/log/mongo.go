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

func (log *MongoLogger) Info(level int, message string, keysAndValues ...interface{}) {
	log.logger.Infof("level=%d msg=%s keysAndValues:%v", level, message, keysAndValues)
}

func (log *MongoLogger) Error(err error, message string, keysAndValues ...interface{}) {
	log.logger.Errorf("err=%d msg=%s keysAndValues:%v", err, message, keysAndValues)
}
