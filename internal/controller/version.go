package controller

import (
	"whisper/internal/dto"
	"whisper/internal/logic"
	"whisper/pkg/context"
	"whisper/pkg/errors"
)

func QueryVersion(ctx *context.Context) {
	v := dto.Version{}
	version := logic.GetVersion(ctx)
	if version != nil {
		v = dto.Version{
			LOL: dto.VersionDetail{
				Version:    version[0].Version,
				UpdateTime: version[0].FileTime,
			},
			LOLM: dto.VersionDetail{
				Version:    version[1].Version,
				UpdateTime: version[1].FileTime,
			},
		}
	}

	ctx.Reply(v, nil)
}

func LOLMVersion(ctx *context.Context) {
	versionList, err := logic.GetLOLMVersionList(ctx)
	ctx.Reply(versionList, errors.New(err))
}
