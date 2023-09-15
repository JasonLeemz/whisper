package controller

import (
	"whisper/internal/dto"
	"whisper/internal/logic"
	"whisper/pkg/context"
	"whisper/pkg/errors"
)

func SearchBox(ctx *context.Context) {
	ctx.Render("index.html", nil)
}

func QueryVersion(ctx *context.Context) {
	// 获取端游版本
	v1 := logic.GetCurrentLOLVersion(ctx)
	// 获取手游版本
	v2 := logic.GetCurrentLOLMVersion(ctx)

	v := dto.Version{
		LOL: dto.VersionDetail{
			Version:    v1.Version,
			UpdateTime: v1.FileTime,
		},
		LOLM: dto.VersionDetail{
			Version:    v2.Version,
			UpdateTime: v2.FileTime,
		},
	}
	ctx.Reply(v, nil)
}

func QueryEquipTypes(ctx *context.Context) {

	types, _ := logic.GetEquipTypes(ctx)

	ctx.Reply(map[string]interface{}{
		"types": types,
	}, nil)
}

type ReqGetRoadmap struct {
	ID       string   `form:"id" json:"id" binding:"required"`
	Version  string   `form:"version" json:"version" binding:"required"`
	Maps     []string `json:"map,omitempty" form:"map,omitempty"`
	Platform int      `form:"platform" json:"platform" binding:"-"`
}

func GetRoadmap(ctx *context.Context) {
	req := &ReqGetRoadmap{}
	if err := ctx.Bind(req); err != nil {
		return
	}
	roadmap, err := logic.GetRoadmap(ctx, req.Version, req.Platform, req.ID, req.Maps)

	ctx.Reply(roadmap, errors.New(err))
}

func GetHotKey(ctx *context.Context) {
	keys := logic.GetHotKey(ctx)
	ctx.Reply(keys, nil)
}
