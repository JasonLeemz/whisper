package model

import (
	"time"
)

type EquipmentDao interface {
	Add(equips []*Equipment) (int64, error)
}
type Equipment struct {
	Id          int64     `gorm:"column:id;primary_key;AUTO_INCREMENT;NOT NULL"`
	Platform    int       `gorm:"column:platform;default:0;NOT NULL"`
	ItemId      string    `gorm:"column:itemId;default:;NOT NULL"`
	Name        string    `gorm:"column:name;default:;NOT NULL"`
	IconPath    string    `gorm:"column:iconPath;default:;NOT NULL"`
	Price       string    `gorm:"column:price;default:;NOT NULL"`
	Description string    `gorm:"column:description;default:;NOT NULL"`
	Plaintext   string    `gorm:"column:plaintext;default:;NOT NULL"`
	Sell        string    `gorm:"column:sell;default:;NOT NULL"`
	Total       string    `gorm:"column:total;default:;NOT NULL"`
	SuitHeroId  string    `gorm:"column:suitHeroId;default:;NOT NULL"`
	Tag         string    `gorm:"column:tag;default:;NOT NULL"`
	Keywords    string    `gorm:"column:keywords;default:;NOT NULL"`
	Maps        string    `gorm:"column:maps;default:;NOT NULL"`
	From        string    `gorm:"column:from;default:;NOT NULL;comment:'合成自'"`
	Into        string    `gorm:"column:into;default:;NOT NULL;comment:'由谁合成'"`
	Types       string    `gorm:"column:types;default:;NOT NULL"`
	Version     string    `gorm:"column:version;default:;NOT NULL"`
	FileTime    string    `gorm:"column:fileTime;default:;NOT NULL"`
	Ctime       time.Time `gorm:"column:ctime;default:current_timestamp();NOT NULL"`
	Utime       time.Time `gorm:"column:utime;default:current_timestamp();NOT NULL"`
}

func (tbl *Equipment) TableName() string {
	return "equipment"
}
