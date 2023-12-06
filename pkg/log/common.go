package log

import (
	"context"
	"strconv"
	"time"
	context2 "whisper/pkg/context"
	"whisper/pkg/errors"
)

const (
	loggerTypeGorm  = "gorm"
	loggerTypeES    = "es"
	loggerTypeMongo = "mongo"
	loggerTypeRPC   = "rpc"
)

var (
	whisperTplStr = "%s\t "

	infoStr      = "INFO "
	warnStr      = "WARN "
	errStr       = "ERROR "
	traceStr     = "[%.3fms] [rows:%v] %s"
	traceWarnStr = "%s %s\t[%.3fms] [rows:%v] %s"
	traceErrStr  = "%s %s\t[%.3fms] [rows:%v] %s"
)

// func genLogTpl(traceID string, startTime time.Time, data []interface{}) string {
func genLogTpl(ctx context.Context, msg string, data []interface{}) string {
	traceID := ""
	if tid, ok := ctx.Value(context2.TraceID).(string); ok {
		traceID = tid
	}
	proc := ""
	if st, ok := ctx.Value(context2.StartTime).(time.Time); ok {
		proc = strconv.FormatFloat(time.Since(st).Seconds(), 'f', -1, 64)
	}
	tpl := whisperTplStr + msg +
		"|trace_id=" + traceID +
		"|proc=" + proc
	for _, v := range data {
		switch v.(type) {
		case string:
			tpl += "|%s"
		case error:
			tpl += "|%+v"
		case errors.Error:
			tpl += "|%+v"
		default:
			tpl += "|%#v"
		}
	}

	return tpl
}
