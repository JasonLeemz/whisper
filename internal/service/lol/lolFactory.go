package service

import (
	"time"
	"whisper/internal/service/common"
	"whisper/pkg/context"
)

type ILOL interface {
	// QueryEquipments 查询所有装备列表
	QueryEquipments(ctx *context.Context) (interface{}, error)
	// QueryHeroes 查询所有英雄
	QueryHeroes(ctx *context.Context) (interface{}, error)
	// QueryRune 查询所有天赋符文
	QueryRune(ctx *context.Context) (interface{}, error)
	// QuerySkill 查询所有召唤师技能
	QuerySkill(ctx *context.Context) (interface{}, error)
	// GetHeroAttribute 查询英雄属性
	GetHeroAttribute(ctx *context.Context, heroID string) (interface{}, error)
	// QueryRuneType 查询天赋符文分类
	QueryRuneType(ctx *context.Context) (interface{}, error)
	// QuerySuitEquip 查询英雄适配装备
	QuerySuitEquip(ctx *context.Context, heroID string) (interface{}, error)
	// HeroRankData 英雄对战数据详情(ChampionFightData)
	HeroRankData(ctx *context.Context, heroID string) (interface{}, error)
	// HeroRankList 手游各位置英雄胜率
	HeroRankList(ctx *context.Context) (interface{}, error)
	// VersionList 手游版本列表
	VersionList(ctx *context.Context) (interface{}, error)
	// VersionDetail 版本更新详情
	VersionDetail(ctx *context.Context, keys []string) (interface{}, error)
	// VersionInfo 版本更新了哪些类别
	VersionInfo(ctx *context.Context, vKey, id string) (interface{}, error)
}

type LOLProductFunc func() ILOL

func CreateLOLProduct(platform int) LOLProductFunc {
	if platform == common.PlatformForLOL {
		return func() ILOL {
			return &LOL{
				platform:      common.PlatformForLOL,
				ts:            time.Now().Unix() / 600,
				yesterday:     time.Now().AddDate(0, 0, -1).Format("20060102"),
				referer:       common.Referer,
				cookie:        common.Cookie,
				commonHeaders: common.CommonHeaders(),
			}
		}
	} else {
		return func() ILOL {
			return &LOLM{
				platform:      common.PlatformForLOLM,
				ts:            time.Now().Unix() / 600,
				referer:       common.Referer,
				cookie:        common.Cookie,
				commonHeaders: common.CommonHeaders(),
			}
		}
	}
}
