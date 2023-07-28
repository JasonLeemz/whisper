package model

import (
	"time"
)

type HeroRole struct {
	Id       uint64    `gorm:"column:id;primary_key;AUTO_INCREMENT;NOT NULL"`
	Platform int       `gorm:"column:platform;default:0;NOT NULL"`
	HeroId   string    `gorm:"column:hero_id;NOT NULL;comment:'英雄id'"`
	Role     string    `gorm:"column:role;NOT NULL;comment:'英雄角色'"`
	Ctime    time.Time `gorm:"column:ctime;default:current_timestamp();NOT NULL"`
	Utime    time.Time `gorm:"column:utime;default:current_timestamp();NOT NULL"`
}

func (h *HeroRole) TableName() string {
	return "hero_role"
}
