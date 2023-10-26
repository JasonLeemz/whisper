package rune

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
		d := dao.NewLOLRuneDAO()
		eVersion, err := d.GetLOLRuneMaxVersion()
		if err != nil {
			return nil, err
		}
		data, err := d.GetLOLRune(eVersion.Version)
		if err != nil {
			return nil, err
		}
		return data, nil
	} else {
		d := dao.NewLOLMRuneDAO()
		v, err := d.GetLOLMRuneMaxVersion()
		if err != nil {
			return nil, err
		}
		data, err := d.GetLOLMRune(v.Version)
		if err != nil {
			return nil, err
		}
		return data, nil
	}

}
