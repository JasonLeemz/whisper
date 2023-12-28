package spider

import (
	context2 "context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
	"whisper/internal/dto"
	"whisper/internal/logic/common"
	"whisper/internal/model"
	dao2 "whisper/internal/model/DAO"
	service "whisper/internal/service/spider"
	"whisper/pkg/context"
	"whisper/pkg/log"
)

func (e *Spider) BilibiliGrab() error {
	// 获取用户空间地址
	dao := dao2.NewAuthorSpaceDAO()
	authors, err := dao.Find([]string{
		"name",
		"space",
		"video_base_url",
		"source",
		"platform",
		"status",
	}, map[string]interface{}{
		"status": 0,
	})
	if err != nil {
		log.Logger.Error(e.ctx, err)
		return err
	}

	// 根据关键字从用户空间URl检索视频列表
	// 获取所有需要抓取的数据
	//sp := make([]*searchParams, 0)
	sp := getAllHeroes(e.ctx)
	sp = append(sp, getAllRunes(e.ctx)...)
	sp = append(sp, getAllEquips(e.ctx)...)

	// 缓存视频列表
	fetchData(e.ctx, authors, sp)
	return nil
}

func fetchData(ctx *context.Context, authors []*model.AuthorSpace, sp []*searchParams) {
	cancelCtx, cancelFunc := context2.WithCancel(ctx)
	defer cancelFunc()

	t := newTask(int32(len(authors)*len(sp)), &sync.WaitGroup{}, make(chan struct{}, 5), make(chan struct{}, 50))

	select {
	case <-cancelCtx.Done():
		log.Logger.Error(ctx, "fail times too much")
		return
	default:
		for _, author := range authors {
			t.wg.Add(1)
			t.ch1 <- struct{}{}

			go func(author *model.AuthorSpace) {
				defer func() {
					t.wg.Done()
					<-t.ch1
				}()
				spider := service.CreateSpiderProduct(author)()

				for _, params := range sp {
					if params.platform != author.Platform {
						atomic.AddInt32(&t.success, 1)
						// 视频博主的游戏平台要和该英雄的所属平台一致
						continue
					}

					t.wg.Add(1)
					t.ch2 <- struct{}{}

					go func(params *searchParams) {
						defer func() {
							t.wg.Done()
							<-t.ch2
						}()

						if t.fail > 25 {
							log.Logger.Error(ctx, "fail times:", t.fail)
							cancelFunc()
							return
						}

						// url 需要验签参数
						data, err := spider.DynamicDecorate(spider.Dynamic)(ctx, author.Space, params.keywords)
						if err != nil {
							atomic.AddInt32(&t.fail, 1)
							log.Logger.Error(ctx, "SearchKeywords", params.keywords, err)
							return
						}

						bdata, ok := data.(*dto.UserDynamic)
						if !ok {
							atomic.AddInt32(&t.fail, 1)
							log.Logger.Error(ctx, "data.(*dto.UserDynamic) assert fail")
							return
						}

						if bdata.Code != 0 {
							atomic.AddInt32(&t.fail, 1)
							log.Logger.Error(ctx, bdata.Message)
							return
						}

						// 写入数据
						err = recodeBilibiliData(ctx, bdata, author.Platform, params.desc)
						if err != nil {
							atomic.AddInt32(&t.fail, 1)
						} else {
							atomic.AddInt32(&t.success, 1)
						}
						// 生成一个1到3秒之间的随机时间间隔
						//rand.Seed(time.Now().UnixNano())
						//duration := time.Duration(rand.Intn(3)+1) * time.Second
						//time.Sleep(duration)
					}(params)

				}
			}(author)
		}
	}

	t.wg.Wait()

	log.Logger.Info(ctx, fmt.Sprintf("BilibiliGrab Done"))
	log.Logger.Info(ctx, fmt.Sprintf("共有: %d 个任务", t.total))
	log.Logger.Info(ctx, fmt.Sprintf("提前结束,执行出错: %d 个任务", t.fail))
	log.Logger.Info(ctx, fmt.Sprintf("成功执行了: %d 个任务", t.success))
	log.Logger.Info(ctx, fmt.Sprintf("剩余: %d 个任务待处理", t.total-t.success-t.fail))
}

func recodeBilibiliData(ctx *context.Context, data *dto.UserDynamic, platform int, hero string) error {
	// https://m.bilibili.com/video/%s => https://m.bilibili.com/video/BV1Zh411C7XS
	baseUrl := "https://m.bilibili.com/video/%s"
	dao := dao2.NewGameStrategyDAO()
	for _, d := range data.Data.Cards {
		card := dto.DynamicCardsCard{}
		e := json.Unmarshal([]byte(d.Card), &card)
		if e != nil {
			log.Logger.Error(ctx, e, "Card:", d.Card)
			continue
		}

		min := ""
		if card.Duration/60 < 10 {
			min = "0" + strconv.Itoa(card.Duration/60)
		} else {
			min = strconv.Itoa(card.Duration / 60)
		}
		length := fmt.Sprintf("%s:%d", min, card.Duration%60)
		strategy := &model.GameStrategy{
			Platform:   platform,
			Source:     common.SourceBilibili,
			Author:     d.Desc.UserProfile.Info.Uname,
			LinkUrl:    fmt.Sprintf(baseUrl, d.Desc.Bvid),
			MainImage:  card.Pic + "@406w_254h_1e_1c.jpg", // https://i2.hdslb.com/bfs/archive/f3122fd3afa2dc3becb4b045055150af410cbd20.jpg@406w_254h_1e_1c.jpg
			PublicDate: time.Unix(d.Desc.Timestamp, 0),    // 1613243324
			Title:      card.Title,
			Subtitle:   "",
			Status:     0,
			Bvid:       d.Desc.Bvid,
			Played:     card.Stat.View,
			Keywords:   hero,
			Length:     length,
		}
		err := dao.InsertORIgnore(strategy)
		if err != nil {
			log.Logger.Error(ctx, err)
			return err
		}
		log.Logger.Info(ctx, d.Desc.UserProfile.Info.Uname, hero, d.Desc.Bvid, "ok")
	}
	return nil
}

type searchParams struct {
	desc     string
	keywords string
	platform int
}

func getAllHeroes(ctx *context.Context) []*searchParams {
	// 获取所有英雄
	adao := dao2.NewHeroAttributeDAO()
	heroes, err := adao.QueryAllHeroes(nil)
	if err != nil {
		log.Logger.Error(ctx, err)
		return nil
	}
	sp := make([]*searchParams, 0, len(heroes)*2)

	for _, hero := range heroes {
		p := &searchParams{}
		if hero.Platform == 0 {
			p.keywords = hero.Title
			p.desc = hero.Title
			p.platform = hero.Platform

			sp = append(sp, p)
			p.keywords = hero.Name
			sp = append(sp, p)
		} else {
			p.keywords = hero.Name
			p.desc = hero.Name
			p.platform = hero.Platform

			sp = append(sp, p)
			p.keywords = hero.Title
			sp = append(sp, p)
		}
	}

	return sp
}

func getAllRunes(ctx *context.Context) []*searchParams {
	sp := make([]*searchParams, 0, 0)
	data1, err := dao2.NewLOLRuneDAO().Find([]string{"*"}, map[string]interface{}{
		"status": 0,
	})
	if err != nil {
		log.Logger.Error(ctx, err)
		return sp
	}
	for _, data := range data1 {
		sp = append(sp, &searchParams{
			desc:     data.Name,
			keywords: data.Name,
			platform: common.PlatformForLOL,
		})
	}
	data2, err := dao2.NewLOLMRuneDAO().Find([]string{"*"}, map[string]interface{}{
		"status": 0,
	})
	if err != nil {
		log.Logger.Error(ctx, err)
		return sp
	}
	for _, data := range data2 {
		sp = append(sp, &searchParams{
			desc:     data.Name,
			keywords: data.Name,
			platform: common.PlatformForLOLM,
		})
	}

	return sp
}

func getAllEquips(ctx *context.Context) []*searchParams {
	sp := make([]*searchParams, 0, 0)
	ed := dao2.NewLOLEquipmentDAO()
	version, _ := ed.GetLOLEquipmentMaxVersion()
	data1, err := ed.GetLOLEquipment(version.Version)
	if err != nil {
		log.Logger.Error(ctx, err)
		return sp
	}
	for _, data := range data1 {
		sp = append(sp, &searchParams{
			desc:     data.Name,
			keywords: data.Name,
			platform: common.PlatformForLOL,
		})
	}

	ed2 := dao2.NewLOLMEquipmentDAO()
	version2, _ := ed2.GetLOLMEquipmentMaxVersion()
	data2, err := ed2.GetLOLMEquipment(version2.Version)
	if err != nil {
		log.Logger.Error(ctx, err)
		return sp
	}
	for _, data := range data2 {
		sp = append(sp, &searchParams{
			desc:     data.Name,
			keywords: data.Name,
			platform: common.PlatformForLOLM,
		})
	}

	return sp
}
