package model

import (
	"time"
)

type RuneType struct {
	Id       uint64    `gorm:"column:id;primary_key;AUTO_INCREMENT;NOT NULL"`
	Name     string    `gorm:"column:name;default:;NOT NULL"`
	SubType  string    `gorm:"column:sub_type;default:;NOT NULL"`
	Type     string    `gorm:"column:type;default:;NOT NULL"`
	Platform int8      `gorm:"column:platform;default:0;NOT NULL"`
	Ctime    time.Time `gorm:"column:ctime;default:current_timestamp();NOT NULL"`
	Utime    time.Time `gorm:"column:utime;default:current_timestamp();NOT NULL"`
}

func (r *RuneType) TableName() string {
	return "rune_type"
}
