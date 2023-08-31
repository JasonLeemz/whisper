package mongo

import (
	"context"
	"fmt"
	"time"
	"whisper/pkg/log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"whisper/pkg/config"
)

var (
	Client   *mongo.Client
	Database *mongo.Database
)

func Init() {
	// 创建连接对象
	uri := fmt.Sprintf("mongodb://%s:%d", config.GlobalConfig.Mongodb.Host, config.GlobalConfig.Mongodb.Port)
	clientOptions := options.Client().
		SetLoggerOptions(
			options.Logger().SetComponentLevel(
				1,
				options.LogLevel(2)).
				SetSink(log.MLogger),
		).
		ApplyURI(uri).
		SetMaxPoolSize(100)

	var err error
	timeoutCtx, cancelFunc := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancelFunc()

	Client, err = mongo.Connect(timeoutCtx, clientOptions)
	if err != nil {
		panic(err)
	}

	// 检查连接
	err = Client.Ping(context.Background(), nil)
	if err != nil {
		panic(err)
	}

	// 选择数据库和集合
	Database = Client.Database("whisper")
}
