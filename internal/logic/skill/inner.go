package skill

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

// GetAll return map[string]*model.LOLSkill|map[string]*model.LOLMSkill
func (e *Inner) GetAll(platform int) interface{} {
	// 获取全部装备
	if platform == common.PlatformForLOL {
		d := dao.NewLOLSkillDAO()
		eVersion, _ := d.GetLOLSkillMaxVersion()
		data, _ := d.GetLOLSkill(eVersion.Version)

		mskill := make(map[string]*model.LOLSkill)
		for _, skill := range data {
			key := fmt.Sprintf(redis.KeyCacheSkill, "召唤师峡谷", strconv.Itoa(common.PlatformForLOL), skill.SkillID)
			mskill[key] = skill
		}
		return mskill
	} else {
		d := dao.NewLOLMSkillDAO()
		v, _ := d.GetLOLMSkillMaxVersion()
		data, _ := d.GetLOLMSkill(v.Version)
		mskill := make(map[string]*model.LOLMSkill)
		for _, skill := range data {
			key := fmt.Sprintf(redis.KeyCacheSkill, "召唤师峡谷", strconv.Itoa(common.PlatformForLOLM), skill.SkillID) // todo maps "召唤师峡谷"
			mskill[key] = skill
		}

		return mskill
	}

}
