package common

import "whisper/pkg/http"

var (
	Referer []http.Header = []http.Header{
		http.Header{
			Key:   "Referer",
			Value: "https://101.qq.com/",
		},
	}

	UA []http.Header = []http.Header{
		http.Header{
			Key:   "user-agent",
			Value: "QTL/9.6.0 (iPhone; IOS 17.0.3; Scale/3.00)",
		},
	}
)

func CommonHeaders() []http.Header {
	var h []http.Header
	h = append(h, Referer...)
	h = append(h, UA...)
	return h
}
