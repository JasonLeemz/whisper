package logic

import (
	"encoding/json"
	"fmt"
	redis2 "github.com/redis/go-redis/v9"
	"whisper/internal/dto"
	"whisper/internal/service"
	"whisper/pkg/context"
	"whisper/pkg/log"
	"whisper/pkg/redis"
)

func GetLOLMVersionList(ctx *context.Context) ([]dto.LOLMVersionListData, error) {
	var queryFromUrl = false
	var vl []dto.LOLMVersionListData
	vd := make(map[string]*dto.LOLMVersionDetail)

	data := redis.RDB.Get(ctx, redis.KeyCacheVersionList)
	result, err := data.Result()
	if err != nil {
		queryFromUrl = true
		log.Logger.Error(ctx, err)
	}
	err = json.Unmarshal([]byte(result), &vl)
	if err != nil {
		queryFromUrl = true
		log.Logger.Error(ctx, err)
	}

	if queryFromUrl {
		versionList, versionDetail, err2 := service.LOLMVersionList(ctx)
		if err2 != nil {
			return nil, err2
		}

		vd = versionDetail
		vl = versionList.Data

		// cache
		// version list
		s, _ := json.Marshal(versionList.Data)
		redis.RDB.Set(ctx, redis.KeyCacheVersionList, s, redis2.KeepTTL)

		// version detail
		for k, detail := range vd {
			key := fmt.Sprintf(redis.KeyCacheVersionDetail, k)
			v, _ := json.Marshal(detail)
			redis.RDB.Set(ctx, key, v, redis2.KeepTTL)
		}
	}

	return vl, err
}
