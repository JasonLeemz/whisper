package controller

import (
	"whisper/internal/logic"
	"whisper/pkg/context"
)

func SearchBox(ctx *context.Context) {
	// 获取端游版本
	v1 := logic.GetCurrentLOLVersion(ctx)
	// 获取手游版本
	v2 := logic.GetCurrentLOLMVersion(ctx)
	ctx.Render("index.html", map[string]interface{}{
		"lol_version":  v1,
		"lolm_version": v2,
	})
}
