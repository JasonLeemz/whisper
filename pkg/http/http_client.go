package http

import (
	"github.com/go-resty/resty/v2"
	"whisper/pkg/context"
	"whisper/pkg/log"
)

type Header struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func GetForm(ctx *context.Context, url string, header ...Header) ([]byte, error) {
	// Create a Resty Client
	client := resty.New()
	cli := client.R()
	for _, h := range header {
		cli.SetHeader(h.Key, h.Value)
	}
	if len(header) == 0 {
		cli.SetHeader("Accept", "application/x-www-form-urlencoded")
	}
	resp, err := cli.Get(url)
	log.Logger.Debug(ctx, "url="+url, "body="+string(resp.Body()))
	return resp.Body(), err
}

func PostForm(ctx *context.Context, url string, data any, header ...Header) ([]byte, error) {
	// Create a Resty Client
	client := resty.New()

	cli := client.R()
	for _, h := range header {
		cli.SetHeader(h.Key, h.Value)
	}
	if len(header) == 0 {
		cli.SetHeader("Accept", "application/x-www-form-urlencoded")
	}

	resp, err := cli.
		SetBody(data).
		Get(url)
	log.Logger.Debug(ctx, "url="+url, "body="+string(resp.Body()))
	return resp.Body(), err
}
