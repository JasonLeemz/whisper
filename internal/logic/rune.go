package logic

import (
	"fmt"
	"time"
	"whisper/internal/dto"
	"whisper/internal/logic/common"
	"whisper/internal/model"
	dao "whisper/internal/model/DAO"
	"whisper/pkg/log"
	"whisper/pkg/pinyin"

	lol "whisper/internal/service/lol"
	"whisper/pkg/context"
	"whisper/pkg/errors"
)

func QueryRune(ctx *context.Context, platform int) (any, error) {
	runes, err := lol.CreateLOLProduct(platform)().QueryRune(ctx)
	if err != nil {
		log.Logger.Error(ctx, err)
		return nil, err
	}

	if platform == common.PlatformForLOL {
		reloadRuneForLOL(ctx, runes.(*dto.LOLRune))
	} else if platform == common.PlatformForLOLM {
		reloadRuneForLOLM(ctx, runes.(*dto.LOLMRune))

	}
	return runes, nil
}

func reloadRuneForLOL(ctx *context.Context, r *dto.LOLRune) {
	// 判断库中是否存在最新版本，如果存在就不更新
	runeDAO := dao.NewLOLRuneDAO()
	result, err := runeDAO.GetLOLRuneMaxVersion()
	if err != nil {
		log.Logger.Error(ctx, errors.New(err))
		return
	}

	if result != nil {
		log.Logger.Info(ctx,
			fmt.Sprintf("DB Version[%s] fileTime[%s],Data Version:[%s] fileTime[%s]", result.Version, result.FileTime, r.Version, r.FileTime),
		)
		x, err := common.CompareTime(result.FileTime, r.FileTime)
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
		if result.Version == r.Version {
			// 将db中该版本的数据status置为1
			cond := map[string]interface{}{
				"version": result.Version,
				"status":  0,
			}
			data := model.LOLRune{
				Status: 1,
			}
			up, err := runeDAO.Update(&data, cond)
			if err != nil {
				log.Logger.Error(ctx, err)
				return
			}
			log.Logger.Info(ctx, "当前版本数据不是最新,已经软删除,生效行数:", up)
		}
	}

	log.Logger.Info(ctx, "running record LOL rune data...")
	startT := time.Now()
	// 入库更新
	rs := make([]*model.LOLRune, 0, len(r.Rune))

	for runeID, rr := range r.Rune {
		namePY, nameF := pinyin.Trans(rr.Name)
		searchKey := rr.Key + "," + rr.StyleName + "," + namePY + "," + nameF
		tmp := model.LOLRune{
			RuneID:    runeID,
			Name:      rr.Name,
			Icon:      rr.Icon,
			Key:       rr.Key,
			Tooltip:   rr.Tooltip,
			Shortdesc: rr.ShortDesc,
			Longdesc:  rr.LongDesc,
			Keywords:  searchKey,
			SlotLabel: rr.SlotLabel,
			StyleName: rr.StyleName,
			Version:   r.Version,
			FileTime:  r.FileTime,
		}
		rs = append(rs, &tmp)
	}

	// 记录英雄列表信息
	_, err = runeDAO.Add(rs)
	if err != nil {
		log.Logger.Error(ctx, errors.New(err))
	}

	log.Logger.Info(ctx, fmt.Sprintf("finish record LOL rune data. Since:%fs", time.Since(startT).Seconds()))
}

func reloadRuneForLOLM(ctx *context.Context, r *dto.LOLMRune) {
	// 判断库中是否存在最新版本，如果存在就不更新
	runeDAO := dao.NewLOLMRuneDAO()
	result, err := runeDAO.GetLOLMRuneMaxVersion()
	if err != nil {
		log.Logger.Error(ctx, errors.New(err))
		return
	}

	if result != nil {
		log.Logger.Info(ctx,
			fmt.Sprintf("DB Version[%s] fileTime[%s],Data Version:[%s] fileTime[%s]", result.Version, result.FileTime, r.Version, r.FileTime),
		)
		x, err := common.CompareTime(result.FileTime, r.FileTime)
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
		if result.Version == r.Version {
			// 将db中该版本的数据status置为1
			cond := map[string]interface{}{
				"version": result.Version,
				"status":  0,
			}
			data := model.LOLMRune{
				Status: 1,
			}
			up, err := runeDAO.Update(&data, cond)
			if err != nil {
				log.Logger.Error(ctx, err)
				return
			}
			log.Logger.Info(ctx, "当前版本数据不是最新,已经软删除,生效行数:", up)
		}
	}

	log.Logger.Info(ctx, "running record LOLM rune data...")
	startT := time.Now()
	// 入库更新
	rs := make([]*model.LOLMRune, 0, len(r.RuneList))

	for _, rr := range r.RuneList {
		namePY, nameF := pinyin.Trans(rr.Name)
		searchKey := rr.AttrName + "," + rr.Type + "," + namePY + "," + nameF

		tmp := model.LOLMRune{
			RuneId:               rr.RuneId,
			Name:                 rr.Name,
			Description:          rr.Description,
			DetailInfo:           rr.DetailInfo,
			AttrName:             rr.AttrName,
			Keywords:             searchKey,
			Type:                 rr.Type,
			IconPath:             rr.IconPath,
			SortOrder:            rr.SortOrder,
			UnlockLv:             rr.UnlockLv,
			PrimarySlotIndex:     rr.PrimarySlotIndex,
			PrimarySlotSortOrder: rr.PrimarySlotSortOrder,
			Version:              r.Version,
			FileTime:             r.FileTime,
		}
		rs = append(rs, &tmp)
	}

	// 记录英雄列表信息
	_, err = runeDAO.Add(rs)
	if err != nil {
		log.Logger.Error(ctx, errors.New(err))
	}

	log.Logger.Info(ctx, fmt.Sprintf("finish record LOLM rune data. Since:%fs", time.Since(startT).Seconds()))
}

func QueryRuneType(ctx *context.Context, platform int) (any, error) {
	runeType, err := lol.CreateLOLProduct(platform)().QueryRuneType(ctx)
	if err != nil {
		log.Logger.Error(ctx, err)
		return nil, err
	}
	if platform == common.PlatformForLOLM {
		reloadRuneTypeForLOLM(ctx, runeType.(*dto.LOLMRuneType))
	}

	return runeType, nil
}

func reloadRuneTypeForLOLM(ctx *context.Context, rt *dto.LOLMRuneType) {
	// 判断库中是否存在最新版本，如果存在就不更新
	rtDAO := dao.NewRuneTypeDAO()
	_, err := rtDAO.DeleteAll(map[string]interface{}{
		"platform": common.PlatformForLOLM,
	})
	if err != nil {
		log.Logger.Error(ctx, errors.New(err))
		return
	}

	// 入库更新
	rs := make([]*model.RuneType, 0, len(rt.RuneTypes))

	for _, t := range rt.RuneTypes {
		tmp := model.RuneType{
			Name:     t.Name,
			SubType:  t.SubType,
			Type:     t.Type,
			Platform: common.PlatformForLOLM,
		}
		rs = append(rs, &tmp)
	}

	_, err = rtDAO.Add(rs)
	if err != nil {
		log.Logger.Error(ctx, errors.New(err))
	}
}
