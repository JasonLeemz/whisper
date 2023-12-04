package spider

import (
	"fmt"
	"sync"
	"sync/atomic"
	"whisper/internal/dto"
	"whisper/internal/model"
	dao2 "whisper/internal/model/DAO"
	service "whisper/internal/service/spider"
	"whisper/pkg/context"
	"whisper/pkg/log"
)

func (e *Spider) BilibiliGrab() error {
	// 获取用户空间地址
	dao := dao2.NewAuthorSpaceDAO()
	authors, err := dao.Find(nil, nil)
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
	t := newTask(int32(len(authors)*len(heroes)), &sync.WaitGroup{}, make(chan struct{}, 2), make(chan struct{}, 5))
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
				url := author.Space // url不打日志了，在service层打过了
				if hero.Platform == 0 {
					url = fmt.Sprintf(url, hero.Title)
				} else {
					url = fmt.Sprintf(url, hero.Name)
				}

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

					// 写入数据
					recodeBilibiliData(ctx, bdata)
				}()
			}
		}(author)
	}
	t.wg.Wait()

	log.Logger.Info(ctx, fmt.Sprintf("处理了: %d 个任务", t.done))
	log.Logger.Info(ctx, fmt.Sprintf("提前结束,执行出错: %d 个任务", t.fail))
	log.Logger.Info(ctx, fmt.Sprintf("成功执行了: %d 个任务", t.success))
	log.Logger.Info(ctx, fmt.Sprintf("剩余: %d 个任务待处理", t.total-t.done))
}

func recodeBilibiliData(ctx *context.Context, data *dto.SearchKeywords) {

}
