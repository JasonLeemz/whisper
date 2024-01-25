package bloomfilter

import (
	"context"
	"github.com/redis/go-redis/v9"
	"hash/fnv"
	"math"
)

type EquipBloomFilter struct {
	redisClient *redis.Client
	hashFuncs   uint
}

func (bf *EquipBloomFilter) Add(data string, bitSetKey string) {
	for i := uint(0); i < bf.hashFuncs; i++ {
		hashValue := bf.hash(data, i)
		bf.redisClient.SetBit(context.Background(), bitSetKey, hashValue, 1)
	}
}

// Contains 判断元素是否存在于布隆过滤器
func (bf *EquipBloomFilter) Contains(data string, bitSetKey string) bool {
	for i := uint(0); i < bf.hashFuncs; i++ {
		hashValue := bf.hash(data, i)
		bit := bf.redisClient.GetBit(context.Background(), bitSetKey, hashValue)
		if bit.Val() == 0 {
			return false
		}
	}
	return true
}

// hash 计算哈希值
func (bf *EquipBloomFilter) hash(data string, index uint) int64 {
	h := fnv.New64a()
	h.Write([]byte(data))
	hashValue := h.Sum64()
	return int64(index) * int64(hashValue%uint64(math.MaxUint32))
}

// NewEquipBloomFilter 创建一个新的布隆过滤器
func NewEquipBloomFilter(redisClient *redis.Client, hashFuncs uint) *EquipBloomFilter {
	return &EquipBloomFilter{
		redisClient: redisClient,
		hashFuncs:   hashFuncs,
	}
}
