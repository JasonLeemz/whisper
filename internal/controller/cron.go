package controller

import (
	"whisper/internal/logic"
	"whisper/pkg/context"
	"whisper/pkg/errors"
)

func Cron(ctx *context.Context) {

	logic.Cron(ctx)
	ctx.Reply(nil, errors.New(nil))
}
