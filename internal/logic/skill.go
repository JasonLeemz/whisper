package logic

import (
	errors2 "errors"
	"fmt"
	"time"
	"whisper/internal/dto"
	"whisper/internal/logic/common"
	"whisper/internal/model"
	dao "whisper/internal/model/DAO"
	"whisper/pkg/log"
	"whisper/pkg/pinyin"

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
	result, err := skillDAO.GetLOLSkillMaxVersion()
	if err != nil {
		log.Logger.Error(ctx, errors.New(err))
		return
	}

	if result != nil {
		log.Logger.Info(ctx,
			fmt.Sprintf("DB Version[%s] fileTime[%s],Data Version:[%s] fileTime[%s]", result.Version, result.FileTime, s.Version, s.FileTime),
		)
		x, err := common.CompareTime(result.FileTime, s.FileTime)
		if err != nil {
			log.Logger.Error(ctx, errors.New(err))
			return
		}
		if x != "<" {
			// 如果原始数据版本和当前获取数据的版本相等，就不更新数据库
			log.Logger.Info(ctx, "原始数据版本和当前获取数据的版本相等,不更新")
			return
		}

		// 有可能日期更新了，但是版本号没变
		if result.Version == s.Version {
			// 将db中该版本的数据status置为1
			cond := map[string]interface{}{
				"version": result.Version,
				"status":  0,
			}
			data := model.LOLSkill{
				Status: 1,
			}
			up, err := skillDAO.Update(&data, cond)
			if err != nil {
				log.Logger.Error(ctx, err)
				return
			}
			log.Logger.Info(ctx, "当前版本数据不是最新,已经软删除,生效行数:", up)
		}
	}

	log.Logger.Info(ctx, "running record LOL skill data...")
	startT := time.Now()
	// 入库更新
	sss := make([]*model.LOLSkill, 0, len(s.SummonerSkill))

	for _, ss := range s.SummonerSkill {
		if ss.Name == "" {
			continue
		}

		namePY, nameF := pinyin.Trans(ss.Name)
		searchKey := namePY + "," + nameF
		tmp := model.LOLSkill{
			Name:          ss.Name,
			Description:   ss.Description,
			Keywords:      searchKey,
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

	log.Logger.Info(ctx, fmt.Sprintf("finish record LOL skill data. Since:%fs", time.Since(startT).Seconds()))
}
func reloadSkillForLOLM(ctx *context.Context, s *dto.LOLMSkill) {
	// 判断库中是否存在最新版本，如果存在就不更新
	skillDAO := dao.NewLOLMSkillDAO()
	result, err := skillDAO.GetLOLMSkillMaxVersion()
	if err != nil {
		log.Logger.Error(ctx, errors.New(err))
		return
	}

	if result != nil {
		log.Logger.Info(ctx,
			fmt.Sprintf("DB Version[%s] fileTime[%s],Data Version:[%s] fileTime[%s]", result.Version, result.FileTime, s.Version, s.FileTime),
		)
		x, err := common.CompareTime(result.FileTime, s.FileTime)
		if err != nil {
			log.Logger.Error(ctx, errors.New(err))
			return
		}
		if x != "<" {
			// 如果原始数据版本和当前获取数据的版本相等，就不更新数据库
			log.Logger.Info(ctx, "原始数据版本和当前获取数据的版本相等,不更新")
			return
		}

		// 有可能日期更新了，但是版本号没变
		if result.Version == s.Version {
			// 将db中该版本的数据status置为1
			cond := map[string]interface{}{
				"version": result.Version,
				"status":  0,
			}
			data := model.LOLMSkill{
				Status: 1,
			}
			up, err := skillDAO.Update(&data, cond)
			if err != nil {
				log.Logger.Error(ctx, err)
				return
			}
			log.Logger.Info(ctx, "当前版本数据不是最新,已经软删除,生效行数:", up)
		}
	}

	log.Logger.Info(ctx, "running record LOLM skill data...")
	startT := time.Now()
	// 入库更新
	ssl := make([]*model.LOLMSkill, 0, len(s.SkillList))

	for _, ss := range s.SkillList {
		namePY, nameF := pinyin.Trans(ss.Name)
		searchKey := namePY + "," + nameF

		tmp := model.LOLMSkill{
			SkillId:   ss.SkillId,
			Name:      ss.Name,
			IconPath:  ss.IconPath,
			FuncDesc:  ss.FuncDesc,
			Keywords:  searchKey,
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

	log.Logger.Info(ctx, fmt.Sprintf("finish record LOLM skill data. Since:%fs", time.Since(startT).Seconds()))
}
