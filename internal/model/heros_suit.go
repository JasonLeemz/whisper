package model

import (
	"time"
)

type HeroesSuit struct {
	Id     uint64    `gorm:"column:id;primary_key;AUTO_INCREMENT;NOT NULL"`
	HeroId string    `gorm:"column:hero_id;NOT NULL;comment:'英雄id'"`
	ItemId string    `gorm:"column:item_id;NOT NULL;comment:'装备id'"`
	Ctime  time.Time `gorm:"column:ctime;default:current_timestamp();NOT NULL"`
	Utime  time.Time `gorm:"column:utime;default:current_timestamp();NOT NULL"`
}

func (h *HeroesSuit) TableName() string {
	return "heroes_suit"
}
