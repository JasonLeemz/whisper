package logic

import (
	context2 "context"
	"encoding/json"
	"fmt"
	"github.com/spf13/cast"
	"html"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"whisper/internal/model/es"

	"whisper/internal/dto"
	"whisper/internal/logic/common"
	"whisper/internal/model"
	dao "whisper/internal/model/DAO"
	"whisper/pkg/config"
	"whisper/pkg/context"
	"whisper/pkg/log"
)

type SearchParams struct {
	KeyWords string   `json:"key_words"`
	Platform string   `json:"platform,omitempty"`
	Category string   `json:"category,omitempty"`
	Version  string   `json:"version,omitempty"`
	Way      []string `json:"way,omitempty"`
	Map      []string `json:"map,omitempty"`
}

func EsSearch(ctx *context.Context, p *SearchParams) (*dto.EsResultHits, error) {
	esBuilder := new(es.Instance).Builder(p.Category).
		// 按名字/介绍
		SetMultiMatchQuery(&es.MultiMatchQuery{
			Text:   p.KeyWords,
			Fields: p.Way,
		}).
		// 端游or手游
		SetTermQuery([]*es.TermQuery{
			&es.TermQuery{
				Name:  "platform",
				Value: p.Platform,
			},
		}).
		SetFieldSort(&es.FieldSort{
			Field:     "_score",
			Direction: "desc",
		})
	indexName := p.Category
	// 按地图
	if indexName == new(model.ESEquipment).GetIndexName() {
		var maps []interface{}
		for _, m := range p.Map {
			maps = append(maps, interface{}(m))
		}
		esBuilder = esBuilder.SetTermsQuery(&es.TermsQuery{
			Name:   "maps",
			Values: maps,
		})
	}
	esInstance, err := esBuilder.Build()
	if err != nil {
		return nil, err
	}

	resp, err := esInstance.Query(ctx)
	if err != nil {
		return nil, err
	}

	switch indexName {
	case dao.ESIndexEquipment:
		for i, hit := range resp.Hits {
			sourceStr, _ := json.Marshal(hit.TmpSource)
			hitData := model.ESEquipment{}
			err = json.Unmarshal(sourceStr, &hitData)
			if err != nil {
				return nil, err
			}
			resp.Hits[i].Source.ID = hitData.ID
			resp.Hits[i].Source.Name = hitData.Name
			resp.Hits[i].Source.IconPath = hitData.IconPath
			resp.Hits[i].Source.Description = hitData.Description
			resp.Hits[i].Source.Plaintext = hitData.Plaintext
			resp.Hits[i].Source.Version = hitData.Version
			resp.Hits[i].Source.ItemId = hitData.ItemId
			resp.Hits[i].Source.Platform = hitData.Platform
			resp.Hits[i].Source.Maps = hitData.Maps

			if hitData.Plaintext != "" && !strings.EqualFold(hitData.Plaintext, hitData.Description) {
				resp.Hits[i].Source.Tags = append(resp.Hits[i].Source.Tags, fmt.Sprintf("%s", hitData.Plaintext))
			}
			//resp.Hits[i].Source.Tags = append(resp.Hits[i].Source.Tags, fmt.Sprintf("Price:%s", hitData.Price))
			//resp.Hits[i].Source.Tags = append(resp.Hits[i].Source.Tags, fmt.Sprintf("Sell:%s", hitData.Sell))
			resp.Hits[i].Source.Tags = append(resp.Hits[i].Source.Tags, fmt.Sprintf("价格:%s", hitData.Total))
			resp.Hits[i].Source.Tags = append(resp.Hits[i].Source.Tags, fmt.Sprintf("Version:%s", hitData.Version))
			if hitData.Platform != cast.ToString(common.PlatformForLOLM) {
				resp.Hits[i].Source.Tags = append(resp.Hits[i].Source.Tags, fmt.Sprintf("%s", hitData.Maps))
			}
		}
	case dao.ESIndexHeroes:
		for i, hit := range resp.Hits {
			sourceStr, _ := json.Marshal(hit.TmpSource)
			hitData := model.ESHeroes{}
			err = json.Unmarshal(sourceStr, &hitData)
			if err != nil {
				return nil, err
			}
			resp.Hits[i].Source.ID = hitData.ID
			resp.Hits[i].Source.Name = hitData.Name
			resp.Hits[i].Source.IconPath = hitData.IconPath
			resp.Hits[i].Source.MainImg = hitData.MainImg
			resp.Hits[i].Source.Description = hitData.Description
			resp.Hits[i].Source.Plaintext = hitData.Plaintext
			resp.Hits[i].Source.Version = hitData.Version
			resp.Hits[i].Source.Platform = hitData.Platform

			resp.Hits[i].Source.Tags = append(resp.Hits[i].Source.Tags, strings.Split(hitData.Roles, ",")...)
			resp.Hits[i].Source.Tags = append(resp.Hits[i].Source.Tags, fmt.Sprintf("Version:%s", hitData.Version))
		}
	case dao.ESIndexRune:
		for i, hit := range resp.Hits {
			sourceStr, _ := json.Marshal(hit.TmpSource)
			hitData := model.ESRune{}
			err = json.Unmarshal(sourceStr, &hitData)
			if err != nil {
				return nil, err
			}
			resp.Hits[i].Source.ID = hitData.ID
			resp.Hits[i].Source.Name = hitData.Name + "(" + hitData.StyleName + ")"
			resp.Hits[i].Source.IconPath = hitData.IconPath
			resp.Hits[i].Source.Description = strings.Replace(hitData.Description, "<hr>", "", -1)
			resp.Hits[i].Source.Plaintext = hitData.Plaintext
			resp.Hits[i].Source.Version = hitData.Version
			resp.Hits[i].Source.Platform = hitData.Platform

			//resp.Hits[i].Source.Tags = append(resp.Hits[i].Source.Tags, fmt.Sprintf("Type:%s", hitData.StyleName))
			if hitData.Tooltip != "" && len(hitData.Tooltip) <= 10 {
				tooltip := html.UnescapeString(hitData.Tooltip)
				//tooltip = strings.Replace(hitData.Tooltip, "&lt;br&gt;", "", -1)
				//tooltip = strings.Replace(tooltip, "&lt;hr&gt;", "", -1)
				tooltip = strings.Replace(tooltip, "<br>", "", -1)
				tooltip = strings.Replace(tooltip, "<hr>", "", -1)
				resp.Hits[i].Source.Tags = append(resp.Hits[i].Source.Tags, tooltip)
			}
			resp.Hits[i].Source.Tags = append(resp.Hits[i].Source.Tags, fmt.Sprintf("Version:%s", hitData.Version))
		}
	case dao.ESIndexSkill:
		for i, hit := range resp.Hits {
			sourceStr, _ := json.Marshal(hit.TmpSource)
			hitData := model.ESSkill{}
			err = json.Unmarshal(sourceStr, &hitData)
			if err != nil {
				return nil, err
			}
			resp.Hits[i].Source.ID = hitData.ID
			resp.Hits[i].Source.Name = hitData.Name
			resp.Hits[i].Source.IconPath = hitData.IconPath
			resp.Hits[i].Source.Description = hitData.Description
			resp.Hits[i].Source.Plaintext = hitData.Plaintext
			resp.Hits[i].Source.Version = hitData.Version
			resp.Hits[i].Source.Platform = hitData.Platform

			if hitData.CoolDown != "" {
				resp.Hits[i].Source.Tags = append(resp.Hits[i].Source.Tags, fmt.Sprintf("冷却:%s", hitData.CoolDown))
			}
			resp.Hits[i].Source.Tags = append(resp.Hits[i].Source.Tags, fmt.Sprintf("Version:%s", hitData.Version))
		}
	}

	return resp, nil
}
func AutoComplete(ctx *context.Context, p *SearchParams) ([]map[string]interface{}, error) {
	esBuilder := new(es.Instance).Builder(p.Category).
		// 端游or手游
		SetTermQuery([]*es.TermQuery{
			&es.TermQuery{
				Name:  "platform",
				Value: p.Platform,
			},
		})
	if p.KeyWords != "" {
		// 按名字/介绍
		esBuilder = esBuilder.SetMultiMatchQuery(&es.MultiMatchQuery{
			Text:   p.KeyWords,
			Fields: p.Way,
		})
	}
	esIns, err := esBuilder.Build()
	if err != nil {
		return nil, err
	}
	resp, err := esIns.Query(ctx)
	if err != nil {
		return nil, err
	}

	var datas []map[string]interface{}
	for _, hit := range resp.Hits {
		sourceStr, _ := json.Marshal(hit.TmpSource)
		var hitData map[string]interface{}
		err = json.Unmarshal(sourceStr, &hitData)
		if err != nil {
			return nil, err
		}
		datas = append(datas, hitData)
	}
	return datas, nil
}

// BuildIndex 重建索引
func BuildIndex(ctx *context.Context, index string, rebuild bool) error {
	queue := make([]string, 0)
	// 如果没有指定index，就重建所有
	if index == "" {
		queue = config.GlobalConfig.ES.BuildIndex
	} else {
		queue = append(queue, index)
	}

	// 删除mapping
	if rebuild {
		if err := deleteIndex(ctx); err != nil {
			return err
		}
	}

	// 如果esmapping不存在就新建
	// 创建索引
	if err := createIndex(ctx); err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	cancelCtx, cancelFunc := context2.WithCancel(ctx)
	defer cancelFunc()

	select {
	case <-cancelCtx.Done():
		break
	default:
		// 表数据直接建索引
		for _, tbl := range queue {
			wg.Add(1)
			go func(tbl string) {
				err := mysql2es(ctx, tbl, &wg) // 这里的wg必须要传地址，不然值传递后传递的是副本
				if err != nil {
					log.Logger.Error(ctx, err)
					cancelFunc()
				}
			}(tbl)
		}
	}

	wg.Wait()
	return nil
}

func mysql2es(ctx *context.Context, tblName string, wg *sync.WaitGroup) error {
	defer wg.Done()

	var (
		heroesModel   model.LOLHeroes
		m_heroesModel model.LOLMHeroes

		equipModel   model.LOLEquipment
		m_equipModel model.LOLMEquipment

		runeModel   model.LOLRune
		m_runeModel model.LOLMRune

		skillModel   model.LOLSkill
		m_skillModel model.LOLMSkill
	)

	switch tblName {
	case heroesModel.TableName():
		log.Logger.Info(ctx, "开始处理:", heroesModel.TableName())
		if err := buildHeroesIndex(ctx); err != nil {
			return err
		}
	case m_heroesModel.TableName():
		log.Logger.Info(ctx, "开始处理:", m_heroesModel.TableName())
		if err := buildMHeroesIndex(ctx); err != nil {
			return err
		}
	case equipModel.TableName():
		log.Logger.Info(ctx, "开始处理:", equipModel.TableName())
		if err := buildEquipIndex(ctx); err != nil {
			return err
		}
	case m_equipModel.TableName():
		log.Logger.Info(ctx, "开始处理:", m_equipModel.TableName())
		if err := buildMEquipIndex(ctx); err != nil {
			return err
		}
	case runeModel.TableName():
		log.Logger.Info(ctx, "开始处理:", runeModel.TableName())
		if err := buildRuneIndex(ctx); err != nil {
			return err
		}
	case m_runeModel.TableName():
		log.Logger.Info(ctx, "开始处理:", m_runeModel.TableName())
		if err := buildMRuneIndex(ctx); err != nil {
			return err
		}
	case skillModel.TableName():
		log.Logger.Info(ctx, "开始处理:", skillModel.TableName())
		if err := buildSkillIndex(ctx); err != nil {
			return err
		}
	case m_skillModel.TableName():
		log.Logger.Info(ctx, "开始处理:", m_skillModel.TableName())
		if err := buildMSkillIndex(ctx); err != nil {
			return err
		}
	}
	return nil
}

func buildHeroesIndex(ctx *context.Context) error {
	d := dao.NewLOLHeroesDAO()
	rs, err := d.GetLOLHeroesMaxVersion()
	if err != nil {
		return err
	}

	data, err := d.GetLOLHeroesWithExt(rs.Version)
	if err != nil {
		return err
	}

	cancelCtx, cancelFunc := context2.WithCancel(ctx)
	defer cancelFunc()
	var (
		runningTask = make(chan struct{}, 10) //最多起100个协程处理索引
		successTask = int32(0)
		failTask    = int32(0)
		doneTask    = int32(0)
		allTask     = len(data)
		wg          = &sync.WaitGroup{}
	)

	select {
	case <-cancelCtx.Done():
		break
	default:
		spellDao := dao.NewHeroSpellDAO()
		hd := dao.CreateEsDao(dao.ESIndexHeroes)().(*dao.ESHeroesDAO)
		for _, row := range data {
			wg.Add(1)

			go func(row *model.LOLHeroesEXT) {
				defer func() {
					<-runningTask // 正在运行 -1
					atomic.AddInt32(&doneTask, 1)
					wg.Done()
				}()
				runningTask <- struct{}{}

				var esData []*model.ESHeroes
				tmp := row
				esHero := model.ESHeroes{
					ID:       tmp.HeroId,
					Name:     tmp.Name + " " + tmp.Title + "(" + tmp.Alias + ")",
					IconPath: tmp.Avatar,
					MainImg:  tmp.MainImg,
					Price:    "GoldPrice:" + tmp.GoldPrice + "/" + "CouponPrice:" + tmp.CouponPrice,
					Roles:    tmp.Roles,
					//Plaintext: "",
					Keywords: tmp.Keywords + "," + tmp.Alias + "," + tmp.Title,
					Version:  tmp.Version,
					FileTime: tmp.FileTime,
					Platform: strconv.Itoa(common.PlatformForLOL),
				}

				spells, err2 := spellDao.GetSpells(tmp.HeroId)
				if err2 != nil {
					log.Logger.Error(ctx, err2)
					atomic.AddInt32(&failTask, 1)
					cancelFunc()
					return
				}
				heroSpell := make([]*dto.HeroDescription, 0, 5)
				for _, spell := range spells {
					desc := &dto.HeroDescription{
						SpellKey:        spell.SpellKey,
						Sort:            spell.Sort,
						Name:            spell.Name,
						Description:     spell.Description,
						AbilityIconPath: spell.AbilityIconPath,
						Detail:          spell.Detail,
						Version:         spell.Version,
					}
					heroSpell = append(heroSpell, desc)
				}
				s, _ := json.Marshal(heroSpell)
				esHero.Description = string(s)
				esData = append(esData, &esHero)

				err3 := hd.Data2ES(ctx, esData)
				if err3 != nil {
					log.Logger.Error(ctx, err3)
					atomic.AddInt32(&failTask, 1)
					cancelFunc()
					return
				} else {
					atomic.AddInt32(&successTask, 1)
				}
			}(row)
		}
	}

	wg.Wait()
	log.Logger.Info(ctx, fmt.Sprintf("LOL Hero Task Done: allTask:%d, success:%d, fail:%d, done:%d", allTask, successTask, failTask, doneTask))
	return nil
}
func buildMHeroesIndex(ctx *context.Context) error {
	d := dao.NewLOLMHeroesDAO()
	rs, err := d.GetLOLMHeroesMaxVersion()
	if err != nil {
		return err
	}

	data, err := d.GetLOLMHeroesWithExt(rs.Version)
	if err != nil {
		return err
	}

	cancelCtx, cancelFunc := context2.WithCancel(ctx)
	defer cancelFunc()
	var (
		runningTask = make(chan struct{}, 10) //最多起100个协程处理索引
		successTask = int32(0)
		failTask    = int32(0)
		doneTask    = int32(0)
		allTask     = len(data)
		wg          = &sync.WaitGroup{}
	)

	select {
	case <-cancelCtx.Done():
		break
	default:
		spellDao := dao.NewHeroSpellDAO()
		hd := dao.CreateEsDao(dao.ESIndexHeroes)().(*dao.ESHeroesDAO)

		for _, row := range data {
			wg.Add(1)

			go func(row *model.LOLMHeroesEXT) {
				defer func() {
					<-runningTask // 正在运行 -1
					atomic.AddInt32(&doneTask, 1)
					wg.Done()
				}()
				runningTask <- struct{}{}

				var esData []*model.ESHeroes
				tmp := row
				esHero := model.ESHeroes{
					ID:       tmp.HeroId,
					Name:     tmp.Title + " " + tmp.Name + "(" + tmp.Alias + ")",
					IconPath: tmp.Avatar,
					MainImg:  tmp.Poster,
					Price:    "GoldPrice:" + tmp.Highlightprice + "/" + "CouponPrice:" + tmp.Couponprice,
					Roles:    tmp.Roles,
					//Plaintext: "",
					Keywords: tmp.Searchkey,
					Version:  tmp.Version,
					FileTime: tmp.FileTime,
					Platform: strconv.Itoa(common.PlatformForLOLM),
				}

				// 查询技能
				spells, err2 := spellDao.GetSpells(tmp.HeroId)
				if err2 != nil {
					log.Logger.Error(ctx, err2)
					atomic.AddInt32(&failTask, 1)
					cancelFunc()
					return
				}

				heroSpell := make([]*dto.HeroDescription, 0, 5)
				for _, spell := range spells {
					desc := &dto.HeroDescription{
						SpellKey:        spell.SpellKey,
						Sort:            spell.Sort,
						Name:            spell.Name,
						Description:     spell.Description,
						AbilityIconPath: spell.AbilityIconPath,
						Detail:          spell.Detail,
						Version:         spell.Version,
					}
					heroSpell = append(heroSpell, desc)
				}
				s, _ := json.Marshal(heroSpell)
				esHero.Description = string(s)
				esData = append(esData, &esHero)

				err3 := hd.Data2ES(ctx, esData)
				if err3 != nil {
					log.Logger.Error(ctx, err3)
					atomic.AddInt32(&failTask, 1)
					cancelFunc()
					return
				} else {
					atomic.AddInt32(&successTask, 1)
				}
			}(row)
		}
	}

	wg.Wait()
	log.Logger.Info(ctx, fmt.Sprintf("LOLM Hero Task Done: allTask:%d, success:%d, fail:%d, done:%d", allTask, successTask, failTask, doneTask))
	return nil
}

func buildEquipIndex(ctx *context.Context) error {
	d := dao.NewLOLEquipmentDAO()
	rs, err := d.GetLOLEquipmentMaxVersion()
	if err != nil {
		return err
	}

	data, err := d.GetLOLEquipmentWithExt(rs.Version)
	if err != nil {
		return err
	}

	cancelCtx, cancelFunc := context2.WithCancel(ctx)
	defer cancelFunc()
	var (
		runningTask = make(chan struct{}, 100) //最多起100个协程处理索引
		successTask = int32(0)
		failTask    = int32(0)
		doneTask    = int32(0)
		allTask     = len(data)
		wg          = &sync.WaitGroup{}
	)

	select {
	case <-cancelCtx.Done():
		break
	default:
		ed := dao.CreateEsDao(dao.ESIndexEquipment)().(*dao.ESEquipmentDAO)

		for _, row := range data {
			wg.Add(1)

			go func(row *model.LOLEquipment) {
				defer func() {
					<-runningTask // 正在运行 -1
					atomic.AddInt32(&doneTask, 1)
					wg.Done()
				}()

				runningTask <- struct{}{}

				var esEquip []*model.ESEquipment
				tmp := row
				esEquip = append(esEquip, &model.ESEquipment{
					ID:          tmp.ItemId,
					ItemId:      tmp.ItemId,
					Name:        tmp.Name,
					IconPath:    tmp.IconPath,
					Price:       tmp.Price,
					Description: tmp.Description,
					Plaintext:   tmp.Plaintext,
					Sell:        tmp.Sell,
					Total:       tmp.Total,
					SuitHeroId:  tmp.SuitHeroId,
					//SuitHeroName: tmp.SuitHeroId, // todo
					//SuitHeroIcon: tmp.SuitHeroId, // todo
					Keywords: tmp.Keywords,
					Maps:     tmp.Maps,
					//From:         tmp.From,  // todo
					//Into:         tmp.Into,  // todo
					//Types:        tmp.Types, // todo
					Version:  tmp.Version,
					FileTime: tmp.FileTime,
					Platform: strconv.Itoa(common.PlatformForLOL),
				})

				err2 := ed.Data2ES(ctx, esEquip)
				if err2 != nil {
					log.Logger.Error(ctx, err2)
					atomic.AddInt32(&failTask, 1)
					cancelFunc()
				} else {
					atomic.AddInt32(&successTask, 1)
				}
			}(row)
		}
	}

	wg.Wait()
	log.Logger.Info(ctx, fmt.Sprintf("LOL Equipment Task Done: allTask:%d, success:%d, fail:%d, done:%d", allTask, successTask, failTask, doneTask))
	return nil
}
func buildMEquipIndex(ctx *context.Context) error {
	d := dao.NewLOLMEquipmentDAO()
	rs, err := d.GetLOLMEquipmentMaxVersion()
	if err != nil {
		return err
	}

	data, err := d.GetLOLMEquipmentWithExt(rs.Version)
	if err != nil {
		return err
	}

	cancelCtx, cancelFunc := context2.WithCancel(ctx)
	defer cancelFunc()
	var (
		runningTask = make(chan struct{}, 100) //最多起100个协程处理索引
		successTask = int32(0)
		failTask    = int32(0)
		doneTask    = int32(0)
		allTask     = len(data)
		wg          = &sync.WaitGroup{}
	)

	select {
	case <-cancelCtx.Done():
		break
	default:
		ed := dao.CreateEsDao(dao.ESIndexEquipment)().(*dao.ESEquipmentDAO)

		for _, row := range data {
			wg.Add(1)

			go func(row *model.LOLMEquipment) {
				defer func() {
					<-runningTask // 正在运行 -1
					atomic.AddInt32(&doneTask, 1)
					wg.Done()
				}()

				runningTask <- struct{}{}

				var esEquip []*model.ESEquipment
				tmp := row

				esEquip = append(esEquip, &model.ESEquipment{
					ID:          tmp.EquipId,
					ItemId:      tmp.EquipId,
					Name:        tmp.Name,
					IconPath:    tmp.IconPath,
					Price:       tmp.Price,
					Description: tmp.Description,
					//Plaintext:    tmp.,
					//Sell:         tmp.,
					Total: tmp.Price,
					//SuitHeroId:   tmp.SuitHeroId,
					//SuitHeroName: tmp.SuitHeroId, // todo
					//SuitHeroIcon: tmp.SuitHeroId, // todo
					Keywords: tmp.SearchKey, // 仅LOLM字段
					Maps:     "召唤师峡谷",
					//From:     tmp.From, // todo
					//Into:     tmp.Into, // todo
					Types:    tmp.Type,
					Version:  tmp.Version,
					FileTime: tmp.FileTime,
					Platform: strconv.Itoa(common.PlatformForLOLM),
				})

				err2 := ed.Data2ES(ctx, esEquip)
				if err2 != nil {
					log.Logger.Error(ctx, err2)
					atomic.AddInt32(&failTask, 1)
					cancelFunc()
				} else {
					atomic.AddInt32(&successTask, 1)
				}
			}(row)
		}
	}

	wg.Wait()
	log.Logger.Info(ctx, fmt.Sprintf("LOLM Equipment Task Done: allTask:%d, success:%d, fail:%d, done:%d", allTask, successTask, failTask, doneTask))
	return nil
}

func buildRuneIndex(ctx *context.Context) error {
	d := dao.NewLOLRuneDAO()
	rs, err := d.GetLOLRuneMaxVersion()
	if err != nil {
		return err
	}

	data, err := d.GetLOLRune(rs.Version)
	if err != nil {
		return err
	}

	cancelCtx, cancelFunc := context2.WithCancel(ctx)
	defer cancelFunc()

	var (
		runningTask = make(chan struct{}, 10) //最多起100个协程处理索引
		successTask = int32(0)
		failTask    = int32(0)
		doneTask    = int32(0)
		allTask     = len(data)
		wg          = &sync.WaitGroup{}
	)

	select {
	case <-cancelCtx.Done():
		break
	default:
		rd := dao.CreateEsDao(dao.ESIndexRune)().(*dao.ESRuneDAO)
		for _, row := range data {
			wg.Add(1)
			go func(row *model.LOLRune) {
				defer func() {
					<-runningTask // 正在运行 -1
					atomic.AddInt32(&doneTask, 1)
					wg.Done()
				}()
				runningTask <- struct{}{} // 正在运行 +1

				var esData []*model.ESRune
				tmp := row
				esData = append(esData, &model.ESRune{
					ID:          tmp.RuneID,
					Name:        tmp.Name,
					IconPath:    tmp.Icon,
					Tooltip:     tmp.Tooltip,
					Description: tmp.Longdesc,
					Plaintext:   tmp.Shortdesc,
					Keywords:    tmp.Keywords,
					SlotLabel:   tmp.SlotLabel,
					StyleName:   tmp.StyleName,
					Version:     tmp.Version,
					FileTime:    tmp.FileTime,
					Platform:    strconv.Itoa(common.PlatformForLOL),
				})

				err := rd.Data2ES(ctx, esData)
				if err != nil {
					log.Logger.Error(ctx, err)
					atomic.AddInt32(&failTask, 1)
					cancelFunc()
					return
				} else {
					atomic.AddInt32(&successTask, 1)
				}
			}(row)
		}
	}
	wg.Wait()
	log.Logger.Info(ctx, fmt.Sprintf("LOL Rune Task Done: allTask:%d, success:%d, fail:%d, done:%d", allTask, successTask, failTask, doneTask))
	return nil
}
func buildMRuneIndex(ctx *context.Context) error {
	d := dao.NewLOLMRuneDAO()
	rs, err := d.GetLOLMRuneMaxVersion()
	if err != nil {
		return err
	}

	data, err := d.GetLOLMRune(rs.Version)
	if err != nil {
		return err
	}

	cancelCtx, cancelFunc := context2.WithCancel(ctx)
	defer cancelFunc()
	var (
		runningTask = make(chan struct{}, 100) //最多起100个协程处理索引
		successTask = int32(0)
		failTask    = int32(0)
		doneTask    = int32(0)
		allTask     = len(data)
		wg          = &sync.WaitGroup{}
	)

	select {
	case <-cancelCtx.Done():
		break
	default:
		rd := dao.CreateEsDao(dao.ESIndexRune)().(*dao.ESRuneDAO)
		for _, row := range data {
			wg.Add(1)

			go func(row *model.LOLMRune) {
				defer func() {
					<-runningTask // 正在运行 -1
					atomic.AddInt32(&doneTask, 1)
					wg.Done()
				}()
				runningTask <- struct{}{} // 正在运行 +1

				var esData []*model.ESRune
				tmp := row

				esData = append(esData, &model.ESRune{
					ID:          tmp.RuneId,
					Name:        tmp.Name,
					IconPath:    tmp.IconPath,
					Tooltip:     tmp.AttrName,
					Description: tmp.DetailInfo,
					Plaintext:   tmp.Description,
					Keywords:    tmp.Keywords,
					StyleName:   tmp.StyleName,
					//SlotLabel:   "",
					//StyleName:   "",
					//Maps:        "",
					Types:    tmp.Type,
					Version:  tmp.Version,
					FileTime: tmp.FileTime,
					Platform: strconv.Itoa(common.PlatformForLOLM),
				})

				err := rd.Data2ES(ctx, esData)
				if err != nil {
					log.Logger.Error(ctx, err)
					atomic.AddInt32(&failTask, 1)
					cancelFunc()
					return
				} else {
					atomic.AddInt32(&successTask, 1)
				}
			}(row)
		}
	}

	wg.Wait()
	log.Logger.Info(ctx, fmt.Sprintf("LOLM Rune Task Done: allTask:%d, success:%d, fail:%d, done:%d", allTask, successTask, failTask, doneTask))
	return nil
}

func buildSkillIndex(ctx *context.Context) error {
	d := dao.NewLOLSkillDAO()
	rs, err := d.GetLOLSkillMaxVersion()
	if err != nil {
		return err
	}

	data, err := d.GetLOLSkill(rs.Version)
	if err != nil {
		return err
	}

	cancelCtx, cancelFunc := context2.WithCancel(ctx)
	defer cancelFunc()
	var (
		runningTask = make(chan struct{}, 100) //最多起100个协程处理索引
		successTask = int32(0)
		failTask    = int32(0)
		doneTask    = int32(0)
		allTask     = len(data)
		wg          = &sync.WaitGroup{}
	)

	select {
	case <-cancelCtx.Done():
		break
	default:
		sd := dao.CreateEsDao(dao.ESIndexSkill)().(*dao.ESSkillDAO)
		for _, row := range data {
			wg.Add(1)

			go func(row *model.LOLSkill) {
				defer func() {
					<-runningTask // 正在运行 -1
					atomic.AddInt32(&doneTask, 1)
					wg.Done()
				}()
				runningTask <- struct{}{}

				var esData []*model.ESSkill
				tmp := row
				esData = append(esData, &model.ESSkill{
					ID:          tmp.SkillID,
					Name:        tmp.Name,
					IconPath:    tmp.Icon,
					Description: tmp.Description,
					Plaintext:   "",
					Keywords:    tmp.Keywords,
					Maps:        tmp.Gamemode,
					CoolDown:    tmp.Cooldown,
					Version:     tmp.Version,
					FileTime:    tmp.FileTime,
					Platform:    strconv.Itoa(common.PlatformForLOL),
				})

				err := sd.Data2ES(ctx, esData)
				if err != nil {
					log.Logger.Error(ctx, err)
					atomic.AddInt32(&failTask, 1)
					cancelFunc()
				} else {
					atomic.AddInt32(&successTask, 1)
				}
			}(row)
		}
	}

	wg.Wait()
	log.Logger.Info(ctx, fmt.Sprintf("LOL Skill Task Done: allTask:%d, success:%d, fail:%d, done:%d", allTask, successTask, failTask, doneTask))
	return nil
}
func buildMSkillIndex(ctx *context.Context) error {
	d := dao.NewLOLMSkillDAO()
	rs, err := d.GetLOLMSkillMaxVersion()
	if err != nil {
		return err
	}

	data, err := d.GetLOLMSkill(rs.Version)
	if err != nil {
		return err
	}

	cancelCtx, cancelFunc := context2.WithCancel(ctx)
	defer cancelFunc()
	var (
		runningTask = make(chan struct{}, 100) //最多起100个协程处理索引
		successTask = int32(0)
		failTask    = int32(0)
		doneTask    = int32(0)
		allTask     = len(data)
		wg          = &sync.WaitGroup{}
	)

	select {
	case <-cancelCtx.Done():
		break
	default:
		sd := dao.CreateEsDao(dao.ESIndexSkill)().(*dao.ESSkillDAO)

		for _, row := range data {
			wg.Add(1)

			go func(row *model.LOLMSkill) {
				defer func() {
					<-runningTask // 正在运行 -1
					atomic.AddInt32(&doneTask, 1)
					wg.Done()
				}()
				runningTask <- struct{}{}

				var esData []*model.ESSkill
				tmp := row
				esData = append(esData, &model.ESSkill{
					ID:          tmp.SkillID,
					Name:        tmp.Name,
					IconPath:    tmp.IconPath,
					Description: tmp.FuncDesc,
					Plaintext:   "",
					Keywords:    tmp.Keywords,
					Maps:        tmp.Mode,
					CoolDown:    tmp.Cd,
					Version:     tmp.Version,
					FileTime:    tmp.FileTime,
					Platform:    strconv.Itoa(common.PlatformForLOLM),
				})

				err := sd.Data2ES(ctx, esData)
				if err != nil {
					log.Logger.Error(ctx, err)
					atomic.AddInt32(&failTask, 1)
					cancelFunc()
				} else {
					atomic.AddInt32(&successTask, 1)
				}

			}(row)
		}
	}

	wg.Wait()
	log.Logger.Info(ctx, fmt.Sprintf("LOLM Skill Task Done: allTask:%d, success:%d, fail:%d, done:%d", allTask, successTask, failTask, doneTask))
	return nil
}

// 创建索引
func createIndex(ctx *context.Context) error {
	// equipment
	if err := dao.NewESEquipmentDAO()().CreateIndex(ctx); err != nil {
		return err
	}
	// heroes
	if err := dao.NewESHeroesDAO()().CreateIndex(ctx); err != nil {
		return err
	}
	// rune
	if err := dao.NewESRuneDAO()().CreateIndex(ctx); err != nil {
		return err
	}

	// skill
	if err := dao.NewESSkillDAO()().CreateIndex(ctx); err != nil {
		return err
	}

	return nil
}

func deleteIndex(ctx *context.Context) error {
	// equipment
	if err := dao.NewESEquipmentDAO()().DeleteIndex(ctx); err != nil {
		return err
	}
	// heroes
	if err := dao.NewESHeroesDAO()().DeleteIndex(ctx); err != nil {
		return err
	}
	// rune
	if err := dao.NewESRuneDAO()().DeleteIndex(ctx); err != nil {
		return err
	}

	// skill
	if err := dao.NewESSkillDAO()().DeleteIndex(ctx); err != nil {
		return err
	}

	return nil
}
