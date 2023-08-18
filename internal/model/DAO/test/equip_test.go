package test

import (
	"fmt"
	"testing"
	run "whisper/init"
	dao2 "whisper/internal/model/DAO"
)

func TestFind(t *testing.T) {
	run.Init()
	dao := dao2.NewLOLEquipmentDAO()
	result, err := dao.Find([]string{
		"max(version) as version",
	}, map[string]interface{}{
		"id": 0,
	})

	l := len(result)

	fmt.Println(l, result, err)
}
