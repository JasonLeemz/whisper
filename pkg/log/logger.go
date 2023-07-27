package log

import (
	"context"
	"errors"
	"fmt"
	"time"
	"whisper/pkg/trace"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"whisper/pkg/config"
)

// Logger 声明日志类全局变量
var Logger *WhisperLogger

type WhisperLogger struct {
	*zap.SugaredLogger
}

type GormLogger struct {
	logger *zap.SugaredLogger

	LogLevel                  logger.LogLevel
	SlowThreshold             time.Duration
	Colorful                  bool
	IgnoreRecordNotFoundError bool
	ParameterizedQueries      bool
}

var (
	whisperTplStr = "%s\t "

	debugStr     = "%s\tDEBUG "
	infoStr      = "%s\tINFO "
	warnStr      = "%s\tWARN "
	errStr       = "%s\tERROR "
	traceStr     = "%s\t[%.3fms] [rows:%v] %s"
	traceWarnStr = "%s %s\t[%.3fms] [rows:%v] %s"
	traceErrStr  = "%s %s\t[%.3fms] [rows:%v] %s"
)

var Glogger *GormLogger

// newGormLogger GormLogger 初始化
func newGormLogger() {
	Glogger = &GormLogger{
		logger:   Logger.SugaredLogger,
		LogLevel: logger.Info,
	}

}

func (glog *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	lg := logger.Default.LogMode(level)
	return lg
}

func (glog *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if glog.LogLevel >= logger.Info {
		glog.logger.Infof(infoStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

func (glog *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if glog.LogLevel >= logger.Warn {
		glog.logger.Warnf(warnStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

func (glog *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if glog.LogLevel >= logger.Error {
		glog.logger.Errorf(errStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

func (glog *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {

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
		if rows == -1 {
			glog.logger.Infof(traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			glog.logger.Infof(traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}

//	func genLogTpl(traceID string, paramsNum int) string {
//		tpl := whisperTplStr + "|trace_id=" + traceID
//		for i := 0; i < paramsNum; i++ {
//			tpl += "|%#v"
//		}
//		return tpl
//	}
func genLogTpl(traceID string, data []interface{}) string {
	tpl := whisperTplStr + "|trace_id=" + traceID

	for _, v := range data {
		switch v.(type) {
		case string:
			tpl += "|%s"
		default:
			tpl += "|%#v"
		}
	}

	return tpl
}

func (log *WhisperLogger) Debug(ctx context.Context, data ...interface{}) {
	tpl := genLogTpl(ctx.Value(trace.TraceID).(string), data)
	log.Debugf(tpl, append([]interface{}{utils.FileWithLineNum()}, data...)...)
}

func (log *WhisperLogger) Info(ctx context.Context, data ...interface{}) {
	//column := make([]interface{}, 0, len(data)+2)
	//column = append(column, "|trace_id="+ctx.Value(trace.TraceID).(string))
	//column = append(column, "|file="+utils.FileWithLineNum()+"|")
	//column = append(column, data...)

	tpl := genLogTpl(ctx.Value(trace.TraceID).(string), data)
	log.Infof(tpl, append([]interface{}{utils.FileWithLineNum()}, data...)...)
}

func (log *WhisperLogger) Warn(ctx context.Context, data ...interface{}) {
	tpl := genLogTpl(ctx.Value(trace.TraceID).(string), data)
	log.Warnf(tpl, append([]interface{}{utils.FileWithLineNum()}, data...)...)
}

func (log *WhisperLogger) Error(ctx context.Context, data ...interface{}) {
	tpl := genLogTpl(ctx.Value(trace.TraceID).(string), data)
	log.Errorf(tpl, append([]interface{}{utils.FileWithLineNum()}, data...)...)
}

// Init 日志类初始化方法
func Init() {
	newZapLogger()
	// gorm logger
	newGormLogger()
}

// 日志记录地址
func getLogWriter() zapcore.WriteSyncer {
	//定义日志文件名，设置权限，当日志文件不存在时创建文件
	lumberJackLogger := &lumberjack.Logger{
		Filename:   config.GlobalConfig.Log.Path,
		MaxSize:    10, // 切割大小M
		MaxBackups: 3,  // 保留最大数量
		MaxAge:     1,
		Compress:   false,
		LocalTime:  true,
	}
	return zapcore.AddSync(lumberJackLogger)

}

// 日志编码方式
func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	// TODO 生产环境 和 开发环境
	//return zapcore.NewJSONEncoder(encoderConfig)
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func newZapLogger() {
	writeSyncer := getLogWriter()
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.Level(config.GlobalConfig.Log.LogLevel))
	zapLogger := zap.New(core, zap.AddCaller())
	Logger = &WhisperLogger{
		zapLogger.Sugar(),
	}
}
