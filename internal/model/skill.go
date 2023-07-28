package model

import "time"

type LOLSkill struct {
	Id            uint64    `gorm:"column:id;primary_key;AUTO_INCREMENT;NOT NULL"`
	Name          string    `gorm:"column:name;default:;NOT NULL"`
	Description   string    `gorm:"column:description;default:;NOT NULL"`
	Summonerlevel string    `gorm:"column:summonerlevel;default:;NOT NULL"`
	Cooldown      string    `gorm:"column:cooldown;default:;NOT NULL"`
	Gamemode      string    `gorm:"column:gamemode;default:;NOT NULL"`
	Icon          string    `gorm:"column:icon;default:;NOT NULL"`
	Version       string    `gorm:"column:version;default:;NOT NULL"`
	FileTime      string    `gorm:"column:fileTime;default:;NOT NULL"`
	Ctime         time.Time `gorm:"column:ctime;default:current_timestamp();NOT NULL"`
	Utime         time.Time `gorm:"column:utime;default:current_timestamp();NOT NULL"`
}

func (s *LOLSkill) TableName() string {
	return "lol_skill"
}

// -------------------------------------

type LOLMSkill struct {
	Id        uint64    `gorm:"column:id;primary_key;AUTO_INCREMENT;NOT NULL"`
	SkillId   string    `gorm:"column:skillId;default:;NOT NULL"`
	Name      string    `gorm:"column:name;default:;NOT NULL"`
	IconPath  string    `gorm:"column:iconPath;default:;NOT NULL"`
	FuncDesc  string    `gorm:"column:funcDesc;default:;NOT NULL"`
	Cd        string    `gorm:"column:cd;default:;NOT NULL"`
	Video     string    `gorm:"column:video;default:;NOT NULL"`
	Unlocklv  string    `gorm:"column:unlocklv;default:;NOT NULL"`
	Mode      string    `gorm:"column:mode;default:;NOT NULL"`
	SortOrder string    `gorm:"column:sortOrder;default:;NOT NULL"`
	Version   string    `gorm:"column:version;default:;NOT NULL"`
	FileTime  string    `gorm:"column:fileTime;default:;NOT NULL"`
	Ctime     time.Time `gorm:"column:ctime;default:current_timestamp();NOT NULL"`
	Utime     time.Time `gorm:"column:utime;default:current_timestamp();NOT NULL"`
}

func (s *LOLMSkill) TableName() string {
	return "lolm_skill"
}
