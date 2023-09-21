package model

import (
	"time"
)

type HeroAttribute struct {
	Id                  uint64    `gorm:"column:id;primary_key;AUTO_INCREMENT;NOT NULL"`
	HeroId              string    `gorm:"column:heroId;default:;NOT NULL"`
	Title               string    `gorm:"column:title;default:;NOT NULL"`
	Name                string    `gorm:"column:name;default:;NOT NULL"`
	Alias               string    `gorm:"column:alias;default:;NOT NULL"`
	ShortBio            string    `gorm:"column:shortBio;default:;NOT NULL"`
	Defense             string    `gorm:"column:defense;default:;NOT NULL"`
	Magic               string    `gorm:"column:magic;default:;NOT NULL"`
	Difficulty          string    `gorm:"column:difficulty;default:;NOT NULL"`
	DifficultyL         string    `gorm:"column:difficultyL;default:;NOT NULL"`
	Attack              string    `gorm:"column:attack;default:;NOT NULL"`
	Attackrange         string    `gorm:"column:attackrange;default:;NOT NULL"`
	Attackdamage        string    `gorm:"column:attackdamage;default:;NOT NULL"`
	Attackspeed         string    `gorm:"column:attackspeed;default:;NOT NULL"`
	Attackspeedperlevel string    `gorm:"column:attackspeedperlevel;default:;NOT NULL"`
	Hp                  string    `gorm:"column:hp;default:;NOT NULL"`
	Hpperlevel          string    `gorm:"column:hpperlevel;default:;NOT NULL"`
	Mp                  string    `gorm:"column:mp;default:;NOT NULL"`
	Mpperlevel          string    `gorm:"column:mpperlevel;default:;NOT NULL"`
	Movespeed           string    `gorm:"column:movespeed;default:;NOT NULL"`
	Armor               string    `gorm:"column:armor;default:;NOT NULL"`
	Armorperlevel       string    `gorm:"column:armorperlevel;default:;NOT NULL"`
	Spellblock          string    `gorm:"column:spellblock;default:;NOT NULL"`
	Spellblockperlevel  string    `gorm:"column:spellblockperlevel;default:;NOT NULL"`
	Hpregen             string    `gorm:"column:hpregen;default:;NOT NULL"`
	Hpregenperlevel     string    `gorm:"column:hpregenperlevel;default:;NOT NULL"`
	Mpregen             string    `gorm:"column:mpregen;default:;NOT NULL"`
	Mpregenperlevel     string    `gorm:"column:mpregenperlevel;default:;NOT NULL"`
	Crit                string    `gorm:"column:crit;default:;NOT NULL"`
	Damage              string    `gorm:"column:damage;default:;NOT NULL"`
	Durability          string    `gorm:"column:durability;default:;NOT NULL"`
	Mobility            string    `gorm:"column:mobility;default:;NOT NULL"`
	Avatar              string    `gorm:"column:avatar;default:;NOT NULL"`
	MainImg             string    `gorm:"column:mainImg;default:;NOT NULL"`
	Highlightprice      string    `gorm:"column:highlightprice;default:;NOT NULL"`
	GoldPrice           string    `gorm:"column:goldPrice;default:;NOT NULL"`
	Couponprice         string    `gorm:"column:couponprice;default:;NOT NULL"`
	IsWeekFree          string    `gorm:"column:isWeekFree;default:;NOT NULL"`
	Platform            int       `gorm:"column:platform;default:;NOT NULL"`
	Version             string    `gorm:"column:version;default:;NOT NULL"`
	FileTime            string    `gorm:"column:fileTime;default:;NOT NULL"`
	Ctime               time.Time `gorm:"column:ctime;default:current_timestamp();NOT NULL"`
	Utime               time.Time `gorm:"column:utime;default:current_timestamp();NOT NULL"`
}

type HeroAttrWithExt struct {
	HeroId              string `gorm:"column:heroId;default:;NOT NULL"`
	Title               string `gorm:"column:title;default:;NOT NULL"`
	Name                string `gorm:"column:name;default:;NOT NULL"`
	Alias               string `gorm:"column:alias;default:;NOT NULL"`
	ShortBio            string `gorm:"column:shortBio;default:;NOT NULL"`
	Defense             string `gorm:"column:defense;default:;NOT NULL"`
	Magic               string `gorm:"column:magic;default:;NOT NULL"`
	Difficulty          string `gorm:"column:difficulty;default:;NOT NULL"`
	DifficultyL         string `gorm:"column:difficultyL;default:;NOT NULL"`
	Attack              string `gorm:"column:attack;default:;NOT NULL"`
	Attackrange         string `gorm:"column:attackrange;default:;NOT NULL"`
	Attackdamage        string `gorm:"column:attackdamage;default:;NOT NULL"`
	Attackspeed         string `gorm:"column:attackspeed;default:;NOT NULL"`
	Attackspeedperlevel string `gorm:"column:attackspeedperlevel;default:;NOT NULL"`
	Hp                  string `gorm:"column:hp;default:;NOT NULL"`
	Hpperlevel          string `gorm:"column:hpperlevel;default:;NOT NULL"`
	Mp                  string `gorm:"column:mp;default:;NOT NULL"`
	Mpperlevel          string `gorm:"column:mpperlevel;default:;NOT NULL"`
	Movespeed           string `gorm:"column:movespeed;default:;NOT NULL"`
	Armor               string `gorm:"column:armor;default:;NOT NULL"`
	Armorperlevel       string `gorm:"column:armorperlevel;default:;NOT NULL"`
	Spellblock          string `gorm:"column:spellblock;default:;NOT NULL"`
	Spellblockperlevel  string `gorm:"column:spellblockperlevel;default:;NOT NULL"`
	Hpregen             string `gorm:"column:hpregen;default:;NOT NULL"`
	Hpregenperlevel     string `gorm:"column:hpregenperlevel;default:;NOT NULL"`
	Mpregen             string `gorm:"column:mpregen;default:;NOT NULL"`
	Mpregenperlevel     string `gorm:"column:mpregenperlevel;default:;NOT NULL"`
	Crit                string `gorm:"column:crit;default:;NOT NULL"`
	Damage              string `gorm:"column:damage;default:;NOT NULL"`
	Durability          string `gorm:"column:durability;default:;NOT NULL"`
	Mobility            string `gorm:"column:mobility;default:;NOT NULL"`
	Avatar              string `gorm:"column:avatar;default:;NOT NULL"`
	Highlightprice      string `gorm:"column:highlightprice;default:;NOT NULL"`
	GoldPrice           string `gorm:"column:goldPrice;default:;NOT NULL"`
	Couponprice         string `gorm:"column:couponprice;default:;NOT NULL"`
	IsWeekFree          string `gorm:"column:isWeekFree;default:;NOT NULL"`
	Platform            int    `gorm:"column:platform;default:;NOT NULL"`
	SkinId              string `gorm:"column:skinId;default:;NOT NULL"`
	HeroName            string `gorm:"column:heroName;default:;NOT NULL"`
	HeroTitle           string `gorm:"column:heroTitle;default:;NOT NULL"`
	SkinName            string `gorm:"column:skin_name;default:;NOT NULL"`
	IsBase              string `gorm:"column:isBase;default:;NOT NULL"`
	EmblemsName         string `gorm:"column:emblemsName;default:;NOT NULL"`
	Description         string `gorm:"column:description;default:;NOT NULL"`
	MainImg             string `gorm:"column:mainImg;default:;NOT NULL"`
	IconImg             string `gorm:"column:iconImg;default:;NOT NULL"`
	LoadingImg          string `gorm:"column:loadingImg;default:;NOT NULL"`
	VideoImg            string `gorm:"column:videoImg;default:;NOT NULL"`
	SourceImg           string `gorm:"column:sourceImg;default:;NOT NULL"`
}

func (h *HeroAttribute) TableName() string {
	return "hero_attribute"
}
