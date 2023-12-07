package equipment

import (
	"errors"
	"github.com/spf13/cast"
	"go.mongodb.org/mongo-driver/bson"
	"strings"
	"whisper/internal/dto"
	"whisper/internal/logic/common"
	"whisper/internal/model"
	dao "whisper/internal/model/DAO"
	"whisper/pkg/context"
	"whisper/pkg/log"
)

func FilterKeyWords(ctx *context.Context, keywords []string, platform int) ([]*model.EquipIntro, error) {
	log.Logger.Info(ctx, keywords)
	// FromMongo
	md := dao.NewMongoEquipmentDAO()

	kw := make([]bson.M, 0, len(keywords))
	for _, words := range keywords {
		in := strings.Split(words, ",")
		kw = append(kw, bson.M{
			"keywords": bson.M{
				"$in": in,
			},
		})
	}
	// 构建查询条件
	filter := bson.M{
		"platform": platform,
		"maps":     "召唤师峡谷",
		"$and":     kw,
	}

	result, err := md.Find(ctx, filter)
	return result, err
}

func GetRoadmap(ctx *context.Context, version string, platform int, equipID string, maps []string) (*dto.RespRoadmap, error) {
	resp := dto.RespRoadmap{}
	if platform == common.PlatformForLOL {
		ed := dao.NewLOLEquipmentDAO()
		roadmap, err := ed.GetRoadmap(version, equipID, maps)
		if err != nil {
			return nil, err
		}
		if len(roadmap["current"]) == 0 {
			return nil, errors.New("current equip can not find data")
		}
		resp.Current = dto.Roadmap{
			ID:        cast.ToInt(roadmap["current"][0].ItemId),
			Name:      roadmap["current"][0].Name,
			Icon:      roadmap["current"][0].IconPath,
			Maps:      roadmap["current"][0].Maps,
			Level:     "",
			Plaintext: roadmap["current"][0].Plaintext,
			Desc:      roadmap["current"][0].Description,
			Price:     cast.ToInt(roadmap["current"][0].Total),
			Sell:      cast.ToInt(roadmap["current"][0].Sell),
			Version:   roadmap["current"][0].Version,
			Platform:  common.PlatformForLOL,
		}

		for _, equip := range roadmap["into"] {
			resp.Into = append(resp.Into, dto.Roadmap{
				ID:        cast.ToInt(equip.ItemId),
				Name:      equip.Name,
				Icon:      equip.IconPath,
				Maps:      equip.Maps,
				Level:     "",
				Plaintext: equip.Plaintext,
				Desc:      equip.Description,
				Price:     cast.ToInt(equip.Total),
				Sell:      cast.ToInt(equip.Sell),
				Version:   equip.Version,
			})
		}

		fromPrice := 0
		for _, equip := range roadmap["from"] {
			price := cast.ToInt(equip.Total)
			resp.From = append(resp.From, dto.Roadmap{
				ID:        cast.ToInt(equip.ItemId),
				Name:      equip.Name,
				Icon:      equip.IconPath,
				Maps:      equip.Maps,
				Level:     "",
				Plaintext: equip.Plaintext,
				Desc:      equip.Description,
				Price:     price,
				Sell:      cast.ToInt(equip.Sell),
				Version:   equip.Version,
			})

			fromPrice += price
		}

		resp.GapPriceFrom = resp.Current.Price - fromPrice
	} else {
		ed := dao.NewLOLMEquipmentDAO()
		roadmap, err := ed.GetRoadmap(version, equipID, maps)
		if err != nil {
			return nil, err
		}
		if len(roadmap["current"]) == 0 {
			return nil, errors.New("current equip can not find data")
		}
		resp.Current = dto.Roadmap{
			ID:        cast.ToInt(roadmap["current"][0].EquipId),
			Name:      roadmap["current"][0].Name,
			Icon:      roadmap["current"][0].IconPath,
			Maps:      "",
			Level:     roadmap["current"][0].Level,
			Plaintext: "",
			Desc:      roadmap["current"][0].Description,
			Price:     cast.ToInt(&roadmap["current"][0].Price),
			Sell:      0,
			Version:   roadmap["current"][0].Version,
			Platform:  common.PlatformForLOLM,
		}

		for _, equip := range roadmap["into"] {
			resp.Into = append(resp.Into, dto.Roadmap{
				ID:        cast.ToInt(equip.EquipId),
				Name:      equip.Name,
				Icon:      equip.IconPath,
				Maps:      "",
				Level:     "",
				Plaintext: "",
				Desc:      equip.Description,
				Price:     cast.ToInt(equip.Price),
				Sell:      0,
				Version:   equip.Version,
			})
		}

		fromPrice := 0
		for _, equip := range roadmap["from"] {
			price := cast.ToInt(equip.Price)
			resp.From = append(resp.From, dto.Roadmap{
				ID:        cast.ToInt(equip.EquipId),
				Name:      equip.Name,
				Icon:      equip.IconPath,
				Maps:      "",
				Level:     "",
				Plaintext: "",
				Desc:      equip.Description,
				Price:     price,
				Sell:      0,
				Version:   equip.Version,
			})
			fromPrice += price
		}

		resp.GapPriceFrom = resp.Current.Price - fromPrice
	}

	// ===================

	return &resp, nil
}
