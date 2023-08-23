package ik

import (
	"fmt"
	"testing"
	"whisper/pkg/context"
)

func TestIK(t *testing.T) {
	ctx := context.NewContext()
	_, err := Analyze(ctx, "<mainText><stats></stats><li><passive>灵兽同伴：</passive>召唤一只踏苔蜥来协助你打野。<li><passive>踏苔蜥的勇气：</passive>在完全成长时，你的灵兽同伴会提供一个可在击杀野怪或脱离战斗后重新生成的<shield>永久护盾</shield>。在持有这个护盾时，会获得20%韧性和减速抗性。</mainText><br>")

	fmt.Println(err)
}
