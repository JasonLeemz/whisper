package model

import (
	"time"
)

type HeroSkin struct {
	Id         uint64    `gorm:"column:id;primary_key;AUTO_INCREMENT;NOT NULL"`
	HeroId     string    `gorm:"column:heroId;default:;NOT NULL"`
	Title      string    `gorm:"column:title;default:;NOT NULL"`
	Name       string    `gorm:"column:name;default:;NOT NULL"`
	MainImg    string    `gorm:"column:mainImg;default:;NOT NULL"`
	IconImg    string    `gorm:"column:iconImg;default:;NOT NULL"`
	LoadingImg string    `gorm:"column:loadingImg;default:;NOT NULL"`
	VideoImg   string    `gorm:"column:videoImg;default:;NOT NULL"`
	SourceImg  string    `gorm:"column:sourceImg;default:;NOT NULL"`
	Platform   int       `gorm:"column:platform;default:0;NOT NULL"`
	Version    string    `gorm:"column:version;default:;NOT NULL"`
	FileTime   string    `gorm:"column:fileTime;default:;NOT NULL"`
	Ctime      time.Time `gorm:"column:ctime;default:current_timestamp();NOT NULL"`
	Utime      time.Time `gorm:"column:utime;default:current_timestamp();NOT NULL"`
}

func (h *HeroSkin) TableName() string {
	return "hero_skin"
}
