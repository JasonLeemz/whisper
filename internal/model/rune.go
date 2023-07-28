package model

import (
	"time"
)

type LOLRune struct {
	Id        uint64    `gorm:"column:id;primary_key;AUTO_INCREMENT;NOT NULL"`
	Name      string    `gorm:"column:name;default:;NOT NULL"`
	Icon      string    `gorm:"column:icon;default:;NOT NULL"`
	Key       string    `gorm:"column:key;default:;NOT NULL"`
	Tooltip   string    `gorm:"column:tooltip;default:;NOT NULL"`
	Shortdesc string    `gorm:"column:shortdesc;default:;NOT NULL"`
	Longdesc  string    `gorm:"column:longdesc"`
	SlotLabel string    `gorm:"column:slotLabel;default:;NOT NULL"`
	StyleName string    `gorm:"column:styleName;default:;NOT NULL"`
	Version   string    `gorm:"column:version;default:;NOT NULL"`
	FileTime  string    `gorm:"column:fileTime;default:;NOT NULL"`
	Ctime     time.Time `gorm:"column:ctime;default:current_timestamp();NOT NULL"`
	Utime     time.Time `gorm:"column:utime;default:current_timestamp();NOT NULL"`
}

func (l *LOLRune) TableName() string {
	return "lol_rune"
}

// --------------------------------

type LOLMRune struct {
	Id                   uint64    `gorm:"column:id;primary_key;AUTO_INCREMENT;NOT NULL"`
	RuneId               string    `gorm:"column:runeId;default:;NOT NULL"`
	Name                 string    `gorm:"column:name;default:;NOT NULL"`
	Description          string    `gorm:"column:description;default:;NOT NULL"`
	DetailInfo           string    `gorm:"column:detailInfo;default:;NOT NULL"`
	AttrName             string    `gorm:"column:attrName;default:;NOT NULL"`
	Type                 string    `gorm:"column:type;default:;NOT NULL"`
	IconPath             string    `gorm:"column:iconPath;default:;NOT NULL"`
	SortOrder            string    `gorm:"column:sortOrder;default:;NOT NULL"`
	UnlockLv             string    `gorm:"column:unlockLv;default:;NOT NULL"`
	PrimarySlotIndex     string    `gorm:"column:primarySlotIndex;primary_key;default:;NOT NULL"`
	PrimarySlotSortOrder string    `gorm:"column:primarySlotSortOrder;primary_key;default:;NOT NULL"`
	Version              string    `gorm:"column:version;default:;NOT NULL"`
	FileTime             string    `gorm:"column:fileTime;default:;NOT NULL"`
	Ctime                time.Time `gorm:"column:ctime;default:current_timestamp();NOT NULL"`
	Utime                time.Time `gorm:"column:utime;default:current_timestamp();NOT NULL"`
}

func (l *LOLMRune) TableName() string {
	return "lolm_rune"
}
