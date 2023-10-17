package log

import (
	"context"
	"strconv"
	"time"
	"whisper/pkg/config"

	context2 "whisper/pkg/context"
	errors2 "whisper/pkg/errors"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	LoggerTypeGorm  = "gorm"
	LoggerTypeES    = "es"
	LoggerTypeMongo = "mongo"
)

// Logger 声明日志类全局变量
var (
	Logger  *WhisperLogger
	GLogger *GormLogger
	ELogger *ESLogger
	MLogger *MongoLogger
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
		case errors2.Error:
			tpl += "|%+v"
		default:
			tpl += "|%#v"
		}
	}

	return tpl
}

// 日志记录地址
func getLogWriter(logPath string) zapcore.WriteSyncer {
	//定义日志文件名，设置权限，当日志文件不存在时创建文件
	lumberJackLogger := &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    10, // 切割大小M
		MaxBackups: 5,  // 保留最大数量
		MaxAge:     30,
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

func createSugarLogger(encoder zapcore.Encoder, writeSyncer zapcore.WriteSyncer, level zapcore.Level) *zap.SugaredLogger {
	return zap.New(
		zapcore.NewCore(
			encoder,
			writeSyncer,
			level,
		), zap.AddCaller()).Sugar()
}

type loggerType string
type CreateLoggerFunc func(logger *zap.SugaredLogger) interface{}

func CreateLogger(loggerType loggerType) CreateLoggerFunc {
	switch loggerType {
	case LoggerTypeGorm:
		return newGormLogger()
	case LoggerTypeES:
		return newESLogger()
	case LoggerTypeMongo:
		return newMongoLogger()
	default:
		return newZapLogger()
	}
}

// Init 日志类初始化方法
func Init() {
	Logger = CreateLogger("")(
		createSugarLogger(
			getEncoder(),
			getLogWriter(config.GlobalConfig.Log.Path),
			zapcore.Level(config.GlobalConfig.Log.LogLevel),
		),
	).(*WhisperLogger)

	// gorm logger
	GLogger = CreateLogger(LoggerTypeGorm)(
		createSugarLogger(
			getEncoder(),
			getLogWriter(config.GlobalConfig.Log.MongoLog),
			zapcore.Level(config.GlobalConfig.Log.LogLevel),
		),
	).(*GormLogger)
	// es logger
	ELogger = CreateLogger(LoggerTypeES)(
		createSugarLogger(
			getEncoder(),
			getLogWriter(config.GlobalConfig.Log.EsLog),
			zapcore.Level(config.GlobalConfig.Log.LogLevel),
		),
	).(*ESLogger)
	// mongo logger
	MLogger = CreateLogger(LoggerTypeMongo)(
		createSugarLogger(
			getEncoder(),
			getLogWriter(config.GlobalConfig.Log.MongoLog),
			zapcore.Level(config.GlobalConfig.Log.LogLevel),
		),
	).(*MongoLogger)
}
