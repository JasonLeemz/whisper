package service

import (
	"fmt"
	"math"
	"testing"
	"time"
)

func TestTs(t *testing.T) {
	t1 := time.Now().Unix()
	t2 := time.Now().Unix() / 600
	t3 := math.Floor(float64(time.Now().Unix() / 600))
	t4 := fmt.Sprintf("%d", time.Now().Unix()/600)
	fmt.Println(t1)
	fmt.Println(t2)
	fmt.Println(t3)
	fmt.Println(t4)
}
