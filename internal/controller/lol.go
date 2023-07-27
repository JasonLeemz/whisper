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
