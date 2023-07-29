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

func QueryRune(ctx *context.Context, platform int) (any, *errors.Error) {

	if platform == common.PlatformForLOL {
		runes, err := service.QueryRuneForLOL(ctx)
		if err != nil {
			log.Logger.Warn(ctx, err)
		}
		reloadRuneForLOL(ctx, runes)
		return runes, nil
	} else if platform == common.PlatformForLOLM {
		runes, err := service.QueryRuneForLOLM(ctx)
		if err != nil {
			log.Logger.Warn(ctx, err)
		}
		reloadRuneForLOLM(ctx, runes)
		return runes, nil
	}

	return nil, errors.New(errors2.New("请指定游戏平台"), errors.ErrNoInvalidInput)
}

func reloadRuneForLOL(ctx *context.Context, r *dto.LOLRune) {
	// 入库更新
	rs := make([]*model.LOLRune, 0, len(r.Rune))
	var err error

	for _, rr := range r.Rune {
		tmp := model.LOLRune{
			Name:      rr.Name,
			Icon:      rr.Icon,
			Key:       rr.Key,
			Tooltip:   rr.Tooltip,
			Shortdesc: rr.ShortDesc,
			Longdesc:  rr.LongDesc,
			SlotLabel: rr.SlotLabel,
			StyleName: rr.StyleName,
			Version:   r.Version,
			FileTime:  r.FileTime,
		}
		rs = append(rs, &tmp)
	}

	// 记录英雄列表信息
	rsd := dao.NewLOLRuneDAO()
	_, err = rsd.Add(rs)
	if err != nil {
		log.Logger.Error(ctx, errors.New(err))
	}
}

func reloadRuneForLOLM(ctx *context.Context, r *dto.LOLMRune) {
	// 入库更新
	rs := make([]*model.LOLMRune, 0, len(r.RuneList))
	var err error

	for _, rr := range r.RuneList {
		tmp := model.LOLMRune{
			RuneId:               rr.RuneId,
			Name:                 rr.Name,
			Description:          rr.Description,
			DetailInfo:           rr.DetailInfo,
			AttrName:             rr.AttrName,
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
	rsDao := dao.NewLOLMRuneDAO()
	_, err = rsDao.Add(rs)
	if err != nil {
		log.Logger.Error(ctx, errors.New(err))
	}
}
