package log

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

type GormLogger struct {
	logger *zap.SugaredLogger

	LogLevel                  logger.LogLevel
	SlowThreshold             time.Duration
	Colorful                  bool
	IgnoreRecordNotFoundError bool
	ParameterizedQueries      bool
}

func (g *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return logger.Default.LogMode(level)
}

func (glog *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if glog.LogLevel >= logger.Info {
		tpl := genLogTpl(ctx, infoStr+msg, data)
		glog.logger.Infof(tpl, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

func (glog *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if glog.LogLevel >= logger.Warn {
		tpl := genLogTpl(ctx, warnStr+msg, data)
		glog.logger.Warnf(tpl, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

func (glog *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if glog.LogLevel >= logger.Error {
		tpl := genLogTpl(ctx, errStr+msg, data)
		glog.logger.Errorf(tpl, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

func (glog *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if glog.LogLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	switch {
	case err != nil && glog.LogLevel >= logger.Error && (!errors.Is(err, logger.ErrRecordNotFound) || !glog.IgnoreRecordNotFoundError):
		sql, rows := fc()
		if rows == -1 {
			glog.logger.Errorf(traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			glog.logger.Errorf(traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case elapsed > glog.SlowThreshold && glog.SlowThreshold != 0 && glog.LogLevel >= logger.Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", glog.SlowThreshold)
		if rows == -1 {
			glog.logger.Warnf(traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			glog.logger.Warnf(traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case glog.LogLevel == logger.Info:
		sql, rows := fc()
		tpl := genLogTpl(ctx, traceStr, nil)
		if rows == -1 {
			glog.logger.Infof(tpl, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			glog.logger.Infof(tpl, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}

// newGormLogger GormLogger 初始化
func newGormLogger() CreateLoggerFunc {
	return func(lg *zap.SugaredLogger) interface{} {
		return &GormLogger{
			logger:   lg,
			LogLevel: logger.Info,

			SlowThreshold:             0,
			Colorful:                  false,
			IgnoreRecordNotFoundError: false,
			ParameterizedQueries:      false,
		}
	}
}
