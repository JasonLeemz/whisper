package service

import (
	"encoding/json"
	"fmt"
	"sync"
	"whisper/internal/dto"
	"whisper/pkg/config"
	"whisper/pkg/context"
	"whisper/pkg/errors"
	"whisper/pkg/http"
	"whisper/pkg/log"
)

type LOLM struct {
	platform      int
	ts            int64
	referer       []http.Header
	cookie        []http.Header
	commonHeaders []http.Header
}

func (lol *LOLM) QueryEquipments(ctx *context.Context) (interface{}, error) {
	url := fmt.Sprintf("%s?ts=%d", config.LOLConfig.LolM.Equipment, lol.ts)
	log.Logger.Info(ctx, "url="+url)

	// 发送 GetForm 请求
	equip := dto.LOLMEquipment{}

	body, err := http.GetForm(ctx, url)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &equip)
	return &equip, err
}

func (lol *LOLM) QueryHeroes(ctx *context.Context) (interface{}, error) {
	url := fmt.Sprintf("%s?ts=%d", config.LOLConfig.LolM.Heroes, lol.ts)
	log.Logger.Info(ctx, "url="+url)

	// 发送 GetForm 请求
	heroes := dto.LOLMHeroes{}

	body, err := http.GetForm(ctx, url)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &heroes)
	return &heroes, err
}

func (lol *LOLM) QueryRune(ctx *context.Context) (interface{}, error) {
	url := fmt.Sprintf("%s?ts=%d", config.LOLConfig.LolM.Rune, lol.ts)
	log.Logger.Info(ctx, "url="+url)

	// 发送 GetForm 请求
	r := dto.LOLMRune{}

	body, err := http.GetForm(ctx, url)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &r)
	return &r, err
}

func (lol *LOLM) QuerySkill(ctx *context.Context) (interface{}, error) {
	url := fmt.Sprintf("%s?ts=%d", config.LOLConfig.LolM.Skill, lol.ts)
	log.Logger.Info(ctx, "url="+url)

	// 发送 GetForm 请求
	s := dto.LOLMSkill{}

	body, err := http.GetForm(ctx, url)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &s)
	return &s, err
}

func (lol *LOLM) GetHeroAttribute(ctx *context.Context, heroID string) (interface{}, error) {
	heroAttrUrl := fmt.Sprintf(config.LOLConfig.LolM.Hero, heroID)
	url := fmt.Sprintf("%s?ts=%d", heroAttrUrl, lol.ts)
	//log.Logger.Info(ctx, "url="+url)

	// 发送 GetForm 请求
	r := dto.HeroAttribute{}

	body, err := http.GetForm(ctx, url)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &r)
	return &r, err
}

func (lol *LOLM) QueryRuneType(ctx *context.Context) (interface{}, error) {
	url := fmt.Sprintf("%s", config.LOLConfig.LolM.RuneType)
	log.Logger.Info(ctx, "url="+url)

	// 发送 GetForm 请求
	runeType := dto.LOLMRuneType{}

	body, err := http.GetForm(ctx, url)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &runeType)
	return &runeType, err
}

func (lol *LOLM) QuerySuitEquip(ctx *context.Context, heroID string) (interface{}, error) {
	heroTechUrl := fmt.Sprintf(config.LOLConfig.LolM.HeroSuit, heroID)
	log.Logger.Info(ctx, "heroTechUrl="+heroTechUrl)

	// 发送 GetForm 请求
	heroTech := dto.HeroTech{}
	// -----------------------------
	body, err := http.GetForm(ctx, heroTechUrl, lol.referer...)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &heroTech)
	if err != nil {
		return nil, err
	}

	// -----------------------------
	wg := sync.WaitGroup{}

	syncMap := sync.Map{}
	for _, eqs := range heroTech.Data.AnchorRecommend.List {
		wg.Add(1)
		go func(eqs dto.AnchorRecommendList) {
			defer wg.Done()

			equipTechUrl := fmt.Sprintf(config.LOLConfig.LolM.HeroEquip, eqs.Head.Id)
			log.Logger.Info(ctx, "equipTechUrl="+equipTechUrl)
			body, err = http.GetForm(ctx, equipTechUrl, lol.referer...)
			if err != nil {
				log.Logger.Error(ctx, err)
				return
			}

			et := dto.EquipTech{}
			err = json.Unmarshal(body, &et)
			if err != nil {
				log.Logger.Error(ctx, err)
				return
			}

			syncMap.Store(eqs.Head.Id, &et)
		}(eqs)
	}
	wg.Wait()

	met := make(map[string]*dto.EquipTech)
	syncMap.Range(func(key, value interface{}) bool {
		met[key.(string)] = value.(*dto.EquipTech)
		return true
	})
	heroTech.EquipData = met
	return &heroTech, nil
}

func (lol *LOLM) HeroRankData(ctx *context.Context, heroID string) (interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (lol *LOLM) HeroRankList(ctx *context.Context) (interface{}, error) {
	url := config.LOLConfig.LolM.HeroWinRate
	log.Logger.Info(ctx, "url="+url)

	// 发送 GetForm 请求
	rankList := dto.HeroRankList{}

	body, err := http.GetForm(ctx, url, lol.referer...)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &rankList)
	if err != nil {
		return nil, err
	}

	return &rankList, nil
}

func (lol *LOLM) VersionList(ctx *context.Context) (interface{}, error) {

	versionListUrl := config.LOLConfig.LolM.VersionList

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
		return nil, errors.New(versionList.Msg, versionList.ErrMsg)
	}

	return &versionList, nil
}

func (lol *LOLM) VersionDetail(ctx *context.Context, keys []string) (interface{}, error) {
	// https://mlol.qt.qq.com/go/database/versionlist?zone=lgame
	versionDetailUrl := config.LOLConfig.LolM.VersionDetail

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

func (lol *LOLM) VersionInfo(ctx *context.Context, vKey, id string) (interface{}, error) {
	versionInfoUrl := fmt.Sprintf(config.LOLConfig.LolM.VersionInfo, "lgame_"+vKey)

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
