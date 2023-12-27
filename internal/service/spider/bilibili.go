package spider

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
	"whisper/internal/dto"
	"whisper/internal/service/spider/tool"
	"whisper/pkg/context"
	"whisper/pkg/http"
	"whisper/pkg/log"
	"whisper/pkg/utils"
)

type Bilibili struct {
	ts            int64
	referer       []http.Header
	commonHeaders []http.Header
}

func (b *Bilibili) Dynamic(ctx *context.Context, space, keywords string) (interface{}, error) {
	// 发送 GetForm 请求
	sk := dto.UserDynamic{}

	body, err := http.GetForm(ctx, space, b.commonHeaders...)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &sk)
	if err != nil {
		log.Logger.Error(ctx, string(body), err)
	}
	if sk.Code != 0 {
		log.Logger.Error(ctx, sk.Code, sk.Message)
		return nil, errors.New(sk.Message)
	}
	return &sk, err
}

func (b *Bilibili) DynamicDecorate(fn DynamicFunc) DynamicFunc {
	return func(ctx *context.Context, space, keywords string) (interface{}, error) {
		return fn(ctx, genDynamicParams(ctx, space, keywords), keywords)
	}
}

func genDynamicParams(ctx *context.Context, space, keywords string) string {
	urlStr := "https://api.bilibili.com/x/space/dynamic/search?keyword=" + keywords + "&pn=1&ps=30&mid=" + space
	return urlStr
}

func (b *Bilibili) SearchKeywords(ctx *context.Context, space, keywords string) (interface{}, error) {
	// 发送 GetForm 请求
	sk := dto.SearchKeywords{}

	body, err := http.GetForm(ctx, space, b.commonHeaders...)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &sk)
	if err != nil {
		log.Logger.Error(ctx, string(body), err)
	}
	if sk.Code != 0 {
		log.Logger.Error(ctx, sk.Code, sk.Message)
		return nil, errors.New(sk.Message)
	}
	return &sk, err
}

func (b *Bilibili) SearchKeywordsDecorate(fn SearchKeywordsFunc) SearchKeywordsFunc {
	return func(ctx *context.Context, space, keywords string) (interface{}, error) {
		return fn(ctx, genSearchParams(ctx, space, keywords), keywords)
	}
}

func genSearchParams(ctx *context.Context, space, keywords string) string {
	// 主页搜索接口
	//urlStr := "https://api.bilibili.com/x/space/wbi/arc/search?mid=" + space

	urlStr := "https://api.bilibili.com/x/space/wbi/arc/search?"
	wts := strconv.FormatInt(time.Now().Unix(), 10)
	key := fmt.Sprintf("dm_cover_img_str=%s&dm_img_list=%s&keyword=%s&mid=%s&order=pubdate&order_avoided=true&platform=web&pn=1&ps=30&tid=0&web_location=1550101&wts=%s",
		getDMCoverImgSTR(), "[]", keywords, space, wts)

	// salt = "72136226c6a73669787ee4fd02a74c27" # 老版本的盐
	//salt := "5a73a9f6609390773b53586cce514c2e" // 2023 0609 新
	//salt := "ea1db124af3c7062474693fa704f4ff8" // 2023 0609 新
	imgKey, subKey := tool.GetWbiKeysCached()
	salt := tool.GetMixinKey(imgKey, subKey)
	wrid := utils.Md5(key + salt)

	p := map[string]string{
		"mid":              space,
		"ps":               "30",
		"tid":              "0",
		"pn":               "1",
		"keyword":          keywords,
		"order":            "pubdate",
		"platform":         "web",
		"web_location":     "1550101",
		"order_avoided":    "true",
		"w_rid":            wrid,
		"wts":              wts,
		"dm_cover_img_str": getDMCoverImgSTR(),
		"dm_img_list":      "[]",
		"dm_img_str":       getDMImgSTR(),
	}

	urlStr += map2urlParams(p)
	fmt.Println(urlStr)

	//gurl, err := tool.SignAndGenerateURL(urlStr)
	//fmt.Println(gurl)
	//if err != nil {
	//	log.Logger.Error(ctx, err)
	//	return urlStr
	//}

	//params := convParams(p)
	//p["dm_cover_img_str"] = getDMCoverImgSTR()
	//p["dm_img_list"] = url.QueryEscape("[]")
	//p["dm_img_str"] = getDMImgSTR()
	//
	//ps := map2urlParams(p)
	//wrid := utils.Md5(ps + "ea1db124af3c7062474693fa704f4ff8")
	//for k, v := range p {
	//	urlStr += k + "=" + v + "&"
	//}
	//urlStr += "w_rid=" + wrid

	//p["dm_cover_img_str"] = getDMCoverImgSTR()
	//p["dm_img_list"] = url.QueryEscape("[]")
	//p["dm_img_str"] = getDMImgSTR()
	//gurl += "&dm_cover_img_str=" + url.QueryEscape("[]") + "&dm_cover_img_str=" + getDMCoverImgSTR() + "&dm_img_str=" + getDMImgSTR()
	//fmt.Println(gurl)
	//log.Logger.Info(ctx, gurl)

	return urlStr
}

func getDMCoverImgSTR() string {
	//return "QU5HTEUgKEFwcGxlLCBBTkdMRSBNZXRhbCBSZW5kZXJlcjogQXBwbGUgTTEsIFVuc3BlY2lmaWVkIFZlcnNpb24pR29vZ2xlIEluYy4gKEFwcGxlKQ"
	// "QU5HTEUgKEFwcGxlLCBBTkdMRSBNZXRhbCBSZW5kZXJlcjogQXBwbGUgTTEsIFVuc3BlY2lmaWVkIFZlcnNpb24pR29vZ2xlIEluYy4gKEFwcGxlKQ"
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	num := r.Intn(3-1) + 1
	//faker := fmt.Sprintf("ANGLE (Intel Inc., Intel(R) Iris(TM) Plus Graphics %d, OpenGL 4.1)Google Inc. (Intel Inc.)", num)
	faker := fmt.Sprintf("ANGLE (Apple, ANGLE Metal Renderer: Apple M%d, Unspecified Version)Google Inc. (Apple)", num)
	dm_cover_img_str := base64.StdEncoding.EncodeToString([]byte(faker))
	dm_cover_img_str = strings.Replace(dm_cover_img_str, "=", "", -1)
	fmt.Println("dm_cover_img_str", dm_cover_img_str)
	return dm_cover_img_str // 将=移除
}
func getDMImgSTR() string {
	//return "V2ViR0wgMS4wIChPcGVuR0wgRVMgMi4wIENocm9taXVtKQ"
	// "V2ViR0wgMS4wIChPcGVuR0wgRVMgMi4wIENocm9taXVtKQ"
	//r := rand.New(rand.NewSource(time.Now().UnixNano()))
	//num := r.Intn(651-350) + 350
	faker := "WebGL 1.0 (OpenGL ES 2.0 Chromium)"
	dm_cover_img_str := base64.StdEncoding.EncodeToString([]byte(faker))
	return strings.Replace(dm_cover_img_str, "=", "", -1) // 将=移除
}

func map2urlParams(mp map[string]string) string {
	params := ""
	for k, v := range mp {
		params += k + "=" + v + "&"
	}

	//return params
	return params[0 : len(params)-1]
}

func convParams(params map[string]string) map[string]string {
	imgKey, subKey := tool.GetWbiKeysCached()
	return tool.EncWbi(params, imgKey, subKey)
}
