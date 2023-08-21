package run

import (
	"whisper/pkg/config"
	"whisper/pkg/es"
	"whisper/pkg/log"
	"whisper/pkg/mongo"
	"whisper/pkg/mysql"
)

func Init() {
	config.Init()
	log.Init()
	mysql.Init()
	es.Init()
	mongo.Init()
}
