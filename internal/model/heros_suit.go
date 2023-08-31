package model

import (
	"time"
)

type HeroesSuit struct {
	Id       uint64    `gorm:"column:id;primary_key;AUTO_INCREMENT;NOT NULL"`
	HeroId   string    `gorm:"column:heroId;default:;NOT NULL"`
	Pos      string    `gorm:"column:pos;default:;NOT NULL"`
	Itemids  string    `gorm:"column:itemids;default:;NOT NULL"`
	Igamecnt int32     `gorm:"column:igamecnt;default:0;NOT NULL"`
	Wincnt   int32     `gorm:"column:wincnt;default:0;NOT NULL"`
	Winrate  int32     `gorm:"column:winrate;default:0;NOT NULL"`
	Allcnt   int32     `gorm:"column:allcnt;default:0;NOT NULL"`
	Showrate int32     `gorm:"column:showrate;default:0;NOT NULL"`
	Type     int32     `gorm:"column:type;default:0;NOT NULL;comment:'0:单件适合 1:鞋子 2:出门装 3:核心三件套'"`
	Platform int       `gorm:"column:platform;default:0;NOT NULL"`
	Version  string    `gorm:"column:version;default:;NOT NULL"`
	FileTime string    `gorm:"column:fileTime;default:;NOT NULL"`
	Ctime    time.Time `gorm:"column:ctime;default:current_timestamp();NOT NULL"`
	Utime    time.Time `gorm:"column:utime;default:current_timestamp();NOT NULL"`
}

func (h *HeroesSuit) TableName() string {
	return "heroes_suit"
}

func (h *HeroesSuit) TypeOther() int32 {
	return 0
}

func (h *HeroesSuit) TypeShoes() int32 {
	return 1
}

func (h *HeroesSuit) TypeOut() int32 {
	return 2
}

func (h *HeroesSuit) TypeCore() int32 {
	return 3
}
