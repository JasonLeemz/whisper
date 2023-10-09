package controller

import (
	"github.com/spf13/cast"
	"whisper/internal/dto"
	"whisper/internal/logic"
	"whisper/pkg/context"
	"whisper/pkg/errors"
)

func QueryVersion(ctx *context.Context) {
	v := dto.Version{}
	version := logic.GetVersion(ctx)
	if version != nil {
		v = dto.Version{
			LOL: dto.LocalVersionDetail{
				Version:    version[0].Version,
				UpdateTime: version[0].FileTime,
			},
			LOLM: dto.LocalVersionDetail{
				Version:    version[1].Version,
				UpdateTime: version[1].FileTime,
			},
		}
	}

	ctx.Reply(v, nil)
}

type VersionListReq struct {
	Platform int `json:"platform" form:"platform"`
}

type VersionListResp struct {
	Tips string                `json:"tips"`
	Data []dto.PageVersionList `json:"data"`
}

func VersionList(ctx *context.Context) {
	req := &VersionListReq{}
	if err := ctx.Bind(req); err != nil {
		return
	}

	platform := cast.ToInt(req.Platform)
	versionList, err := logic.GetVersionList(ctx, platform)
	pvl := make([]dto.PageVersionList, 0, len(versionList))
	for _, item := range versionList {
		pvl = append(pvl, dto.PageVersionList{
			Id:           item.Id,
			Name:         item.Name,
			Title:        item.Title,
			Vkey:         item.Vkey,
			Introduction: item.Introduction,
			Isnew:        item.Isnew,
			Image:        item.Image,
			PublicDate:   item.PublicDate,
			Platform:     platform,
		})
	}
	data := VersionListResp{
		Tips: "",
		Data: pvl,
	}
	ctx.Reply(data, errors.New(err))
}

type VersionDetailReq struct {
	Platform int    `json:"platform" form:"platform"`
	Version  string `json:"version" form:"version"`
	ID       string `json:"id" form:"id"`
}

func VersionDetail(ctx *context.Context) {
	req := &VersionDetailReq{}
	if err := ctx.Bind(req); err != nil {
		return
	}
	detail, err := logic.VersionDetail(ctx, req.Platform, req.Version, req.ID)
	ctx.Reply(detail, errors.New(err))
}
