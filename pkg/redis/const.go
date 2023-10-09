package redis

const (
	// KeyHotSearchSearchBox HSET
	KeyHotSearchSearchBox = "hot:search"

	// KeyHotSearchEquipBox HSET
	KeyHotSearchEquipBox = "hot:equip"

	// KeyCacheEquip SET cache:equip:map:platform:item
	KeyCacheEquip = "cache:equip:%s:%s:%s"
	// KeyCacheRune SET cache:rune:map:platform:item
	KeyCacheRune = "cache:rune:%s:%s:%s"
	// KeyCacheSkill SET cache:skill:map:platform:item
	KeyCacheSkill = "cache:skill:%s:%s:%s"

	// KeyCacheHeroEquip HSET cache:hero_equip
	KeyCacheHeroEquip = "cache:hero_equip"

	// KeyCacheEquipHeroSuit HSET cache:equip_hero:platform:item
	KeyCacheEquipHeroSuit    = "cache:equip_hero:%d:%s"
	KeyCacheEquipHeroSuitAll = "cache:equip_hero:*"
	// KeyCacheRuneHeroSuit HSET cache:rune_hero:platform:item
	KeyCacheRuneHeroSuit    = "cache:rune_hero:%d:%s"
	KeyCacheRuneHeroSuitAll = "cache:rune_hero:*"
	// KeyCacheSkillHeroSuit HSET cache:skill_hero:platform:item
	KeyCacheSkillHeroSuit    = "cache:skill_hero:%d:%s"
	KeyCacheSkillHeroSuitAll = "cache:skill_hero:*"

	// KeyCacheVersionList SET cache:version:1
	KeyCacheVersionList = "cache:version:%d"
	// KeyCacheVersionDetail SET cache:version:detail:4.2_hero
	KeyCacheVersionDetail = "cache:version:detail:%s"
)
