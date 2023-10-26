package skill

import (
	"whisper/internal/logic/common"
	dao "whisper/internal/model/DAO"
	"whisper/pkg/context"
)

type Inner struct {
	ctx      *context.Context
	platform int
}

func (e *Inner) GetAll(platform int) (interface{}, error) {
	// 获取全部装备
	if platform == common.PlatformForLOL {
		d := dao.NewLOLSkillDAO()
		eVersion, err := d.GetLOLSkillMaxVersion()
		if err != nil {
			return nil, err
		}
		data, err := d.GetLOLSkill(eVersion.Version)
		if err != nil {
			return nil, err
		}
		return data, nil
	} else {
		d := dao.NewLOLMSkillDAO()
		v, err := d.GetLOLMSkillMaxVersion()
		if err != nil {
			return nil, err
		}
		data, err := d.GetLOLMSkill(v.Version)
		if err != nil {
			return nil, err
		}
		return data, nil
	}

}
