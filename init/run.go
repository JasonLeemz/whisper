package run

import (
	"whisper/pkg/config"
	"whisper/pkg/log"
	"whisper/pkg/mysql"
)

func Init() {
	config.Init()
	log.Init()
	mysql.Init()
}
