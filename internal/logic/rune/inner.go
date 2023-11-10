package rune

import (
	"fmt"
	"strconv"
	"whisper/internal/logic/common"
	"whisper/internal/model"
	dao "whisper/internal/model/DAO"
	"whisper/pkg/context"
	"whisper/pkg/redis"
)

type Inner struct {
	ctx      *context.Context
	platform int
}

func NewInnerIns(ctx *context.Context) *Inner {
	return &Inner{ctx: ctx}
}

// GetAll return map[string]*model.LOLRune|map[string]*model.LOLMRune
func (e *Inner) GetAll(platform int) interface{} {
	// 获取全部装备
	if platform == common.PlatformForLOL {
		d := dao.NewLOLRuneDAO()
		eVersion, _ := d.GetLOLRuneMaxVersion()
		data, _ := d.GetLOLRune(eVersion.Version)

		mrune := make(map[string]*model.LOLRune)
		for _, r := range data {
			key := fmt.Sprintf(redis.KeyCacheRune, "召唤师峡谷", strconv.Itoa(common.PlatformForLOL), r.RuneID)
			mrune[key] = r
		}
		return mrune
	} else {
		d := dao.NewLOLMRuneDAO()
		v, _ := d.GetLOLMRuneMaxVersion()
		data, _ := d.GetLOLMRune(v.Version)

		mrune := make(map[string]*model.LOLMRune)
		for _, r := range data {
			key := fmt.Sprintf(redis.KeyCacheRune, "召唤师峡谷", strconv.Itoa(common.PlatformForLOLM), r.RuneId)
			mrune[key] = r
		}
		return mrune
	}

}
