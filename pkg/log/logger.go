package log

import (
	"whisper/pkg/config"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger 声明日志类全局变量
var (
	Logger    *WhisperLogger
	RpcLogger *RPCLogger
	GLogger   *GormLogger
	ELogger   *ESLogger
	MLogger   *MongoLogger
	//Xlogger *XXXLogger
)

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
	case loggerTypeGorm:
		return newGormLogger()
	case loggerTypeES:
		return newESLogger()
	case loggerTypeMongo:
		return newMongoLogger()
	case loggerTypeRPC:
		return newRPCLogger()
	//case loggerTypeXXX:
	//	return newXXXLogger()
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

	RpcLogger = CreateLogger(loggerTypeRPC)(
		createSugarLogger(
			getEncoder(),
			getLogWriter(config.GlobalConfig.Log.RpcLog),
			zapcore.Level(config.GlobalConfig.Log.LogLevel),
		),
	).(*RPCLogger)

	// gorm logger
	GLogger = CreateLogger(loggerTypeGorm)(
		createSugarLogger(
			getEncoder(),
			getLogWriter(config.GlobalConfig.Log.SqlLog),
			zapcore.Level(config.GlobalConfig.Log.LogLevel),
		),
	).(*GormLogger)
	// es logger
	ELogger = CreateLogger(loggerTypeES)(
		createSugarLogger(
			getEncoder(),
			getLogWriter(config.GlobalConfig.Log.EsLog),
			zapcore.Level(config.GlobalConfig.Log.LogLevel),
		),
	).(*ESLogger)
	// mongo logger
	MLogger = CreateLogger(loggerTypeMongo)(
		createSugarLogger(
			getEncoder(),
			getLogWriter(config.GlobalConfig.Log.MongoLog),
			zapcore.Level(config.GlobalConfig.Log.LogLevel),
		),
	).(*MongoLogger)

	//// XXX logger
	//XLogger = CreateLogger(loggerTypeXXX)(
	//	createSugarLogger(
	//		getEncoder(),
	//		getLogWriter(config.GlobalConfig.Log.MongoLog),
	//		zapcore.Level(config.GlobalConfig.Log.LogLevel),
	//	),
	//).(*XXXLogger)
}
