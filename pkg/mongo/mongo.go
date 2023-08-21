package mongo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
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
	clientOptions := options.Client().ApplyURI(uri)

	var err error
	Client, err = mongo.Connect(context.Background(), clientOptions)
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
	//.Collection("mycollection")

	// 创建要插入的文档
	document := bson.D{{"name", "John"}, {"age", 30}}

	// 插入文档
	_, err = Database.Collection("equipment").InsertOne(context.Background(), document)
	if err != nil {
		panic(err)
	}
}
