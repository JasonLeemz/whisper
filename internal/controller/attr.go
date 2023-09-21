package controller

import (
	"whisper/internal/logic"
	"whisper/pkg/context"
	"whisper/pkg/errors"
)

func AttrData2Redis(ctx *context.Context) {
	err := logic.AttrData2Redis(ctx)
	ctx.Reply(nil, errors.New(err))
}
