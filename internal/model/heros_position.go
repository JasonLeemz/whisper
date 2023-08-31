package model

import (
	"time"
)

type HeroesPosition struct {
	Id       uint64    `gorm:"column:id;primary_key;AUTO_INCREMENT;NOT NULL"`
	HeroId   string    `gorm:"column:heroId;default:;NOT NULL"`
	Pos      string    `gorm:"column:pos;default:;NOT NULL"`
	Platform int       `gorm:"column:platform;default:0;NOT NULL"`
	Version  string    `gorm:"column:version;default:;NOT NULL"`
	FileTime string    `gorm:"column:fileTime;default:;NOT NULL"`
	Ctime    time.Time `gorm:"column:ctime;default:current_timestamp();NOT NULL"`
	Utime    time.Time `gorm:"column:utime;default:current_timestamp();NOT NULL"`
}

func (h *HeroesPosition) TableName() string {
	return "heroes_position"
}
