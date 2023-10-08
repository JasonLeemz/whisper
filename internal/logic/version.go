package logic

import (
	"encoding/json"
	"errors"
	"fmt"
	redis2 "github.com/redis/go-redis/v9"
	"whisper/internal/dto"
	"whisper/internal/logic/common"
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
func VersionDetail(ctx *context.Context, version string) (map[string]dto.LOLMVersionDetail, error) {
	data := make(map[string]dto.LOLMVersionDetail)
	//var keys []string // cache:version:detail:lolm:4.3c_hero
	keys := map[string]string{
		"hero":   fmt.Sprintf(redis.KeyCacheVersionDetail, version+"_hero"),
		"prop":   fmt.Sprintf(redis.KeyCacheVersionDetail, version+"_prop"),
		"rune":   fmt.Sprintf(redis.KeyCacheVersionDetail, version+"_rune"),
		"system": fmt.Sprintf(redis.KeyCacheVersionDetail, version+"_system"),
	}

	queryFromUrl := false
	// 尝试从redis获取
	for k, rk := range keys {
		d := redis.RDB.Get(ctx, rk)
		result, err := d.Result()
		if err != nil {
			queryFromUrl = true
			log.Logger.Warn(ctx, err)
			break
		} else {
			vd := dto.LOLMVersionDetail{}
			_ = json.Unmarshal([]byte(result), &vd)
			data[k] = vd
		}
	}

	if queryFromUrl {
		detail, err := service.VersionDetail(ctx, common.PlatformForLOLM, version)
		if err != nil {
			log.Logger.Error(ctx, err)
			return nil, err
		}

		for k, rk := range keys {
			// lgame_%s_hero
			key := fmt.Sprintf("lgame_%s_%s", version, k)
			if d, ok := detail[key]; ok {
				s, _ := json.Marshal(d)
				redis.RDB.Set(ctx, rk, s, redis2.KeepTTL)
			}
		}
		return nil, errors.New("try again")
	} else {
		return data, nil
	}
}
