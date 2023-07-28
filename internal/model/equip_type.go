package model

import (
	"time"
)

type EquipType struct {
	Id       uint64    `gorm:"column:id;primary_key;AUTO_INCREMENT;NOT NULL"`
	Platform int       `gorm:"column:platform;default:0;NOT NULL"`
	ItemId   string    `gorm:"column:item_id;NOT NULL;comment:'装备id'"`
	Types    string    `gorm:"column:types;NOT NULL;comment:'装备类型'"`
	Ctime    time.Time `gorm:"column:ctime;default:current_timestamp();NOT NULL"`
	Utime    time.Time `gorm:"column:utime;default:current_timestamp();NOT NULL"`
}

func (e *EquipType) TableName() string {
	return "equip_type"
}
