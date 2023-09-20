package logic

import (
	"whisper/internal/model"
	dao "whisper/internal/model/DAO"
	"whisper/pkg/context"
)

func GetHeroSkins(ctx *context.Context, platform int, heroID string) ([]*model.HeroSkin, error) {
	sd := dao.NewHeroSkinDAO()
	ret, err := sd.Find([]string{
		"heroId", "mainImg", "iconImg", "loadingImg", "videoImg", "sourceImg", "videoImg", "platform",
	}, map[string]interface{}{
		"heroId": heroID,
	})
	return ret, err
}
