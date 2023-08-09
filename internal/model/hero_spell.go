package model

import (
	"time"
)

type HeroSpell struct {
	Id              uint64    `gorm:"column:id;primary_key;AUTO_INCREMENT;NOT NULL"`
	HeroId          string    `gorm:"column:heroId;default:;NOT NULL"`
	SpellKey        string    `gorm:"column:spellKey;default:;NOT NULL"`
	Sort            int       `gorm:"column:sort;default:;NOT NULL"` // 用于展示数据的时候排序
	Name            string    `gorm:"column:name;default:;NOT NULL"`
	Description     string    `gorm:"column:description;default:;NOT NULL"`
	AbilityIconPath string    `gorm:"column:abilityIconPath;default:;NOT NULL"`
	Detail          string    `gorm:"column:detail;NOT NULL"`
	Platform        int       `gorm:"column:platform;default:;NOT NULL"`
	Version         string    `gorm:"column:version;default:;NOT NULL"`
	FileTime        string    `gorm:"column:fileTime;default:;NOT NULL"`
	Ctime           time.Time `gorm:"column:ctime;default:current_timestamp();NOT NULL"`
	Utime           time.Time `gorm:"column:utime;default:current_timestamp();NOT NULL"`
}

func (h *HeroSpell) TableName() string {
	return "hero_spell"
}
