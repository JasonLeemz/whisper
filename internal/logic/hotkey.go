package logic

import (
	"whisper/pkg/context"
	"whisper/pkg/redis"

	redis2 "github.com/redis/go-redis/v9"
)

func GetHotKey(ctx *context.Context) any {
	// ZREVRANGE my_rankings 0 2 WITHSCORES
	score := redis.RDB.ZRevRangeByScoreWithScores(ctx, "SearchKey", &redis2.ZRangeBy{
		Min:    "0",
		Max:    "10",
		Offset: 0,
		Count:  0,
	})

	return score.Val()
}
