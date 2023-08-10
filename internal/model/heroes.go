package model

import (
	"time"
)

type LOLHeroesEXT struct {
	Avatar      string `json:"avatar"`
	HeroId      string `gorm:"column:heroId;default:;NOT NULL"`
	Name        string `gorm:"column:name;default:;NOT NULL"`
	Alias       string `gorm:"column:alias;default:;NOT NULL"`
	Title       string `gorm:"column:title;default:;NOT NULL"`
	Roles       string `gorm:"column:roles;default:;NOT NULL"`
	Difficulty  string `gorm:"column:difficulty;default:;NOT NULL"`
	GoldPrice   string `gorm:"column:goldPrice;default:;NOT NULL"`
	CouponPrice string `gorm:"column:couponPrice;default:;NOT NULL"`
	Keywords    string `gorm:"column:keywords;default:;NOT NULL"`
	Version     string `gorm:"column:version;default:;NOT NULL"`
	FileTime    string `gorm:"column:fileTime;default:;NOT NULL"`
}

type LOLHeroes struct {
	Id                  uint64    `gorm:"column:id;primary_key;AUTO_INCREMENT;NOT NULL"`
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
	IsARAMWeekFree      string    `gorm:"column:isARAMweekfree;default:;NOT NULL"`
	IsPermanentWeekFree string    `gorm:"column:ispermanentweekfree;default:;NOT NULL"`
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

func (h *LOLHeroes) TableName() string {
	return "lol_heroes"
}

// -----------

type LOLMHeroesEXT struct {
	HeroId         string `gorm:"column:heroId;default:;NOT NULL"`
	Name           string `gorm:"column:name;default:;NOT NULL"`
	Title          string `gorm:"column:title;default:;NOT NULL"`
	Roles          string `gorm:"column:roles;default:;NOT NULL"`
	Intro          string `gorm:"column:intro;default:;NOT NULL"`
	Avatar         string `gorm:"column:avatar;default:;NOT NULL"`
	Card           string `gorm:"column:card;default:;NOT NULL"`
	Poster         string `gorm:"column:poster;default:;NOT NULL"`
	Highlightprice string `gorm:"column:highlightprice;default:;NOT NULL"`
	Couponprice    string `gorm:"column:couponprice;default:;NOT NULL"`
	Alias          string `gorm:"column:alias;default:;NOT NULL"`
	Lane           string `gorm:"column:lane;default:;NOT NULL"`
	Tags           string `gorm:"column:tags;default:;NOT NULL"`
	Searchkey      string `gorm:"column:searchkey;default:;NOT NULL"`
	DifficultyL    string `gorm:"column:difficultyL;default:;NOT NULL"`
	Version        string `gorm:"column:version;default:;NOT NULL"`
	FileTime       string `gorm:"column:fileTime;default:;NOT NULL"`
}

type LOLMHeroes struct {
	Id             uint64    `gorm:"column:id;primary_key;AUTO_INCREMENT;NOT NULL"`
	HeroId         string    `gorm:"column:heroId;default:;NOT NULL"`
	Name           string    `gorm:"column:name;default:;NOT NULL"`
	Title          string    `gorm:"column:title;default:;NOT NULL"`
	Roles          string    `gorm:"column:roles;default:;NOT NULL"`
	Intro          string    `gorm:"column:intro;default:;NOT NULL"`
	Avatar         string    `gorm:"column:avatar;default:;NOT NULL"`
	Card           string    `gorm:"column:card;default:;NOT NULL"`
	Poster         string    `gorm:"column:poster;default:;NOT NULL"`
	Highlightprice string    `gorm:"column:highlightprice;default:;NOT NULL"`
	Couponprice    string    `gorm:"column:couponprice;default:;NOT NULL"`
	Alias          string    `gorm:"column:alias;default:;NOT NULL"`
	Lane           string    `gorm:"column:lane;default:;NOT NULL"`
	Tags           string    `gorm:"column:tags;default:;NOT NULL"`
	Searchkey      string    `gorm:"column:searchkey;default:;NOT NULL"`
	IsWeekFree     string    `gorm:"column:isWeekFree;default:;NOT NULL"`
	DifficultyL    string    `gorm:"column:difficultyL;default:;NOT NULL"`
	Damage         string    `gorm:"column:damage;default:;NOT NULL"`
	SurviveL       string    `gorm:"column:surviveL;default:;NOT NULL"`
	AssistL        string    `gorm:"column:assistL;default:;NOT NULL"`
	Version        string    `gorm:"column:version;default:;NOT NULL"`
	FileTime       string    `gorm:"column:fileTime;default:;NOT NULL"`
	Ctime          time.Time `gorm:"column:ctime;default:current_timestamp();NOT NULL"`
	Utime          time.Time `gorm:"column:utime;default:current_timestamp();NOT NULL"`
}

func (l *LOLMHeroes) TableName() string {
	return "lolm_heroes"
}
