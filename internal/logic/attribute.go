package logic

import (
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"
	"whisper/internal/dto"
	"whisper/internal/logic/common"
	"whisper/internal/model"
	dao "whisper/internal/model/DAO"
	"whisper/internal/service"
	"whisper/pkg/context"
	"whisper/pkg/log"
)

// HeroAttribute
// 根据传过来的id获取数据
// 只有一个id时直接返回
// 多个id时放入队列: 等待全部执行结束，或者有报错时退出
func HeroAttribute(ctx *context.Context, heroID string, platform int) (*dto.HeroAttribute, error) {
	if heroID != "0" {
		attribute, err := QueryHeroAttribute(ctx, heroID, platform)
		if err != nil {
			log.Logger.Error(ctx, err)
			return nil, err
		}
		err = recordHeroRoleAndSpellAndSkin(ctx, attribute, platform)
		if err != nil {
			log.Logger.Error(ctx, err)
			return nil, err
		}
		return attribute, err
	}

	var (
		runningErr  = make(chan error)        // 记录任务报错
		runningChan = make(chan struct{}, 50) // 控制同时并行的任务
		wg          = sync.WaitGroup{}        // 等待所有任务完成
		succNum     = int32(0)                // 记录成功完成的任务数
		doneNum     = int32(0)                // 记录做过的任务数
		failNum     = int32(0)                // 失败任务数
	)

	heroIDs, err := getAllHeroIDs(platform)
	if err != nil {
		return nil, err
	}

	// 防止runningChan阻塞影响主协程运行，放到go func中执行
	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
		}()

		for _, id := range heroIDs {
			runningChan <- struct{}{} // 正在运行+1
			doneNum++
			wg.Add(1)

			// do task
			go func(id string) {
				defer func() {
					wg.Done()
					<-runningChan // 正在运行-1
				}()

				attribute, err := QueryHeroAttribute(ctx, id, platform)
				if err == nil {
					if err := recordHeroRoleAndSpellAndSkin(ctx, attribute, platform); err != nil {
						runningErr <- err
						atomic.AddInt32(&failNum, 1)
					} else {
						atomic.AddInt32(&succNum, 1) // 执行成功的task数量+1
					}
				} else {
					atomic.AddInt32(&failNum, 1)
					runningErr <- err
				}
			}(id)
		}
	}()

	// 等待任务结束
	go func() {
		wg.Wait()
		close(runningErr)
	}()

	// 打印执行过程中的错误
	for err := range runningErr {
		log.Logger.Error(ctx, err)
	}

	log.Logger.Info(ctx, fmt.Sprintf("HeroAttribute Done"))
	log.Logger.Info(ctx, fmt.Sprintf("处理了: %d 个任务", doneNum))
	log.Logger.Info(ctx, fmt.Sprintf("提前结束,执行出错: %d 个任务", failNum))
	log.Logger.Info(ctx, fmt.Sprintf("成功执行了: %d 个任务", succNum))
	log.Logger.Info(ctx, fmt.Sprintf("剩余: %d 个任务待处理", int32(len(heroIDs))-doneNum))

	return nil, nil
}

func QueryHeroAttribute(ctx *context.Context, heroID string, platform int) (*dto.HeroAttribute, error) {
	if platform == common.PlatformForLOL {
		return service.GetLOLHeroAttribute(ctx, heroID)
	} else {
		return service.GetLOLMHeroAttribute(ctx, heroID)
	}
}

func recordHeroRoleAndSpellAndSkin(ctx *context.Context, data *dto.HeroAttribute, platform int) error {
	// 记录HeroRole
	if err := recordHeroRole(ctx, data, platform); err != nil {
		return err
	}

	// 记录HeroSpell
	if err := recordHeroSpell(ctx, data, platform); err != nil {
		return err
	}

	// 记录HeroSkin
	if err := recordHeroSkin(ctx, data, platform); err != nil {
		return err
	}

	// 记录HeroAttr
	avatar := data.Hero.Avatar
	mainImg := ""
	if len(data.Skins) > 0 {
		avatar = data.Skins[0].IconImg
		mainImg = data.Skins[0].MainImg
	}
	attr := &model.HeroAttribute{
		HeroId:              data.Hero.HeroId,
		Title:               data.Hero.Title,
		Name:                data.Hero.Name,
		Alias:               data.Hero.Alias,
		ShortBio:            data.Hero.ShortBio,
		Defense:             data.Hero.Defense,
		Magic:               data.Hero.Magic,
		Difficulty:          data.Hero.Difficulty,
		DifficultyL:         data.Hero.DifficultyL,
		Attack:              data.Hero.Attack,
		Attackrange:         data.Hero.AttackRange,
		Attackdamage:        data.Hero.AttackDamage,
		Attackspeed:         data.Hero.Attackspeed,
		Attackspeedperlevel: data.Hero.Attackspeedperlevel,
		Hp:                  data.Hero.Hp,
		Hpperlevel:          data.Hero.Hpperlevel,
		Mp:                  data.Hero.Mp,
		Mpperlevel:          data.Hero.Mpperlevel,
		Movespeed:           data.Hero.Movespeed,
		Armor:               data.Hero.Armor,
		Armorperlevel:       data.Hero.Armorperlevel,
		Spellblock:          data.Hero.Spellblock,
		Spellblockperlevel:  data.Hero.Spellblockperlevel,
		Hpregen:             data.Hero.Hpregen,
		Hpregenperlevel:     data.Hero.Hpregenperlevel,
		Mpregen:             data.Hero.Mpregen,
		Mpregenperlevel:     data.Hero.Mpregenperlevel,
		Crit:                data.Hero.Crit,
		Damage:              data.Hero.Damage,
		Durability:          data.Hero.Durability,
		Mobility:            data.Hero.Mobility,
		Avatar:              avatar,
		MainImg:             mainImg,
		Highlightprice:      data.Hero.Highlightprice,
		GoldPrice:           data.Hero.GoldPrice,
		Couponprice:         data.Hero.Couponprice,
		IsWeekFree:          data.Hero.IsWeekFree,
		Platform:            platform,
		Version:             data.Version,
		FileTime:            data.FileTime,
	}
	ha := dao.NewHeroAttributeDAO()
	err3 := ha.DeleteAndInsert(map[string]interface{}{
		"heroId": data.Hero.HeroId,
	}, []*model.HeroAttribute{attr})
	if err3 != nil {
		return err3
	}

	return nil
}
func recordHeroRole(ctx *context.Context, data *dto.HeroAttribute, platform int) error {
	hrs := make([]*model.HeroRole, 0, len(data.Hero.Roles))
	for i, _ := range data.Hero.Roles {
		hrs = append(hrs, &model.HeroRole{
			Platform: platform,
			HeroId:   data.Hero.HeroId,
			Role:     data.Hero.Roles[i],
		})
	}
	hrdao := dao.NewHeroRoleDAO()
	err := hrdao.DeleteAndInsert(map[string]interface{}{
		"hero_id": data.Hero.HeroId,
	}, hrs)
	if err != nil {
		return err
	}

	return nil
}

// recordHeroSpell ...
func recordHeroSpell(ctx *context.Context, data *dto.HeroAttribute, platform int) error {
	hrs := make([]*model.HeroSpell, 0, len(data.Spells))
	var keyMap = map[string]int{
		"passive": 0,
		"q":       1,
		"w":       2,
		"e":       3,
		"r":       4,
	}
	apellMap := make(map[string]bool)
	for i, spell := range data.Spells {
		if _, ok := apellMap[spell.Name]; ok {
			continue
		} else {
			apellMap[spell.Name] = true
		}
		detail, _ := json.Marshal(spell.Detail)
		hs := &model.HeroSpell{
			HeroId:          data.Hero.HeroId,
			SpellKey:        spell.SpellKey,
			Sort:            keyMap[spell.SpellKey],
			Name:            spell.Name,
			Description:     spell.Description,
			AbilityIconPath: spell.AbilityIconPath,
			Detail:          string(detail),
			Platform:        platform,
			Version:         data.Version,
			FileTime:        data.FileTime,
		}
		if sort, ok := keyMap[spell.SpellKey]; ok {
			hs.Sort = sort
		} else {
			hs.Sort = -1
			if platform == common.PlatformForLOLM {
				hs.Sort = i
			}
		}
		hrs = append(hrs, hs)
	}

	hrdao := dao.NewHeroSpellDAO()
	err := hrdao.DeleteAndInsert(map[string]interface{}{
		"heroId": data.Hero.HeroId,
	}, hrs)
	if err != nil {
		return err
	}

	return nil
}

// recordHeroSkin ...
func recordHeroSkin(ctx *context.Context, data *dto.HeroAttribute, platform int) error {
	if platform == common.PlatformForLOLM {
		// 手游皮肤获取逻辑不同，这里只处理端游
		return nil
	}
	hrs := make([]*model.HeroSkin, 0, len(data.Skins))

	for _, skin := range data.Skins {
		if skin.MainImg == "" {
			continue
		}
		hs := &model.HeroSkin{
			HeroId:      data.Hero.HeroId,
			SkinId:      skin.SkinId,
			HeroName:    skin.HeroName,
			HeroTitle:   skin.HeroTitle,
			Name:        skin.Name,
			IsBase:      skin.IsBase,
			EmblemsName: skin.EmblemsName,
			Description: skin.Description,
			MainImg:     skin.MainImg,
			IconImg:     skin.IconImg,
			LoadingImg:  skin.LoadingImg,
			VideoImg:    skin.VideoImg,
			SourceImg:   skin.SourceImg,
			Platform:    platform,
			Version:     data.Version,
			FileTime:    data.FileTime,
		}

		hrs = append(hrs, hs)
	}

	hrdao := dao.NewHeroSkinDAO()
	err := hrdao.DeleteAndInsert(map[string]interface{}{
		"heroId": data.Hero.HeroId,
	}, hrs)
	if err != nil {
		return err
	}

	return nil
}

func getAllHeroIDs(platform int) ([]string, error) {
	// 这里必须从heroes表获取
	var heroIDs []string
	// 获取所有英雄ID
	if platform == common.PlatformForLOL {
		heroesDao := dao.NewLOLHeroesDAO()
		v, err := heroesDao.GetLOLHeroesMaxVersion()
		if err != nil {
			return nil, err
		}
		heroes, err := heroesDao.GetLOLHeroes(v.Version)
		if err != nil {
			return nil, err
		}
		heroIDs = make([]string, 0, len(heroes))
		for i, _ := range heroes {
			heroIDs = append(heroIDs, heroes[i].HeroId)
		}
	} else {
		heroesDao := dao.NewLOLMHeroesDAO()
		v, err := heroesDao.GetLOLMHeroesMaxVersion()
		if err != nil {
			return nil, err
		}
		heroes, err := heroesDao.GetLOLMHeroes(v.Version)
		if err != nil {
			return nil, err
		}
		heroIDs = make([]string, 0, len(heroes))
		for i, _ := range heroes {
			heroIDs = append(heroIDs, heroes[i].HeroId)
		}
	}

	return heroIDs, nil
}

func GetVersion(ctx *context.Context) []*model.HeroAttribute {
	ad := dao.NewHeroAttributeDAO()
	result, err := ad.GetMaxVersion()
	if err != nil {
		log.Logger.Error(ctx, err)
		return nil
	}
	if result != nil {
		return result
	}
	return nil
}

func GetAllHeroesFromAttr(ctx *context.Context, platform []int) ([]*model.HeroAttribute, error) {
	ad := dao.NewHeroAttributeDAO()
	return ad.Find([]string{
		"heroId",
	}, map[string]interface{}{
		"platform": platform,
	})
}

func GetAttribute(ctx *context.Context, platform int, heroID string) (*model.HeroAttribute, error) {
	ad := dao.NewHeroAttributeDAO()
	ret, err := ad.Find([]string{
		"heroId", "title", "name", "alias", "shortBio", "defense", "magic", "difficulty", "difficultyL", "attack", "attackrange", "attackdamage", "attackspeed", "attackspeedperlevel", "hp", "hpperlevel", "mp", "mpperlevel", "movespeed", "armor", "armorperlevel", "spellblock", "spellblockperlevel", "hpregen", "hpregenperlevel", "mpregen", "mpregenperlevel", "crit", "damage", "durability", "mobility", "avatar", "highlightprice", "goldPrice", "couponprice", "isWeekFree", "platform", "version", "fileTime", "ctime", "utime",
	}, map[string]interface{}{
		"platform": platform,
		"heroId":   heroID,
	})
	if err != nil {
		return nil, err
	}
	if len(ret) > 0 {
		return ret[0], nil
	}
	return nil, nil
}

// AttrData2Redis todo 未完成
func AttrData2Redis(ctx *context.Context) error {
	ad := dao.NewHeroAttributeDAO()
	attrs, err := ad.FindWithExt(nil)
	if err != nil {
		return err
	}

	data := make(map[string]dto.AttrWithExt)
	for _, attr := range attrs {
		data[attr.HeroId] = dto.AttrWithExt{}
	}

	return nil
}
