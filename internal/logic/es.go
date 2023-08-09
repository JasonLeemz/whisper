package logic

import (
	"encoding/json"
	errors2 "errors"
	"fmt"
	"strconv"
	"whisper/internal/logic/common"
	"whisper/internal/model"
	"whisper/internal/model/DAO"
	"whisper/pkg/es"

	"github.com/olivere/elastic/v7"
	"whisper/pkg/config"
	"whisper/pkg/context"
	"whisper/pkg/log"
)

type SearchParams struct {
	KeyWords string   `json:"key_words"`
	Platform string   `json:"platform,omitempty"`
	Category string   `json:"category,omitempty"`
	Way      []string `json:"way,omitempty"`
	Map      []string `json:"map,omitempty"`
}

func EsSearch(ctx *context.Context, p *SearchParams) (*common.EsResultHits, error) {
	indexName := p.Category
	if indexName == "" {
		return nil, errors2.New("indexName is nil")
	}

	query := elastic.NewBoolQuery()

	// 高亮搜索结果
	hl := elastic.NewHighlight()
	fields := make([]*elastic.HighlighterField, 0, len(p.Way))

	// 按名字介绍
	way := make([]string, 0)
	for _, w := range p.Way {
		way = append(way, w)
		fields = append(fields, elastic.NewHighlighterField(w))
	}
	query = query.Must(elastic.NewMultiMatchQuery(p.KeyWords, way...))
	hl = hl.Fields(fields...)
	hl = hl.PreTags("<em>").PostTags("</em>")

	equipModel := new(model.ESEquipment)
	// 按地图
	if indexName == equipModel.GetIndexName() {
		maps := make([]interface{}, 0)
		for _, m := range p.Map {
			maps = append(maps, m)
		}
		query = query.Must(elastic.NewTermsQuery("maps", maps...))
	}

	// 端游or手游
	query = query.Must(elastic.NewTermQuery("platform", p.Platform))

	//query = query.Filter(elastic.NewRangeQuery("id").Gte(0))
	//query = query.Filter(elastic.NewRangeQuery("id").Lte(9999999))

	sortByScore := elastic.NewFieldSort("_score").Desc()

	res, err := es.ESClient.Search().
		Index(indexName).
		Highlight(hl).
		Query(query).
		SortBy(sortByScore).
		From(0).Size(10).
		Pretty(true).
		Do(ctx)
	if err != nil {
		return nil, err
	}

	resp := common.EsResultHits{}
	data, _ := json.Marshal(res.Hits)
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}

	switch indexName {
	case equipModel.GetIndexName():
		for i, hit := range resp.Hits {
			sourceStr, _ := json.Marshal(hit.TmpSource)
			hitData := model.ESEquipment{}
			err = json.Unmarshal(sourceStr, &hitData)
			if err != nil {
				return nil, err
			}
			resp.Hits[i].Source.Name = hitData.Name
			resp.Hits[i].Source.IconPath = hitData.IconPath
			resp.Hits[i].Source.Description = hitData.Description
			resp.Hits[i].Source.Plaintext = hitData.Plaintext
			resp.Hits[i].Source.Version = hitData.Version
			resp.Hits[i].Source.Tags = append(resp.Hits[i].Source.Tags, fmt.Sprintf("Price:%s", hitData.Price))
			resp.Hits[i].Source.Tags = append(resp.Hits[i].Source.Tags, fmt.Sprintf("Sell:%s", hitData.Sell))
			resp.Hits[i].Source.Tags = append(resp.Hits[i].Source.Tags, fmt.Sprintf("Total:%s", hitData.Total))
			resp.Hits[i].Source.Tags = append(resp.Hits[i].Source.Tags, fmt.Sprintf("Maps:%s", hitData.Maps))
			resp.Hits[i].Source.Tags = append(resp.Hits[i].Source.Tags, fmt.Sprintf("Version:%s", hitData.Version))
		}
	case new(model.ESHeroes).GetIndexName():
		for i, hit := range resp.Hits {
			sourceStr, _ := json.Marshal(hit.TmpSource)
			hitData := model.ESHeroes{}
			err = json.Unmarshal(sourceStr, &hitData)
			if err != nil {
				return nil, err
			}
			resp.Hits[i].Source.Name = hitData.Name
			resp.Hits[i].Source.IconPath = hitData.IconPath
			resp.Hits[i].Source.Description = hitData.Description
			resp.Hits[i].Source.Plaintext = hitData.Plaintext
			resp.Hits[i].Source.Version = hitData.Version
			resp.Hits[i].Source.Tags = append(resp.Hits[i].Source.Tags, fmt.Sprintf("Price:%s", hitData.Price))
			resp.Hits[i].Source.Tags = append(resp.Hits[i].Source.Tags, fmt.Sprintf("Types:%s", hitData.Types))
			resp.Hits[i].Source.Tags = append(resp.Hits[i].Source.Tags, fmt.Sprintf("Version:%s", hitData.Version))
		}
	case new(model.ESRune).GetIndexName():
		for i, hit := range resp.Hits {
			sourceStr, _ := json.Marshal(hit.TmpSource)
			hitData := model.ESRune{}
			err = json.Unmarshal(sourceStr, &hitData)
			if err != nil {
				return nil, err
			}
			resp.Hits[i].Source.Name = hitData.Name
			resp.Hits[i].Source.IconPath = hitData.IconPath
			resp.Hits[i].Source.Description = hitData.Description
			resp.Hits[i].Source.Plaintext = hitData.Plaintext
			resp.Hits[i].Source.Version = hitData.Version
			resp.Hits[i].Source.Tags = append(resp.Hits[i].Source.Tags, fmt.Sprintf("主系:%s", hitData.StyleName))
			resp.Hits[i].Source.Tags = append(resp.Hits[i].Source.Tags, fmt.Sprintf("SlotLabel:%s", hitData.SlotLabel))
			resp.Hits[i].Source.Tags = append(resp.Hits[i].Source.Tags, fmt.Sprintf("Version:%s", hitData.Version))
		}
	case new(model.ESSkill).GetIndexName():
		for i, hit := range resp.Hits {
			sourceStr, _ := json.Marshal(hit.TmpSource)
			hitData := model.ESSkill{}
			err = json.Unmarshal(sourceStr, &hitData)
			if err != nil {
				return nil, err
			}
			resp.Hits[i].Source.Name = hitData.Name
			resp.Hits[i].Source.IconPath = hitData.IconPath
			resp.Hits[i].Source.Description = hitData.Description
			resp.Hits[i].Source.Plaintext = hitData.Plaintext
			resp.Hits[i].Source.Version = hitData.Version
			resp.Hits[i].Source.Tags = append(resp.Hits[i].Source.Tags, fmt.Sprintf("冷却:%s", hitData.CoolDown))
			resp.Hits[i].Source.Tags = append(resp.Hits[i].Source.Tags, fmt.Sprintf("Version:%s", hitData.Version))
		}
	}

	return &resp, nil
}

// BuildIndex 重建索引
func BuildIndex(ctx *context.Context, index string) error {
	queue := make([]string, 0)
	// 如果没有指定index，就重建所有
	if index == "" {
		queue = config.GlobalConfig.ES.BuildIndex
	} else {
		queue = append(queue, index)
	}

	// 表数据直接建索引
	for _, tbl := range queue {
		err := mysql2es(ctx, tbl)
		if err != nil {
			log.Logger.Error(ctx, err)
			return err
		}
	}
	return nil
}

var equipChan = make(chan []*model.ESEquipment, 100) // 最多起100个协程处理索引
var heroesChan = make(chan []*model.ESHeroes, 100)   // 最多起100个协程处理索引
var runeChan = make(chan []*model.ESRune, 100)       // 最多起100个协程处理索引
var skillChan = make(chan []*model.ESSkill, 100)     // 最多起100个协程处理索引

func mysql2es(ctx *context.Context, tblName string) error {
	// 创建索引
	if err := createIndex(ctx); err != nil {
		return err
	}
	go equipCustomer(ctx)
	go heroesCustomer(ctx)
	go runeCustomer(ctx)
	go skillCustomer(ctx)

	var equipModel model.LOLEquipment
	var m_equipModel model.LOLMEquipment

	var heroesModel model.LOLHeroes
	var m_heroesModel model.LOLMHeroes

	var runeModel model.LOLRune
	var m_runeModel model.LOLMRune

	var skillModel model.LOLSkill
	var m_skillModel model.LOLMSkill

	switch tblName {
	case runeModel.TableName():
		if err := runeProduce(); err != nil {
			return err
		}
	case m_runeModel.TableName():
		if err := runeMProduce(); err != nil {
			return err
		}
	case skillModel.TableName():

		if err := skillProduce(); err != nil {
			return err
		}
	case m_skillModel.TableName():
		if err := skillMProduce(); err != nil {
			return err
		}
	case heroesModel.TableName():
		if err := heroesProduce(); err != nil {
			return err
		}
	case m_heroesModel.TableName():
		if err := heroesMProduce(); err != nil {
			return err
		}
	case m_equipModel.TableName():
		if err := equipMProduce(); err != nil {
			return err
		}
	case equipModel.TableName():
		if err := equipProduce(); err != nil {
			return err
		}
	}
	return nil
}
func heroesMProduce() error {
	d := dao.NewLOLMHeroesDAO()
	rs, err := d.GetLOLMHeroesMaxVersion()
	if err != nil {
		return err
	}

	data, err := d.GetLOLMHeroes(rs.Version)
	if err != nil {
		return err
	}

	go func() {
		for _, d := range data {
			var esData []*model.ESHeroes
			tmp := d
			esData = append(esData, &model.ESHeroes{
				ID:       tmp.HeroId + "_" + tmp.Version,
				Name:     tmp.Title + " " + tmp.Name,
				IconPath: tmp.Avatar,
				Price:    "GoldPrice:" + tmp.Highlightprice + "/" + "CouponPrice:" + tmp.Couponprice,
				//Description: tmp., TODO
				//Plaintext: "",
				Keywords: tmp.Alias + "," + tmp.Title,
				//Maps:      "",
				//Types:    "",
				Version:  tmp.Version,
				FileTime: tmp.FileTime,
				Platform: strconv.Itoa(common.PlatformForLOLM),
			})

			// 放入阻塞队列
			heroesChan <- esData

		}
	}()
	return nil
}
func heroesProduce() error {
	d := dao.NewLOLHeroesDAO()
	rs, err := d.GetLOLHeroesMaxVersion()
	if err != nil {
		return err
	}

	data, err := d.GetLOLHeroes(rs.Version)
	if err != nil {
		return err
	}

	go func() {
		for _, d := range data {
			var esData []*model.ESHeroes
			tmp := d
			esData = append(esData, &model.ESHeroes{
				ID:       tmp.HeroId + "_" + tmp.Version,
				Name:     tmp.Name + " " + tmp.Title,
				IconPath: fmt.Sprintf("https://game.gtimg.cn/images/lol/act/img/skin/small%s000.jpg", tmp.HeroId),
				Price:    "GoldPrice:" + tmp.GoldPrice + "/" + "CouponPrice:" + tmp.CouponPrice,
				//Description: tmp., TODO
				//Plaintext: "",
				Keywords: tmp.Keywords + "," + tmp.Alias + "," + tmp.Title + "," + tmp.Camp,
				//Maps:      "",
				//Types:    "",
				Version:  tmp.Version,
				FileTime: tmp.FileTime,
				Platform: strconv.Itoa(common.PlatformForLOL),
			})

			// 放入阻塞队列
			heroesChan <- esData

		}
	}()
	return nil
}
func heroesCustomer(ctx *context.Context) {
	var esData []*model.ESHeroes
	ed := dao.NewESHeroesDAO()
	for {
		esData = <-heroesChan
		// LOLEquipment2ES 支持批量索引，这里避免占用内存过大每次只处理一行数据
		err := ed.Heroes2ES(ctx, esData)
		if err != nil {
			log.Logger.Error(ctx, err)
		}

		// TODO
	}
}

func equipProduce() error {
	d := dao.NewLOLEquipmentDAO()
	rs, err := d.GetLOLEquipmentMaxVersion()
	if err != nil {
		return err
	}

	equipment, err := d.GetLOLEquipment(rs.Version)
	if err != nil {
		return err
	}

	go func() {
		for _, equip := range equipment {
			var esEquip []*model.ESEquipment
			tmp := equip
			esEquip = append(esEquip, &model.ESEquipment{
				ID:           tmp.ItemId + "_" + tmp.Maps,
				EquipId:      tmp.ItemId,
				Name:         tmp.Name,
				IconPath:     tmp.IconPath,
				Price:        tmp.Price,
				Description:  tmp.Description,
				Plaintext:    tmp.Plaintext,
				Sell:         tmp.Sell,
				Total:        tmp.Total,
				SuitHeroId:   tmp.SuitHeroId,
				SuitHeroName: tmp.SuitHeroId, // todo
				SuitHeroIcon: tmp.SuitHeroId, // todo
				Keywords:     tmp.Keywords,
				Maps:         tmp.Maps,
				From:         tmp.From,  // todo
				Into:         tmp.Into,  // todo
				Types:        tmp.Types, // todo
				Version:      tmp.Version,
				FileTime:     tmp.FileTime,
				Platform:     strconv.Itoa(common.PlatformForLOL),
			})

			// 放入阻塞队列
			equipChan <- esEquip

		}
	}()
	return nil
}
func equipMProduce() error {
	d := dao.NewLOLMEquipmentDAO()
	rs, err := d.GetLOLMEquipmentMaxVersion()
	if err != nil {
		return err
	}

	equipment, err := d.GetLOLMEquipment(rs.Version)
	if err != nil {
		return err
	}

	go func() {
		for _, equip := range equipment {
			var esEquip []*model.ESEquipment
			tmp := equip
			esEquip = append(esEquip, &model.ESEquipment{
				ID:          tmp.EquipId + "_lolm",
				EquipId:     tmp.EquipId,
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
				Keywords: tmp.SearchKey, // 数据为空
				Maps:     "召唤师峡谷",
				From:     tmp.From, // todo
				Into:     tmp.Into, // todo
				Types:    tmp.Type,
				Version:  tmp.Version,
				FileTime: tmp.FileTime,
				Platform: strconv.Itoa(common.PlatformForLOLM),
			})

			// 放入阻塞队列
			equipChan <- esEquip
		}
	}()

	return nil
}
func equipCustomer(ctx *context.Context) {
	var equips []*model.ESEquipment
	ed := dao.NewESEquipmentDAO()
	for {
		equips = <-equipChan
		// LOLEquipment2ES 支持批量索引，这里避免占用内存过大每次只处理一行数据
		err := ed.Equipment2ES(ctx, equips)
		if err != nil {
			log.Logger.Error(ctx, err)
		}

		// TODO
	}
}

func runeProduce() error {
	d := dao.NewLOLRuneDAO()
	rs, err := d.GetLOLRuneMaxVersion()
	if err != nil {
		return err
	}

	data, err := d.GetLOLRune(rs.Version)
	if err != nil {
		return err
	}

	go func() {
		for _, row := range data {
			var esData []*model.ESRune
			tmp := row
			esData = append(esData, &model.ESRune{
				ID:          strconv.Itoa(int(tmp.Id)) + "_" + tmp.Name + "_" + tmp.Version,
				Name:        tmp.Name,
				IconPath:    tmp.Icon,
				Tooltip:     tmp.Tooltip,
				Description: tmp.Longdesc,
				Plaintext:   tmp.Shortdesc,
				Keywords:    tmp.Key,
				SlotLabel:   tmp.SlotLabel,
				StyleName:   tmp.StyleName,
				Version:     tmp.Version,
				FileTime:    tmp.FileTime,
				Platform:    strconv.Itoa(common.PlatformForLOL),
			})

			// 放入阻塞队列
			runeChan <- esData

		}
	}()
	return nil
}
func runeMProduce() error {
	d := dao.NewLOLMRuneDAO()
	rs, err := d.GetLOLMRuneMaxVersion()
	if err != nil {
		return err
	}

	data, err := d.GetLOLMRune(rs.Version)
	if err != nil {
		return err
	}

	go func() {
		for _, row := range data {
			var esData []*model.ESRune
			tmp := row
			esData = append(esData, &model.ESRune{
				ID:          tmp.RuneId + "_" + tmp.Name + "_" + tmp.Version,
				Name:        tmp.Name,
				IconPath:    tmp.IconPath,
				Tooltip:     tmp.Description,
				Description: tmp.DetailInfo,
				Plaintext:   tmp.AttrName,
				Keywords:    tmp.Name,
				SlotLabel:   "",
				StyleName:   tmp.Type,
				Maps:        "",
				Types:       tmp.Type,
				Version:     tmp.Version,
				FileTime:    tmp.FileTime,
				Platform:    strconv.Itoa(common.PlatformForLOLM),
			})

			// 放入阻塞队列
			runeChan <- esData

		}
	}()
	return nil
}
func runeCustomer(ctx *context.Context) {
	var esData []*model.ESRune
	ed := dao.NewESRuneDAO()
	for {
		esData = <-runeChan
		// Rune2ES 支持批量索引，这里避免占用内存过大每次只处理一行数据
		err := ed.Rune2ES(ctx, esData)
		if err != nil {
			log.Logger.Error(ctx, err)
		}

		// TODO
	}
}

func skillProduce() error {
	d := dao.NewLOLSkillDAO()
	rs, err := d.GetLOLSkillMaxVersion()
	if err != nil {
		return err
	}

	data, err := d.GetLOLSkill(rs.Version)
	if err != nil {
		return err
	}

	go func() {
		for _, row := range data {
			var esData []*model.ESSkill
			tmp := row
			esData = append(esData, &model.ESSkill{
				ID:          strconv.Itoa(int(tmp.Id)) + "_" + tmp.Name,
				Name:        tmp.Name,
				IconPath:    tmp.Icon,
				Description: tmp.Description,
				Plaintext:   "",
				Keywords:    tmp.Name,
				Maps:        tmp.Gamemode,
				CoolDown:    tmp.Cooldown,
				Version:     tmp.Version,
				FileTime:    tmp.FileTime,
				Platform:    strconv.Itoa(common.PlatformForLOL),
			})

			// 放入阻塞队列
			skillChan <- esData

		}
	}()
	return nil
}
func skillMProduce() error {
	d := dao.NewLOLMSkillDAO()
	rs, err := d.GetLOLMSkillMaxVersion()
	if err != nil {
		return err
	}

	data, err := d.GetLOLMSkill(rs.Version)
	if err != nil {
		return err
	}

	go func() {
		for _, row := range data {
			var esData []*model.ESSkill
			tmp := row
			esData = append(esData, &model.ESSkill{
				ID:          tmp.SkillId + "_" + tmp.Name,
				Name:        tmp.Name,
				IconPath:    tmp.IconPath,
				Description: tmp.FuncDesc,
				Plaintext:   "",
				Keywords:    tmp.Name,
				Maps:        tmp.Mode,
				CoolDown:    tmp.Cd,
				Version:     tmp.Version,
				FileTime:    tmp.FileTime,
				Platform:    strconv.Itoa(common.PlatformForLOLM),
			})

			// 放入阻塞队列
			skillChan <- esData

		}
	}()
	return nil
}
func skillCustomer(ctx *context.Context) {
	var esData []*model.ESSkill
	ed := dao.NewESSkillDAO()
	for {
		esData = <-skillChan
		// Rune2ES 支持批量索引，这里避免占用内存过大每次只处理一行数据
		err := ed.Skill2ES(ctx, esData)
		if err != nil {
			log.Logger.Error(ctx, err)
		}

		// TODO
	}
}

// 创建索引
func createIndex(ctx *context.Context) error {
	// equipment
	if err := dao.NewESEquipmentDAO().CreateIndex(ctx); err != nil {
		return err
	}
	// heroes
	if err := dao.NewESHeroesDAO().CreateIndex(ctx); err != nil {
		return err
	}
	// rune
	if err := dao.NewESRuneDAO().CreateIndex(ctx); err != nil {
		return err
	}

	// skill
	if err := dao.NewESSkillDAO().CreateIndex(ctx); err != nil {
		return err
	}

	return nil
}
