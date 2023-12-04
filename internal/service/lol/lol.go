package lol

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"sync"
	"whisper/internal/dto"
	"whisper/pkg/config"
	"whisper/pkg/context"
	"whisper/pkg/http"
	"whisper/pkg/log"
)

type LOL struct {
	platform      int
	ts            int64
	yesterday     string
	referer       []http.Header
	cookie        []http.Header
	commonHeaders []http.Header
}

func (lol *LOL) QueryEquipments(ctx *context.Context) (interface{}, error) {
	url := fmt.Sprintf("%s?ts=%d", config.LOLConfig.Lol.Equipment, lol.ts)
	log.Logger.Info(ctx, "url="+url)

	// 发送 GetForm 请求
	equip := dto.LOLEquipment{}

	body, err := http.GetForm(ctx, url)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &equip)
	return &equip, err
}

func (lol *LOL) QueryHeroes(ctx *context.Context) (interface{}, error) {
	url := fmt.Sprintf("%s?ts=%d", config.LOLConfig.Lol.Heroes, lol.ts)
	log.Logger.Info(ctx, "url="+url)

	// 发送 GetForm 请求
	heroes := dto.LOLHeroes{}

	body, err := http.GetForm(ctx, url)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &heroes)
	return &heroes, err
}

func (lol *LOL) QueryRune(ctx *context.Context) (interface{}, error) {
	url := fmt.Sprintf("%s?ts=%d", config.LOLConfig.Lol.Rune, lol.ts)
	log.Logger.Info(ctx, "url="+url)

	// 发送 GetForm 请求
	r := dto.LOLRune{}

	body, err := http.GetForm(ctx, url)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &r)
	return &r, err
}

func (lol *LOL) QuerySkill(ctx *context.Context) (interface{}, error) {
	url := fmt.Sprintf("%s?ts=%d", config.LOLConfig.Lol.Skill, lol.ts)
	log.Logger.Info(ctx, "url="+url)

	// 发送 GetForm 请求
	s := dto.LOLSkill{}

	body, err := http.GetForm(ctx, url)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &s)
	return &s, err
}

func (lol *LOL) GetHeroAttribute(ctx *context.Context, heroID string) (interface{}, error) {
	heroAttrUrl := fmt.Sprintf(config.LOLConfig.Lol.Hero, heroID)
	url := fmt.Sprintf("%s?ts=%d", heroAttrUrl, lol.ts)
	log.Logger.Debug(ctx, "url="+url)

	// 发送 GetForm 请求
	r := dto.HeroAttribute{}

	body, err := http.GetForm(ctx, url)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &r)
	return &r, err
}

func (lol *LOL) QueryRuneType(ctx *context.Context) (interface{}, error) {
	return nil, errors.New("QueryRuneType implement me")
}

func (lol *LOL) QuerySuitEquip(ctx *context.Context, heroID string) (interface{}, error) {
	url := fmt.Sprintf(config.LOLConfig.Lol.SuitEquip, lol.yesterday, heroID)
	log.Logger.Info(ctx, "url="+url)

	// 发送 GetForm 请求
	suitEquip := dto.HeroSuitEquip{}

	body, err := http.GetForm(ctx, url, lol.referer...)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &suitEquip)
	if err != nil {
		return nil, err
	}

	result := dto.JDataDataResult{}
	err = json.Unmarshal([]byte(suitEquip.JData.Data.Result), &result)
	if err != nil {
		return nil, err
	}

	return &result, err
}

func (lol *LOL) HeroRankData(ctx *context.Context, heroID string) (interface{}, error) {
	//jsonpResponse := `var CHAMPION_DETAIL_17={"gameVer":"13.16","date":"2023-08-30 16:15:26"};/*  |xGv00|b214aa8b2b62d14489dce9170b96cdee */`
	champDetailUrl := fmt.Sprintf(config.LOLConfig.Lol.ChampDetail, heroID)
	url := fmt.Sprintf("%s?ts=%d", champDetailUrl, lol.ts)
	log.Logger.Info(ctx, "url="+url)

	// 发送 GetForm 请求
	championFightData := dto.ChampionFightData{}

	body, err := http.GetForm(ctx, url, lol.referer...)
	if err != nil {
		return nil, err
	}

	// 使用正则表达式提取 JSON 数据部分
	re := regexp.MustCompile(`{.*}`)
	match := re.FindString(string(body))

	err = json.Unmarshal([]byte(match), &championFightData)
	if err != nil {
		return nil, err
	}

	return &championFightData, nil
}

func (lol *LOL) HeroRankList(ctx *context.Context) (interface{}, error) {
	return nil, errors.New("implement me")
}

func (lol *LOL) VersionList(ctx *context.Context) (interface{}, error) {
	versionListUrl := config.LOLConfig.Lol.VersionList

	log.Logger.Info(ctx, "versionListUrl="+versionListUrl)

	// 发送 GetForm 请求
	versionList := dto.VersionList{}
	body, err := http.GetForm(ctx, versionListUrl, lol.cookie...)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &versionList)
	if err != nil {
		return nil, err
	}

	if versionList.Result != 0 {
		return nil, errors.New(versionList.Msg + "|" + versionList.ErrMsg)
	}

	return &versionList, nil
}

func (lol *LOL) VersionDetail(ctx *context.Context, keys []string) (interface{}, error) {
	// https://mlol.qt.qq.com/go/database/versiondetail?key=%s
	versionDetailUrl := config.LOLConfig.Lol.VersionDetail

	wg := sync.WaitGroup{}
	syncMap := sync.Map{}
	for _, k := range keys {
		wg.Add(1)
		go func(k string) {
			defer wg.Done()

			detailUrl := fmt.Sprintf(versionDetailUrl, k)
			log.Logger.Info(ctx, "detailUrl="+detailUrl)
			body, err := http.GetForm(ctx, detailUrl, lol.commonHeaders...)
			if err != nil {
				log.Logger.Error(ctx, err)
				return
			}

			detail := dto.VersionDetail{}
			err = json.Unmarshal(body, &detail)
			if err != nil {
				log.Logger.Error(ctx, err)
				return
			}

			syncMap.Store(k, &detail)
		}(k)
	}
	wg.Wait()

	vd := make(map[string]*dto.VersionDetail)
	syncMap.Range(func(key, value interface{}) bool {
		vd[key.(string)] = value.(*dto.VersionDetail)
		return true
	})

	return vd, nil
}

func (lol *LOL) VersionInfo(ctx *context.Context, vKey, id string) (interface{}, error) {
	// https://mlol.qt.qq.com/go/database/versioninfo?key=%s # lol_20170111_10
	versionInfoUrl := fmt.Sprintf(config.LOLConfig.Lol.VersionInfo, "lol_"+vKey+"_"+id)
	log.Logger.Info(ctx, "versionInfoUrl="+versionInfoUrl)

	// 发送 GetForm 请求
	versionInfo := dto.VersionInfo{}
	body, err := http.GetForm(ctx, versionInfoUrl, lol.commonHeaders...)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &versionInfo)
	if err != nil {
		return nil, err
	}

	return &versionInfo, nil
}
