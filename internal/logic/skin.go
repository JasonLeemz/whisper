package logic

import (
	"whisper/internal/model"
	dao "whisper/internal/model/DAO"
	"whisper/pkg/context"
)

func GetHeroSkins(ctx *context.Context, platform int, heroID string) ([]*model.HeroSkin, error) {
	sd := dao.NewHeroSkinDAO()
	ret, err := sd.Find([]string{
		"heroId", "skinId", "heroTitle", "heroName", "name", "description", "emblemsName", "mainImg", "iconImg", "loadingImg", "videoImg", "sourceImg", "isBase", "platform", "version", "fileTime",
	}, map[string]interface{}{
		"heroId": heroID,
	})
	return ret, err
}
