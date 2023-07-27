package http

import (
	"github.com/go-resty/resty/v2"
	"whisper/pkg/context"
	"whisper/pkg/log"
)

func GetForm(ctx *context.Context, url string) ([]byte, error) {
	// Create a Resty Client
	client := resty.New()

	resp, err := client.R().
		SetHeader("Accept", "application/x-www-form-urlencoded").
		Get(url)
	//log.Logger.Debug(ctx, "url="+url, "body="+string(resp.Body()))
	return resp.Body(), err
}

func PostForm(ctx *context.Context, url string, data any) ([]byte, error) {
	// Create a Resty Client
	client := resty.New()

	resp, err := client.R().
		SetHeader("Accept", "application/x-www-form-urlencoded").
		SetBody(data).
		Get(url)
	log.Logger.Debug(ctx, "url="+url, "body="+string(resp.Body()))
	return resp.Body(), err
}
