package service

import (
	"time"
	"whisper/internal/service/common"
	"whisper/pkg/context"
)

type ISpider interface {
	// SearchKeywords 根据关键字检索视频列表
	SearchKeywords(ctx *context.Context, url string) (interface{}, error)
}

type SpiderProductFunc func() ISpider

func CreateSpiderProduct(source int) SpiderProductFunc {
	if source == common.SourceBilibili {
		return func() ISpider {
			return &Bilibili{
				ts:            time.Now().Unix(),
				referer:       common.Referer,
				cookie:        common.Cookie,
				commonHeaders: common.CommonHeaders(),
			}
		}
	} else {
		return func() ISpider {
			// todo other source
			return nil
		}
	}
}
