package spider

import (
	"encoding/json"
	"fmt"
	"testing"
	run "whisper/init"
	"whisper/internal/model"
	"whisper/pkg/context"
	"whisper/pkg/utils"
)

func TestSpider(t *testing.T) {
	ctx := context.NewContext()
	run.Init()

	author := &model.AuthorSpace{
		Space:    "424730226",
		Source:   0,
		Platform: 1,
	}
	s := CreateSpiderProduct(author)()
	for i := 0; i < 1; i++ {
		data, err := s.DynamicDecorate(s.Dynamic)(ctx, author.Space, "婕拉")
		p, _ := json.Marshal(data)
		fmt.Println(string(p), err)
	}
}

func TestWBI(t *testing.T) {
	r := utils.Md5("limingze")
	fmt.Println(r)
	//tool.Test()
}
