package controller

import (
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

	ctx.Reply(map[string]interface{}{
		"lol_version":  v1,
		"lolm_version": v2,
	}, nil)
}

func QueryEquipTypes(ctx *context.Context) {

	types, _ := logic.GetEquipTypes(ctx)

	ctx.Reply(map[string]interface{}{
		"types": types,
	}, nil)
}

type ReqGetRoadmap struct {
	ID       string `form:"id" json:"id" binding:"required"`
	Version  string `form:"version" json:"version" binding:"required"`
	Platform int    `form:"platform" json:"platform" binding:"-"`
}

func GetRoadmap(ctx *context.Context) {
	req := &ReqGetRoadmap{}
	if err := ctx.Bind(req); err != nil {
		return
	}
	roadmap, err := logic.GetRoadmap(ctx, req.Version, req.Platform, req.ID)

	ctx.Reply(roadmap, errors.New(err))
}
