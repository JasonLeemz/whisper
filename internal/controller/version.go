package controller

import (
	"whisper/internal/dto"
	"whisper/internal/logic"
	"whisper/internal/logic/common"
	"whisper/pkg/context"
	"whisper/pkg/errors"
)

func QueryVersion(ctx *context.Context) {
	v := dto.Version{}
	version := logic.GetVersion(ctx)
	if version != nil {
		v = dto.Version{
			LOL: dto.VersionDetail{
				Version:    version[0].Version,
				UpdateTime: version[0].FileTime,
			},
			LOLM: dto.VersionDetail{
				Version:    version[1].Version,
				UpdateTime: version[1].FileTime,
			},
		}
	}

	ctx.Reply(v, nil)
}

type LOLMVersionResp struct {
	Tips string                `json:"tips"`
	Data []dto.PageVersionList `json:"data"`
}

func LOLMVersion(ctx *context.Context) {
	versionList, err := logic.GetLOLMVersionList(ctx)
	pvl := make([]dto.PageVersionList, 0, len(versionList))
	for _, item := range versionList {
		pvl = append(pvl, dto.PageVersionList{
			Name:         item.Name,
			Title:        item.Title,
			Vkey:         item.Vkey,
			Introduction: item.Introduction,
			Isnew:        item.Isnew,
			Image:        item.Image,
			PublicDate:   item.PublicDate,
			Platform:     common.PlatformForLOLM,
		})
	}
	data := LOLMVersionResp{
		Tips: "手游版本列表",
		Data: pvl,
	}
	ctx.Reply(data, errors.New(err))
}

func VersionDetail(ctx *context.Context) {
	ctx.Reply(nil, nil)
}
