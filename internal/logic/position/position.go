package position

import (
	"errors"
	"whisper/internal/dto"
	"whisper/internal/logic/common"
	"whisper/internal/model"
	dao "whisper/internal/model/DAO"
	lol "whisper/internal/service/lol"
	"whisper/pkg/context"
	"whisper/pkg/log"
	"whisper/pkg/utils"
)

type position struct {
	ctx      *context.Context
	platform int
}

type Func func(ctx *context.Context, platform int) *position

func NewPosition() Func {
	return func(ctx *context.Context, platform int) *position {
		return &position{ctx: ctx, platform: platform}
	}
}

func (pos *position) HeroesPosition() (*dto.HeroRankList, error) {
	list, err := lol.CreateLOLProduct(pos.platform)().HeroRankList(pos.ctx)
	if err != nil {
		log.Logger.Error(pos.ctx, err)
		return nil, err
	}

	rankList := list.(*dto.HeroRankList)

	// 删除旧数据
	cond := map[string]interface{}{
		"platform": common.PlatformForLOLM,
	}

	hp := make([]*model.HeroesPosition, 0)
	// 只取钻石以上分段
	if levData, ok := rankList.Data[common.LevelDiamond]; ok {
		for pos, heroes := range levData {
			posName := common.PositionNameEN[pos]
			for _, data := range heroes {
				hp = append(hp, &model.HeroesPosition{
					HeroId:   data.HeroId,
					Pos:      posName,
					ShowRate: utils.Str2Int(data.AppearRate),
					WinRate:  utils.Str2Int(data.WinRate),
					Platform: common.PlatformForLOLM,
					Version:  data.Dtstatdate,
					FileTime: data.Dtstatdate,
				})
			}
		}
	}

	err = pos.UpdateHeroesPosition(cond, hp)
	if err != nil {
		log.Logger.Error(pos.ctx, err)
		return nil, err
	}
	return rankList, nil
}

func (pos *position) UpdateHeroesPosition(cond map[string]interface{}, posData []*model.HeroesPosition) error {
	hpd := dao.NewHeroesPositionDAO()
	if len(posData) == 0 {
		log.Logger.Warn(pos.ctx, "posData is nil", "cond:", cond)
		return nil
	}

	err := hpd.DeleteAndInsert(cond, posData)
	if err != nil {
		return errors.New("Add HeroesPosition " + err.Error())
	}

	return nil
}
