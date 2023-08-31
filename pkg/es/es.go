package es

import (
	"fmt"
	"github.com/olivere/elastic/v7"
	"whisper/pkg/config"
	"whisper/pkg/log"
)

var ESClient *elastic.Client

func Init() {
	var err error

	// 初始化一个连接
	ESClient, err = elastic.NewClient(
		elastic.SetURL(fmt.Sprintf("http://%s:%s", config.GlobalConfig.ES.Host, config.GlobalConfig.ES.Port)),
		elastic.SetTraceLog(log.ELogger),
		elastic.SetInfoLog(log.ELogger),
		elastic.SetErrorLog(log.ELogger),
		elastic.SetSniff(false),
	)

	if err != nil {
		panic(err)
	}
}
