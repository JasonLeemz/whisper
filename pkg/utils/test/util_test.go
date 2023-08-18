package test

import (
	"fmt"
	"testing"
	"whisper/pkg/utils"
)

func TestCompareVersion(t *testing.T) {
	res := utils.CompareVersion("1.2a", "1.2b")

	fmt.Println("result:", res)

}
