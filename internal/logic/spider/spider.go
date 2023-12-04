package spider

import "whisper/pkg/context"

type Spider struct {
	ctx *context.Context
}

func NewSpider(ctx *context.Context) *Spider {
	return &Spider{ctx: ctx}
}
