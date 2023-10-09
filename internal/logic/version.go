package logic

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
	"whisper/internal/dto"
	"whisper/internal/logic/common"
	"whisper/internal/service"
	"whisper/pkg/context"
	"whisper/pkg/log"
	"whisper/pkg/redis"
)

func GetVersionList(ctx *context.Context, platform int) ([]dto.VersionListData, error) {
	var queryFromUrl = false
	var vl []dto.VersionListData

	key := fmt.Sprintf(redis.KeyCacheVersionList, platform)
	data := redis.RDB.Get(ctx, key)
	result, err := data.Result()
	if err != nil {
		queryFromUrl = true
		log.Logger.Error(ctx, err)
	} else {
		err = json.Unmarshal([]byte(result), &vl)
		if err != nil {
			queryFromUrl = true
			log.Logger.Error(ctx, err)
		}
	}

	if queryFromUrl {
		versionList, err2 := service.VersionList(ctx, platform)
		if err2 != nil {
			return nil, err2
		}

		vl = versionList.Data

		// cache
		s, _ := json.Marshal(versionList.Data)
		redis.RDB.Set(ctx, key, s, time.Hour*24)
	}

	return vl, err
}

func VersionDetail(ctx *context.Context, platform int, vkey, id string) (map[string]*dto.VersionDetail, error) {

	// 获取该版本下更新的类别
	cates, err := GetUpdateCates(ctx, platform, vkey, id)
	if err != nil {
		return nil, err
	}

	var keys []string
	keyTpl := ""
	if platform == common.PlatformForLOL {
		// lol_20230830_rune_157
		keyTpl = "lol_" + vkey + "_%s_" + id
	} else {
		// lgame_4.3c_hero
		keyTpl = "lgame_" + vkey + "_%s"
	}

	keyCateName := make(map[string]string)
	for cate, cateName := range cates {
		key := fmt.Sprintf(keyTpl, cate)
		keyCateName[key] = cateName
		keys = append(keys, key)
	}

	log.Logger.Info(ctx, keys)

	details, err := service.VersionDetail(ctx, platform, keys)
	if err != nil {
		return nil, err
	}

	resp := make(map[string]*dto.VersionDetail)
	for key, detail := range details {
		resp[keyCateName[key]] = detail
	}
	return resp, nil
}

// GetUpdateCates 获取更新类别
func GetUpdateCates(ctx *context.Context, platform int, vkey, id string) (map[string]string, error) {
	info, err := service.VersionInfo(ctx, platform, vkey, id)
	if err != nil {
		return nil, err
	}
	if info.Code != 0 {
		return nil, errors.New(info.Msg)
	}

	cates := make(map[string]string)
	for _, tab := range info.Data.Tabs {
		/*
			"tabs": [
			      {
			        "title": "英雄",
			        "schemeUrl": "qtpage://version/detail/list?key=lgame_4.3c_hero\u0026textColor=%232896F5",
			        "is_default_tab": 1,
			        "key": "hero"
			      }
			    ]
		*/
		// {"hero":"英雄"}
		cates[tab.Key] = tab.Title
	}

	return cates, nil
}
