package spider

import (
	"fmt"
	"math/rand"
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
	// 获取所有英雄
	adao := dao2.NewHeroAttributeDAO()
	heroes, err := adao.QueryAllHeroes(nil)
	if err != nil {
		log.Logger.Error(e.ctx, err)
		return err
	}

	// 缓存视频列表
	fetchData(e.ctx, authors, heroes)
	return nil
}

func fetchData(ctx *context.Context, authors []*model.AuthorSpace, heroes []*model.HeroAttribute) {
	t := newTask(int32(len(authors)*len(heroes)), &sync.WaitGroup{}, make(chan struct{}, 5), make(chan struct{}, 10))
	for i, author := range authors {
		log.Logger.Info(ctx, ">>>>>>>>>>开始处理 hero:<<<<<<<<<<<", i, "/", author.Name, "/", author.Space)
		t.wg.Add(1)
		t.ch1 <- struct{}{}

		go func(author *model.AuthorSpace) {
			defer func() {
				t.wg.Done()
				<-t.ch1
			}()
			for _, hero := range heroes {
				if hero.Platform != author.Platform {
					atomic.AddInt32(&t.done, 1) // 提前结束要把完成数+1
					// 视频博主的游戏平台要和该英雄的所属平台一致
					continue
				}

				url := author.Space // url不打日志了，在service层打过了
				name := ""
				if hero.Platform == 0 {
					name = hero.Title
				} else {
					name = hero.Name
				}
				url = fmt.Sprintf(url, name)

				t.wg.Add(1)
				t.ch2 <- struct{}{}

				go func() {
					defer func() {
						t.wg.Done()
						<-t.ch2
						atomic.AddInt32(&t.done, 1)
					}()

					data, err := service.CreateSpiderProduct(author.Source)().SearchKeywords(ctx, url)
					if err != nil {
						atomic.AddInt32(&t.fail, 1)
						log.Logger.Error(ctx, "SearchKeywords", err)
						return
					}

					bdata, ok := data.(*dto.SearchKeywords)
					if !ok {
						atomic.AddInt32(&t.fail, 1)
						log.Logger.Error(ctx, "data.(*dto.SearchKeywords) assert fail")
						return
					}

					if bdata.Code != 0 {
						log.Logger.Error(ctx, bdata.Message)
						return
					}

					// 写入数据
					err = recodeBilibiliData(ctx, bdata, author.Platform, name)
					if err != nil {
						atomic.AddInt32(&t.fail, 1)
					} else {
						atomic.AddInt32(&t.success, 1)
					}

					// 生成一个1到3秒之间的随机时间间隔
					rand.Seed(time.Now().UnixNano())
					duration := time.Duration(rand.Intn(3)+1) * time.Second
					time.Sleep(duration)
				}()

			}
		}(author)
	}
	t.wg.Wait()

	log.Logger.Info(ctx, fmt.Sprintf("共有: %d 个任务", t.total))
	log.Logger.Info(ctx, fmt.Sprintf("处理了: %d 个任务", t.done))
	log.Logger.Info(ctx, fmt.Sprintf("提前结束,执行出错: %d 个任务", t.fail))
	log.Logger.Info(ctx, fmt.Sprintf("成功执行了: %d 个任务", t.success))
	log.Logger.Info(ctx, fmt.Sprintf("剩余: %d 个任务待处理", t.total-t.done))
}

func recodeBilibiliData(ctx *context.Context, data *dto.SearchKeywords, platform int, hero string) error {
	// https://m.bilibili.com/video/%s => https://m.bilibili.com/video/BV1Zh411C7XS
	baseUrl := "https://m.bilibili.com/video/%s"
	dao := dao2.NewGameStrategyDAO()
	for _, d := range data.Data.List.Vlist {
		strategy := &model.GameStrategy{
			Platform:   platform,
			Source:     common.SourceBilibili,
			Author:     d.Author,
			LinkUrl:    fmt.Sprintf(baseUrl, d.Bvid),
			MainImage:  d.Pic,
			PublicDate: time.Unix(d.Created, 0), // 1613243324
			Title:      d.Title,
			Subtitle:   d.Subtitle,
			Status:     0,
			Bvid:       d.Bvid,
			Played:     d.Play,
			Hero:       hero,
		}
		rows, err := dao.InsertORIgnore(strategy)
		if err != nil {
			log.Logger.Error(ctx, err)
			return err
		}
		if rows == 0 {
			log.Logger.Info(ctx, hero, "未更新")
			return err
		} else {
			log.Logger.Info(ctx, d.Author, hero, "ok")
		}
	}
	return nil
}
