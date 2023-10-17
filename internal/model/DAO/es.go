package dao

import (
	. "whisper/internal/model/common"
	"whisper/pkg/context"
)

const (
	ESIndexEquipment = IndexEquipment
	ESIndexHeroes    = IndexHeroes
	ESIndexRune      = IndexRune
	ESIndexSkill     = IndexSkill
)

type EsDaoFunc func() ESIndex

func CreateEsDao(t string) EsDaoFunc {
	switch t {
	case ESIndexEquipment:
		return NewESEquipmentDAO()
	case ESIndexHeroes:
		return NewESHeroesDAO()
	case ESIndexRune:
		return NewESRuneDAO()
	case ESIndexSkill:
		return NewESSkillDAO()

	}
	return nil
}

type ESIndex interface {
	CreateIndex(ctx *context.Context) error
	DeleteIndex(ctx *context.Context) error
	Data2ES(ctx *context.Context, data interface{}) error
	Find(ctx *context.Context, cond *QueryCond) ([]map[string]interface{}, error)
}
