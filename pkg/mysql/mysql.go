package mysql

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"whisper/pkg/config"
	"whisper/pkg/log"
)

var DB *gorm.DB

func Init() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.GlobalConfig.Database.Username,
		config.GlobalConfig.Database.Password,
		config.GlobalConfig.Database.Host,
		config.GlobalConfig.Database.Port,
		config.GlobalConfig.Database.DB,
	)
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: log.Glogger,
	})
	if err != nil {
		panic(err)
	}
}
