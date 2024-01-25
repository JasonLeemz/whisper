package bloomfilter

import (
	"fmt"
	"hash/fnv"
	"math"
	"testing"
)

func TestHash(t *testing.T) {
	data := "123456789"
	index := 1

	h := fnv.New64a()
	h.Write([]byte(data))
	hashValue := h.Sum64()

	p := uint64(math.MaxUint32)

	rs := int64(index) * int64(hashValue%p)

	fmt.Println(rs)
}
