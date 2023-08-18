package jieba

import "testing"

func BenchmarkJieba(t *testing.B) {
	for i := 0; i < t.N; i++ {
		Analyzer("描述: 30攻击力 500生命值 巨像：提供相当于2%最大生命值的攻击力。 顺劈：攻击附带额外的(5 + 1.5%最大生命值)物理伤害攻击特效，并生成一道冲击波来对目标身后的敌人们造成(40 + 3%最大生命值)物理伤害。 远程英雄造成75%伤害。攻击特效伤害也可作用于建筑物。")
	}
}
