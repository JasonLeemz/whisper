package equipment

import (
	"go.mongodb.org/mongo-driver/bson"
	"strings"
	"whisper/internal/model"
	dao "whisper/internal/model/DAO"
	"whisper/pkg/context"
	"whisper/pkg/log"
)

func FilterKeyWords(ctx *context.Context, keywords []string, platform int) ([]*model.EquipIntro, error) {
	log.Logger.Info(ctx, keywords)
	// FromMongo
	md := dao.NewMongoEquipmentDAO()

	kw := make([]bson.M, 0, len(keywords))
	for _, words := range keywords {
		in := strings.Split(words, ",")
		kw = append(kw, bson.M{
			"keywords": bson.M{
				"$in": in,
			},
		})
	}
	// 构建查询条件
	filter := bson.M{
		"platform": platform,
		"maps":     "召唤师峡谷",
		"$and":     kw,
	}

	result, err := md.Find(ctx, filter)
	return result, err
}
