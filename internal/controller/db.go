package controller

import (
	"whisper/internal/logic"
	"whisper/pkg/context"
	"whisper/pkg/errors"
)

func AliasHeroes(ctx *context.Context) {
	result, err := logic.AliasHeroes(ctx)
	ctx.Reply(result, errors.New(err))
}

func AliasEquip(ctx *context.Context) {
	result, err := logic.AliasEquip(ctx)
	ctx.Reply(result, errors.New(err))
}
