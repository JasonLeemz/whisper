package log

import (
	"context"
	"errors"
	"fmt"
	"time"
	errors2 "whisper/pkg/errors"
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
var GLogger *GormLogger
var ELogger *ESLogger

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

type ESLogger struct {
	logger *zap.SugaredLogger
}

func (log ESLogger) Printf(format string, v ...interface{}) {
	log.logger.Info(fmt.Sprintf(format, v...))
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

// newGormLogger GormLogger 初始化
func newGormLogger() {
	// 希望sql的log单独打印日志，这里不复用Logger.SugaredLogger
	writeSyncer := getLogWriter(config.GlobalConfig.Log.SqlLog)
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.Level(config.GlobalConfig.Log.LogLevel))
	zapLogger := zap.New(core, zap.AddCaller())
	gormLogger := &WhisperLogger{
		zapLogger.Sugar(),
	}

	GLogger = &GormLogger{
		logger:   gormLogger.SugaredLogger,
		LogLevel: logger.Info,

		SlowThreshold:             0,
		Colorful:                  false,
		IgnoreRecordNotFoundError: false,
		ParameterizedQueries:      false,
	}

}

// newESLogger ESLogger 初始化
func newESLogger() {
	// 希望es的log单独打印日志，这里不复用Logger.SugaredLogger
	writeSyncer := getLogWriter(config.GlobalConfig.Log.EsLog)
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.Level(config.GlobalConfig.Log.LogLevel))
	zapLogger := zap.New(core, zap.AddCaller())
	esLogger := &WhisperLogger{
		zapLogger.Sugar(),
	}

	ELogger = &ESLogger{
		logger: esLogger.SugaredLogger,
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

func genLogTpl(traceID string, data []interface{}) string {
	tpl := whisperTplStr + "|trace_id=" + traceID

	for _, v := range data {
		switch v.(type) {
		case string:
			tpl += "|%s"
		case error:
			tpl += "|%+v"
		case errors2.Error:
			tpl += "|%+v"
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
	// es logger
	newESLogger()
}

// 日志记录地址
func getLogWriter(logPath string) zapcore.WriteSyncer {
	//定义日志文件名，设置权限，当日志文件不存在时创建文件
	lumberJackLogger := &lumberjack.Logger{
		Filename:   logPath,
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
	writeSyncer := getLogWriter(config.GlobalConfig.Log.Path)
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.Level(config.GlobalConfig.Log.LogLevel))
	zapLogger := zap.New(core, zap.AddCaller())
	Logger = &WhisperLogger{
		zapLogger.Sugar(),
	}
}
