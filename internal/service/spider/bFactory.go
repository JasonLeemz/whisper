package spider

import (
	"time"
	"whisper/internal/model"
	"whisper/internal/service/common"
	"whisper/pkg/context"
)

type SearchKeywordsFunc func(ctx *context.Context, url, keywords string) (interface{}, error)
type DynamicFunc func(ctx *context.Context, space, keywords string) (interface{}, error)

type ISpider interface {
	SearchKeywordsDecorate(fn SearchKeywordsFunc) SearchKeywordsFunc
	// SearchKeywords 根据关键字检索视频列表
	SearchKeywords(ctx *context.Context, space, keywords string) (interface{}, error)

	DynamicDecorate(fn DynamicFunc) DynamicFunc
	// Dynamic 用户动态
	Dynamic(ctx *context.Context, space, keywords string) (interface{}, error)
}

type SpiderProductFunc func() ISpider

func CreateSpiderProduct(author *model.AuthorSpace) SpiderProductFunc {
	if author.Source == common.SourceBilibili {
		return func() ISpider {
			referer[0].Value += author.Space
			return &Bilibili{
				ts:            time.Now().Unix(),
				referer:       referer,
				commonHeaders: commonHeaders(),
			}
		}
	} else {
		return func() ISpider {
			// todo other source
			return nil
		}
	}
}
