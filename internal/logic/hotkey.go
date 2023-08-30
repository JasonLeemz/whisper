package logic

import (
	redis2 "github.com/redis/go-redis/v9"
	"whisper/pkg/context"
	"whisper/pkg/redis"
)

func GetHotKey(ctx *context.Context) any {
	min := "-inf"
	max := "+inf"
	// ZREVRANGE my_rankings 0 2 WITHSCORES
	score := redis.RDB.ZRevRangeByScoreWithScores(ctx, redis.KeyHotSearchSearchBox, &redis2.ZRangeBy{
		Min:    min,
		Max:    max,
		Offset: 0,
		Count:  5,
	})

	var hotkey []string
	for _, key := range score.Val() {
		hotkey = append(hotkey, key.Member.(string))
	}

	return hotkey
}
