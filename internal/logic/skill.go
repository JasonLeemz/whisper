package logic

import (
	errors2 "errors"
	"fmt"
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

	// 判断库中是否存在最新版本，如果存在就不更新
	skillDAO := dao.NewLOLSkillDAO()
	result, err := skillDAO.Find([]string{
		"max(ctime) as ctime",
		"fileTime",
		"version",
	}, nil)
	if err != nil {
		log.Logger.Error(ctx, errors.New(err))
		return
	}

	if len(result) > 0 {
		log.Logger.Info(ctx,
			fmt.Sprintf("DB Version[%s] fileTime[%s],Data Version:[%s] fileTime[%s]", result[0].Version, result[0].FileTime, s.Version, s.FileTime),
		)
		x, err := common.CompareTime(result[0].FileTime, s.FileTime)
		if err != nil {
			log.Logger.Error(ctx, errors.New(err))
			return
		}
		if x != "<" {
			// 如果原始数据版本和当前获取数据的版本相等，就不更新数据库
			log.Logger.Info(ctx, "原始数据版本和当前获取数据的版本相等,不更新")
			return
		}
	}

	// 入库更新
	sss := make([]*model.LOLSkill, 0, len(s.SummonerSkill))

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
	_, err = skillDAO.Add(sss)
	if err != nil {
		log.Logger.Error(ctx, errors.New(err))
	}
}
func reloadSkillForLOLM(ctx *context.Context, s *dto.LOLMSkill) {
	// 判断库中是否存在最新版本，如果存在就不更新
	skillDAO := dao.NewLOLMSkillDAO()
	result, err := skillDAO.Find([]string{
		"max(ctime) as ctime",
		"fileTime",
		"version",
	}, nil)
	if err != nil {
		log.Logger.Error(ctx, errors.New(err))
		return
	}

	if len(result) > 0 {
		log.Logger.Info(ctx,
			fmt.Sprintf("DB Version[%s] fileTime[%s],Data Version:[%s] fileTime[%s]", result[0].Version, result[0].FileTime, s.Version, s.FileTime),
		)
		x, err := common.CompareTime(result[0].FileTime, s.FileTime)
		if err != nil {
			log.Logger.Error(ctx, errors.New(err))
			return
		}
		if x != "<" {
			// 如果原始数据版本和当前获取数据的版本相等，就不更新数据库
			log.Logger.Info(ctx, "原始数据版本和当前获取数据的版本相等,不更新")
			return
		}
	}

	// 入库更新
	ssl := make([]*model.LOLMSkill, 0, len(s.SkillList))

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
	_, err = skillDAO.Add(ssl)
	if err != nil {
		log.Logger.Error(ctx, errors.New(err))
	}
}
