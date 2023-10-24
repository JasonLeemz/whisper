package command

import (
	"sync"
)

// 调度器
// 作用是加入command和运行command

// Invoker 创建调度器
type Invoker struct {
	cmds []ICommand
}

type ICommand interface {
	ExecEXT() (interface{}, error)
	ExecE() error
	Exec()
}

type IInvoker interface {
	AddCommand(cmds ...ICommand) // 调度器应该具有将命令加入执行队列的能力
	BlockRun() error             // 调度器可以执行内部的cmd(阻塞执行)
	NonBlockRun()                // 调度器可以执行内部的cmd(非阻塞执行)
}

func (e *Invoker) AddCommand(cmds ...ICommand) {
	e.cmds = append(e.cmds, cmds...)
}

func (e *Invoker) BlockRun() error {
	for _, cmd := range e.cmds {
		if err := cmd.ExecE(); err != nil {
			return err
		}
	}
	return nil
}

func (e *Invoker) NonBlockRun() {
	wg := sync.WaitGroup{}
	wg.Add(len(e.cmds))
	for _, cmd := range e.cmds {
		go func(cmd ICommand) {
			defer func() {
				wg.Done()
			}()
			_ = cmd.ExecE()
		}(cmd)
	}
	wg.Wait()
}
