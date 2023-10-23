package equipment

import (
	"sync"
)

// 调度器
// 作用是加入command和运行command

// EquipInvoker 创建调度器
type EquipInvoker struct {
	cmds []ICommand
}

type ICommand interface {
	Exec() error // cmd应该具有exec方法
}

type IEquipInvoker interface {
	AddCommand() // 调度器应该具有将命令加入执行队列的能力
	Run()        // 调度器可以执行内部的cmd
}

func (e *EquipInvoker) AddCommand(cmds ...ICommand) {
	e.cmds = append(e.cmds, cmds...)
}

func (e *EquipInvoker) NonBlockRun() {
	wg := sync.WaitGroup{}
	wg.Add(len(e.cmds))
	for _, cmd := range e.cmds {
		go func(cmd ICommand) {
			defer func() {
				wg.Done()
			}()
			_ = cmd.Exec()
		}(cmd)
	}
	wg.Wait()
}
