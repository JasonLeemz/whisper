package spider

import (
	"whisper/pkg/config"
	"whisper/pkg/http"
)

var (
	referer []http.Header = []http.Header{
		http.Header{
			Key:   "Referer",
			Value: "https://space.bilibili.com",
		},
	}

	commonHeader []http.Header = []http.Header{
		http.Header{
			Key:   "User-Agent",
			Value: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36",
		},
		http.Header{
			Key:   "Content-Type",
			Value: "application/json; charset=utf-8",
		},
		http.Header{
			Key:   "Accept",
			Value: "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
		},
	}
)

func commonHeaders() []http.Header {
	var h []http.Header
	h = append(h, referer...)
	h = append(h, http.Header{
		Key:   "Cookie",
		Value: config.SpiderCfg.Bilibili.Cookie,
	})
	h = append(h, commonHeader...)
	return h
}
