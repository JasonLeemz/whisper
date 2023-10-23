package logic

import (
	"errors"
	"github.com/spf13/cast"
	"whisper/internal/dto"
	"whisper/internal/logic/common"
	dao "whisper/internal/model/DAO"
	"whisper/pkg/context"
)

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
