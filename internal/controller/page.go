package controller

import (
	"whisper/internal/logic"
	"whisper/internal/logic/equipment"
	"whisper/pkg/context"
	"whisper/pkg/errors"
)

func SearchBox(ctx *context.Context) {
	ctx.Render("index.html", nil)
}

type ReqGetRoadmap struct {
	ID       string   `form:"id" json:"id" binding:"required"`
	Version  string   `form:"version" json:"version" binding:"required"`
	Maps     []string `json:"map,omitempty" form:"map,omitempty" binding:"required"`
	Platform int      `form:"platform" json:"platform" binding:"-"`
}

func GetRoadmap(ctx *context.Context) {
	req := &ReqGetRoadmap{}
	if err := ctx.Bind(req); err != nil {
		return
	}
	roadmap, err := equipment.GetRoadmap(ctx, req.Version, req.Platform, req.ID, req.Maps)

	ctx.Reply(roadmap, errors.New(err))
}

func GetHotKey(ctx *context.Context) {
	keys := logic.GetHotKey(ctx)
	ctx.Reply(keys, nil)
}
