package spider

import (
	"whisper/pkg/http"
)

var (
	faker []http.Header = []http.Header{
		http.Header{
			Key:   "authority",
			Value: "api.bilibili.com",
		},
		http.Header{
			Key:   "cache-control",
			Value: "max-age=0",
		},
		http.Header{
			Key:   "sec-ch-ua",
			Value: "\" Not;A Brand\";v=\"99\", \"Google Chrome\";v=\"91\", \"Chromium\";v=\"91\"",
		},
		http.Header{
			Key:   "sec-ch-ua-mobile",
			Value: "?0",
		},
		http.Header{
			Key:   "upgrade-insecure-requests",
			Value: "1",
		},
		http.Header{
			Key:   "user-agent",
			Value: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36",
		},
		http.Header{
			Key:   "accept",
			Value: "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
		},
		http.Header{
			Key:   "sec-fetch-site",
			Value: "none",
		},
		http.Header{
			Key:   "sec-fetch-mode",
			Value: "navigate",
		},
		http.Header{
			Key:   "sec-fetch-user",
			Value: "?1",
		},
		http.Header{
			Key:   "sec-fetch-dest",
			Value: "document",
		},
		http.Header{
			Key:   "accept-language",
			Value: "zh-CN,zh;q=0.9",
		},
	}

	referer []http.Header = []http.Header{
		http.Header{
			Key:   "Referer",
			Value: "https://m.bilibili.com/space/",
		},
	}

	commonHeader []http.Header = []http.Header{
		http.Header{
			Key:   "User-Agent",
			Value: "Mozilla/5.0",
		},
		//http.Header{
		//	Key:   "Content-Type",
		//	Value: "application/json; charset=utf-8",
		//},
		//http.Header{
		//	Key:   "Accept",
		//	Value: "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
		//},
		http.Header{
			Key:   "accept-language",
			Value: "en,zh-CN;q=0.9,zh;q=0.8",
		},
	}
)

func commonHeaders() []http.Header {
	var h []http.Header
	//h = append(h, faker...)
	//h = append(h, http.Header{
	//	Key: "Cookie",
	//	//Value: config.SpiderCfg.Bilibili.Cookie,
	//	Value: "bsource=search_baidu; innersign=1",
	//})
	h = append(h, commonHeader...)
	return h
}
