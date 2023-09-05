package redis

const (
	KeyHotSearchSearchBox = "hot:search"
	KeyHotSearchEquipBox  = "hot:equip"

	// cache:equip:map:platform:item
	KeyCacheEquip = "cache:equip:%s:%s:%s"
	// cache:hero_equip
	KeyCacheHeroEquip = "cache:hero_equip"
)
