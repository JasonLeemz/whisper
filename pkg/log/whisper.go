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

func (lg *WhisperLogger) Debug(ctx *context.Context, data ...interface{}) {
	tpl := genLogTpl(ctx, "", data)
	lg.Debugf(tpl, append([]interface{}{utils.FileWithLineNum()}, data...)...)
}

func (lg *WhisperLogger) Info(ctx *context.Context, data ...interface{}) {
	tpl := genLogTpl(ctx, "", data)
	lg.Infof(tpl, append([]interface{}{utils.FileWithLineNum()}, data...)...)
}

func (lg *WhisperLogger) Warn(ctx *context.Context, data ...interface{}) {
	tpl := genLogTpl(ctx, "", data)
	lg.Warnf(tpl, append([]interface{}{utils.FileWithLineNum()}, data...)...)
}

func (lg *WhisperLogger) Error(ctx *context.Context, data ...interface{}) {
	tpl := genLogTpl(ctx, "", data)
	lg.Errorf(tpl, append([]interface{}{utils.FileWithLineNum()}, data...)...)
}
