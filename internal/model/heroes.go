package model

import (
	"time"
)

type Heroes struct {
	Id                  uint64    `gorm:"column:id;primary_key;AUTO_INCREMENT;NOT NULL"`
	Platform            int       `gorm:"column:platform;default:0;NOT NULL"`
	HeroId              string    `gorm:"column:heroId;default:;NOT NULL"`
	Name                string    `gorm:"column:name;default:;NOT NULL"`
	Alias               string    `gorm:"column:alias;default:;NOT NULL"`
	Title               string    `gorm:"column:title;default:;NOT NULL"`
	Roles               string    `gorm:"column:roles;default:;NOT NULL"`
	IsWeekFree          string    `gorm:"column:isWeekFree;default:;NOT NULL"`
	Attack              string    `gorm:"column:attack;default:;NOT NULL"`
	Defense             string    `gorm:"column:defense;default:;NOT NULL"`
	Magic               string    `gorm:"column:magic;default:;NOT NULL"`
	Difficulty          string    `gorm:"column:difficulty;default:;NOT NULL"`
	SelectAudio         string    `gorm:"column:selectAudio;default:;NOT NULL"`
	BanAudio            string    `gorm:"column:banAudio;default:;NOT NULL"`
	IsARAMweekfree      string    `gorm:"column:isARAMweekfree;default:;NOT NULL"`
	Ispermanentweekfree string    `gorm:"column:ispermanentweekfree;default:;NOT NULL"`
	ChangeLabel         string    `gorm:"column:changeLabel;default:;NOT NULL"`
	GoldPrice           string    `gorm:"column:goldPrice;default:;NOT NULL"`
	CouponPrice         string    `gorm:"column:couponPrice;default:;NOT NULL"`
	Camp                string    `gorm:"column:camp;default:;NOT NULL"`
	CampId              string    `gorm:"column:campId;default:;NOT NULL"`
	Keywords            string    `gorm:"column:keywords;default:;NOT NULL"`
	InstanceId          string    `gorm:"column:instance_id;default:;NOT NULL"`
	Version             string    `gorm:"column:version;default:;NOT NULL"`
	FileTime            string    `gorm:"column:fileTime;default:;NOT NULL"`
	Ctime               time.Time `gorm:"column:ctime;default:current_timestamp();NOT NULL"`
	Utime               time.Time `gorm:"column:utime;default:current_timestamp();NOT NULL"`
}

func (h *Heroes) TableName() string {
	return "heroes"
}
