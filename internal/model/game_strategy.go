package model

import "time"

type GameStrategy struct {
	Id          int64     `gorm:"column:id;primary_key;AUTO_INCREMENT;NOT NULL"`
	Platform    int       `gorm:"column:platform;default:0;NOT NULL;comment:'手游or端游'"`
	Source      int32     `gorm:"column:source;default:0;NOT NULL;comment:'视频来源'"`
	Author      string    `gorm:"column:author;default:;NOT NULL;comment:'作者'"`
	LinkUrl     string    `gorm:"column:link_url;default:;NOT NULL;comment:'视频链接'"`
	MainImage   string    `gorm:"column:main_image;default:;NOT NULL;comment:'主图'"`
	PublicDate  time.Time `gorm:"column:public_date;default:current_timestamp();NOT NULL;comment:'视频发布日期'"`
	Ctime       time.Time `gorm:"column:ctime;default:current_timestamp();NOT NULL"`
	Utime       time.Time `gorm:"column:utime;default:current_timestamp();NOT NULL"`
	UserProfile string    `gorm:"column:user_profile;default:{};NOT NULL"`
	Title       string    `gorm:"column:title;default:;NOT NULL;comment:'视频标题'"`
	Subtitle    string    `gorm:"column:subtitle;default:;NOT NULL;comment:'视频副标题'"`
	Status      int8      `gorm:"column:status;default:0;NOT NULL"`
	Bvid        string    `gorm:"column:bvid;default:0;NOT NULL"`
	Length      string    `gorm:"column:length;default:0;NOT NULL"`
	Played      int64     `gorm:"column:played;default:0;NOT NULL"`
	Hero        string    `gorm:"column:hero;default:0;NOT NULL"`
}

func (g *GameStrategy) TableName() string {
	return "game_strategy"
}
