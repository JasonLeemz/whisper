package controller

import (
	"whisper/internal/logic"
	"whisper/pkg/errors"

	"whisper/pkg/context"
)

type ReqEquipment struct {
	Platform int `form:"platform" json:"platform" binding:"-"`
	//Platform int `form:"platform" json:"platform" binding:"required"`
}

func Equipment(ctx *context.Context) {

	req := &ReqEquipment{}
	if err := ctx.Bind(req); err != nil {
		return
	}

	equip, err := logic.QueryEquipments(ctx, req.Platform)

	ctx.Reply(equip, errors.New(err))
}

type ReqHeroes struct {
	Platform int `form:"platform" json:"platform" binding:"-"`
}

func Heroes(ctx *context.Context) {

	req := &ReqHeroes{}
	if err := ctx.Bind(req); err != nil {
		return
	}

	equip, err := logic.QueryHeroes(ctx, req.Platform)

	ctx.Reply(equip, errors.New(err))
}

type ReqRune struct {
	Platform int `form:"platform" json:"platform" binding:"-"`
}

func Rune(ctx *context.Context) {

	req := &ReqRune{}
	if err := ctx.Bind(req); err != nil {
		return
	}

	equip, err := logic.QueryRune(ctx, req.Platform)

	ctx.Reply(equip, errors.New(err))
}

type ReqSkill struct {
	Platform int `form:"platform" json:"platform" binding:"-"`
}

func Skill(ctx *context.Context) {

	req := &ReqSkill{}
	if err := ctx.Bind(req); err != nil {
		return
	}

	equip, err := logic.QuerySkill(ctx, req.Platform)

	ctx.Reply(equip, errors.New(err))
}
