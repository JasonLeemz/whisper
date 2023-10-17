package log

import (
	"go.uber.org/zap"
	"gorm.io/gorm/utils"
	"whisper/pkg/context"
)

type WhisperLogger struct {
	*zap.SugaredLogger
}

func newZapLogger() CreateLoggerFunc {
	return func(logger *zap.SugaredLogger) interface{} {
		return &WhisperLogger{
			logger,
		}
	}
}

func (log *WhisperLogger) Debug(ctx *context.Context, data ...interface{}) {
	tpl := genLogTpl(ctx, "", data)
	log.Debugf(tpl, append([]interface{}{utils.FileWithLineNum()}, data...)...)
}

func (log *WhisperLogger) Info(ctx *context.Context, data ...interface{}) {
	tpl := genLogTpl(ctx, "", data)
	log.Infof(tpl, append([]interface{}{utils.FileWithLineNum()}, data...)...)
}

func (log *WhisperLogger) Warn(ctx *context.Context, data ...interface{}) {
	tpl := genLogTpl(ctx, "", data)
	log.Warnf(tpl, append([]interface{}{utils.FileWithLineNum()}, data...)...)
}

func (log *WhisperLogger) Error(ctx *context.Context, data ...interface{}) {
	tpl := genLogTpl(ctx, "", data)
	log.Errorf(tpl, append([]interface{}{utils.FileWithLineNum()}, data...)...)
}
