package model

import (
	"time"
)

type EquipAlias struct {
	Id         uint64    `gorm:"column:id;primary_key;AUTO_INCREMENT;NOT NULL"`
	Name       string    `gorm:"column:name;default:;NOT NULL"`
	Keywords   string    `gorm:"column:keywords;default:;NOT NULL"`
	KeywordsPy string    `gorm:"column:keywords_py;default:;NOT NULL"`
	Platform   int8      `gorm:"column:platform;default:0;NOT NULL"`
	Ctime      time.Time `gorm:"column:ctime;default:current_timestamp();NOT NULL"`
	Utime      time.Time `gorm:"column:utime;default:current_timestamp();NOT NULL"`
}

func (e *EquipAlias) TableName() string {
	return "equip_alias"
}
