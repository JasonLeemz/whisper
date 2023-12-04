package model

import "time"

type AuthorSpace struct {
	Id           int64     `gorm:"column:id;primary_key;AUTO_INCREMENT;NOT NULL"`
	Name         string    `gorm:"column:name;NOT NULL"`
	Space        string    `gorm:"column:space;default:;NOT NULL"`
	VideoBaseUrl string    `gorm:"column:video_base_url;NOT NULL"`
	Source       int       `gorm:"column:source;default:0;NOT NULL"`
	Ctime        time.Time `gorm:"column:ctime;default:current_timestamp();NOT NULL"`
	Utime        time.Time `gorm:"column:utime;default:current_timestamp();NOT NULL"`
	Status       int       `gorm:"column:status;default:0;NOT NULL"`
	Platform     int       `gorm:"column:platform;default:0;NOT NULL"`
}

func (a *AuthorSpace) TableName() string {
	return "author_space"
}
