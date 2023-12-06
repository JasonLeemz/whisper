package log

import (
	"go.uber.org/zap"
	"gorm.io/gorm/utils"
	"whisper/pkg/context"
)

type RPCLogger struct {
	*zap.SugaredLogger
}

// newRPCLogger RPCLogger 初始化
func newRPCLogger() CreateLoggerFunc {
	return func(logger *zap.SugaredLogger) interface{} {
		return &RPCLogger{
			logger,
		}
	}
}

func (lg *RPCLogger) Debug(ctx *context.Context, data ...interface{}) {
	tpl := genLogTpl(ctx, "", data)
	lg.Debugf(tpl, append([]interface{}{utils.FileWithLineNum()}, data...)...)
}

func (lg *RPCLogger) Info(ctx *context.Context, data ...interface{}) {
	tpl := genLogTpl(ctx, "", data)
	lg.Infof(tpl, append([]interface{}{utils.FileWithLineNum()}, data...)...)
}

func (lg *RPCLogger) Warn(ctx *context.Context, data ...interface{}) {
	tpl := genLogTpl(ctx, "", data)
	lg.Warnf(tpl, append([]interface{}{utils.FileWithLineNum()}, data...)...)
}

func (lg *RPCLogger) Error(ctx *context.Context, data ...interface{}) {
	tpl := genLogTpl(ctx, "", data)
	lg.Errorf(tpl, append([]interface{}{utils.FileWithLineNum()}, data...)...)
}
