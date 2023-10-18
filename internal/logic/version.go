package logic

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
	"whisper/internal/dto"
	"whisper/internal/logic/common"
	lol "whisper/internal/service/lol"
	"whisper/pkg/context"
	"whisper/pkg/log"
	"whisper/pkg/redis"
)

//Decorator

type GetVersionListFunc func(ctx *context.Context, platform int) ([]dto.VersionListData, error)

// DecorateGetVersionList 装饰结果
func DecorateGetVersionList(fn GetVersionListFunc) GetVersionListFunc {
	return func(ctx *context.Context, platform int) ([]dto.VersionListData, error) {
		list, err := fn(ctx, platform)
		if err != nil {
			return list, err
		}

		for k, _ := range list {
			list[k].PublicDate = "更新时间: " + list[k].PublicDate
		}

		return list, nil
	}
}

func GetVersionList(ctx *context.Context, platform int) ([]dto.VersionListData, error) {
	var queryFromUrl = false
	var vl []dto.VersionListData

	key := fmt.Sprintf(redis.KeyCacheVersionList, platform)
	data := redis.RDB.Get(ctx, key)
	result, err := data.Result()
	if err != nil {
		queryFromUrl = true
		log.Logger.Warn(ctx, err)
	} else {
		err = json.Unmarshal([]byte(result), &vl)
		if err != nil {
			queryFromUrl = true
			log.Logger.Error(ctx, err)
		}
	}

	if queryFromUrl {
		list, err := lol.CreateLOLProduct(platform)().VersionList(ctx)
		versionList := list.(*dto.VersionList)
		if err != nil {
			return nil, err
		}

		vl = versionList.Data

		// cache
		s, _ := json.Marshal(versionList.Data)
		redis.RDB.Set(ctx, key, s, time.Hour*24)
	}

	return vl, err
}

func VersionDetail(ctx *context.Context, platform int, vKey, id string) (map[string]*dto.VersionDetail, error) {

	// 获取该版本下更新的类别
	cates, err := GetUpdateCats(ctx, platform, vKey, id)
	if err != nil {
		return nil, err
	}

	var keys []string
	keyTpl := ""
	if platform == common.PlatformForLOL {
		// lol_20230830_rune_157
		keyTpl = "lol_" + vKey + "_%s_" + id
	} else {
		// lgame_4.3c_hero
		keyTpl = "lgame_" + vKey + "_%s"
	}

	keyCateName := make(map[string]string)
	for cate, cateName := range cates {
		key := fmt.Sprintf(keyTpl, cate)
		keyCateName[key] = cateName
		keys = append(keys, key)
	}

	log.Logger.Info(ctx, keys)

	details, err := lol.CreateLOLProduct(platform)().VersionDetail(ctx, keys)
	if err != nil {
		return nil, err
	}

	resp := make(map[string]*dto.VersionDetail)
	for key, detail := range details.(map[string]*dto.VersionDetail) {
		resp[keyCateName[key]] = detail
	}
	return resp, nil
}

// GetUpdateCats 获取更新类别
func GetUpdateCats(ctx *context.Context, platform int, vKey, id string) (map[string]string, error) {
	info, err := lol.CreateLOLProduct(platform)().VersionInfo(ctx, vKey, id)
	if err != nil {
		return nil, err
	}
	versionInfo := info.(*dto.VersionInfo)
	if versionInfo.Code != 0 {
		return nil, errors.New(versionInfo.Msg)
	}

	cates := make(map[string]string)
	for _, tab := range versionInfo.Data.Tabs {
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
