package log

import (
	"context"
	"errors"
	"fmt"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"strconv"
	"sync"
	"time"
	"whisper/pkg/config"
	context2 "whisper/pkg/context"
	errors2 "whisper/pkg/errors"
)

// Logger 声明日志类全局变量
var Logger *WhisperLogger
var GLogger *GormLogger
var ELogger *ESLogger
var MLogger *MongoLogger

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

type MongoLogger struct {
	logger *zap.SugaredLogger
	mu     sync.Mutex
}

func (log *MongoLogger) Info(level int, message string, keysAndValues ...interface{}) {
	log.logger.Infof("level=%d msg=%s keysAndValues:%v", level, message, keysAndValues)
}

func (log *MongoLogger) Error(err error, message string, keysAndValues ...interface{}) {
	log.logger.Errorf("err=%d msg=%s keysAndValues:%v", err, message, keysAndValues)
}

type ESLogger struct {
	logger *zap.SugaredLogger
}

func (log ESLogger) Printf(format string, v ...interface{}) {
	log.logger.Info(fmt.Sprintf(format, v...))
}

var (
	whisperTplStr = "%s\t "

	infoStr      = "INFO "
	warnStr      = "WARN "
	errStr       = "ERROR "
	traceStr     = "[%.3fms] [rows:%v] %s"
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

	GLogger = &GormLogger{
		logger:   zapLogger.Sugar(),
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

	ELogger = &ESLogger{
		logger: zapLogger.Sugar(),
	}

}

// newMongoLogger MongoLogger 初始化
func newMongoLogger() {
	// 希望es的log单独打印日志，这里不复用Logger.SugaredLogger
	writeSyncer := getLogWriter(config.GlobalConfig.Log.MongoLog)
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.Level(config.GlobalConfig.Log.LogLevel))
	zapLogger := zap.New(core, zap.AddCaller())

	MLogger = &MongoLogger{
		logger: zapLogger.Sugar(),
		mu:     sync.Mutex{},
	}

}

func (glog *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	lg := logger.Default.LogMode(level)
	return lg
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
		tpl := genLogTpl(ctx, traceStr, nil)
		if rows == -1 {
			glog.logger.Infof(tpl, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			glog.logger.Infof(tpl, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}

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

func (log *WhisperLogger) Debug(ctx context.Context, data ...interface{}) {
	tpl := genLogTpl(ctx, "", data)
	log.Debugf(tpl, append([]interface{}{utils.FileWithLineNum()}, data...)...)
}

func (log *WhisperLogger) Info(ctx context.Context, data ...interface{}) {
	tpl := genLogTpl(ctx, "", data)
	log.Infof(tpl, append([]interface{}{utils.FileWithLineNum()}, data...)...)
}

func (log *WhisperLogger) Warn(ctx context.Context, data ...interface{}) {
	tpl := genLogTpl(ctx, "", data)
	log.Warnf(tpl, append([]interface{}{utils.FileWithLineNum()}, data...)...)
}

func (log *WhisperLogger) Error(ctx context.Context, data ...interface{}) {
	tpl := genLogTpl(ctx, "", data)
	log.Errorf(tpl, append([]interface{}{utils.FileWithLineNum()}, data...)...)
}

// Init 日志类初始化方法
func Init() {
	newZapLogger()
	// gorm logger
	newGormLogger()
	// es logger
	newESLogger()
	// mongo logger
	newMongoLogger()
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

func newZapLogger() {
	writeSyncer := getLogWriter(config.GlobalConfig.Log.Path)
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.Level(config.GlobalConfig.Log.LogLevel))
	zapLogger := zap.New(core, zap.AddCaller())
	Logger = &WhisperLogger{
		zapLogger.Sugar(),
	}
}
