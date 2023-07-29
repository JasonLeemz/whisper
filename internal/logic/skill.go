package logic

import (
	errors2 "errors"
	"whisper/internal/dto"
	"whisper/internal/logic/common"
	"whisper/internal/model"
	dao "whisper/internal/model/DAO"
	"whisper/pkg/log"

	"whisper/internal/service"
	"whisper/pkg/context"
	"whisper/pkg/errors"
)

func QuerySkill(ctx *context.Context, platform int) (any, *errors.Error) {

	if platform == common.PlatformForLOL {
		skills, err := service.QuerySkillForLOL(ctx)
		if err != nil {
			log.Logger.Warn(ctx, err)
		}
		reloadSkillForLOL(ctx, skills)
		return skills, nil
	} else if platform == common.PlatformForLOLM {
		skills, err := service.QuerySkillForLOLM(ctx)
		if err != nil {
			log.Logger.Warn(ctx, err)
		}
		reloadSkillForLOLM(ctx, skills)
		return skills, nil
	}

	return nil, errors.New(errors2.New("请指定游戏平台"), errors.ErrNoInvalidInput)
}

func reloadSkillForLOL(ctx *context.Context, s *dto.LOLSkill) {
	// 入库更新
	sss := make([]*model.LOLSkill, 0, len(s.SummonerSkill))
	var err error

	for _, ss := range s.SummonerSkill {
		tmp := model.LOLSkill{
			Name:          ss.Name,
			Description:   ss.Description,
			Summonerlevel: ss.SummonerLevel,
			Cooldown:      ss.CoolDown,
			Gamemode:      ss.GameMode,
			Icon:          ss.Icon,
			Version:       s.Version,
			FileTime:      s.FileTime,
		}
		sss = append(sss, &tmp)
	}

	// 记录英雄列表信息
	sd := dao.NewLOLSkillDAO()
	_, err = sd.Add(sss)
	if err != nil {
		log.Logger.Error(ctx, errors.New(err))
	}
}
func reloadSkillForLOLM(ctx *context.Context, s *dto.LOLMSkill) {
	// 入库更新
	ssl := make([]*model.LOLMSkill, 0, len(s.SkillList))
	var err error

	for _, ss := range s.SkillList {
		tmp := model.LOLMSkill{
			SkillId:   ss.SkillId,
			Name:      ss.Name,
			IconPath:  ss.IconPath,
			FuncDesc:  ss.FuncDesc,
			Cd:        ss.Cd,
			Video:     ss.Video,
			Unlocklv:  ss.UnlockLv,
			Mode:      ss.Mode,
			SortOrder: ss.SortOrder,
			Version:   s.Version,
			FileTime:  s.FileTime,
		}
		ssl = append(ssl, &tmp)
	}

	// 记录英雄列表信息
	sd := dao.NewLOLMSkillDAO()
	_, err = sd.Add(ssl)
	if err != nil {
		log.Logger.Error(ctx, errors.New(err))
	}
}
