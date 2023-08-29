package run

import (
	"whisper/pkg/config"
	"whisper/pkg/es"
	"whisper/pkg/log"
	"whisper/pkg/mongo"
	"whisper/pkg/mq"
	"whisper/pkg/mysql"
	"whisper/pkg/redis"

	mq2 "whisper/internal/service/mq"
)

func Init() {
	config.Init()
	log.Init()
	mysql.Init()
	es.Init()
	mongo.Init()
	mq.Init()
	redis.Init()

	consumerInit()
}

func consumerInit() {
	for _, f := range mq2.ConsumerFunc {
		if err := f(); err != nil {
			panic(err)
		}
	}
}
