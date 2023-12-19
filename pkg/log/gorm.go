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
	*zap.SugaredLogger

	LogLevel                  logger.LogLevel
	SlowThreshold             time.Duration
	Colorful                  bool
	IgnoreRecordNotFoundError bool
	ParameterizedQueries      bool
}

func (lg *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return logger.Default.LogMode(level)
}

func (lg *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if lg.LogLevel >= logger.Info {
		tpl := genLogTpl(ctx, infoStr+msg, data)
		lg.Infof(tpl, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

func (lg *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if lg.LogLevel >= logger.Warn {
		tpl := genLogTpl(ctx, warnStr+msg, data)
		lg.Warnf(tpl, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

func (lg *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if lg.LogLevel >= logger.Error {
		tpl := genLogTpl(ctx, errStr+msg, data)
		lg.Errorf(tpl, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

func (lg *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if lg.LogLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	switch {
	case err != nil && lg.LogLevel >= logger.Error && (!errors.Is(err, logger.ErrRecordNotFound) || !lg.IgnoreRecordNotFoundError):
		sql, rows := fc()
		if rows == -1 {
			lg.Errorf(traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			lg.Errorf(traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case elapsed > lg.SlowThreshold && lg.SlowThreshold != 0 && lg.LogLevel >= logger.Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", lg.SlowThreshold)
		if rows == -1 {
			lg.Warnf(traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			lg.Warnf(traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case lg.LogLevel == logger.Info:
		sql, rows := fc()
		tpl := genLogTpl(ctx, traceStr, nil)
		if rows == -1 {
			lg.Infof(tpl, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			lg.Infof(tpl, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}

// newGormLogger GormLogger 初始化
func newGormLogger() CreateLoggerFunc {
	return func(lg *zap.SugaredLogger) interface{} {
		return &GormLogger{
			SugaredLogger: lg,
			LogLevel:      logger.Info,

			SlowThreshold:             0,
			Colorful:                  false,
			IgnoreRecordNotFoundError: false,
			ParameterizedQueries:      false,
		}
	}
}
