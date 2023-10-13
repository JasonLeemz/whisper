package dao

import . "whisper/internal/model/common"

const (
	ESIndexEquipment = IndexEquipment
	ESIndexHeroes    = IndexHeroes
	ESIndexRune      = IndexRune
	ESIndexSkill     = IndexSkill
)

func CreateEsDao(t string) interface{} {
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
