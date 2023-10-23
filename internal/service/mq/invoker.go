package mq

import (
	"sync"
	dao "whisper/internal/model/DAO"
)

type ICommand interface {
	Exec(args ...interface{}) error
}

type Invoker struct {
	cmds []ICommand
}

func (i *Invoker) BlockRun(args ...interface{}) {
	for _, cmd := range i.cmds {
		if err := cmd.Exec(args...); err != nil {
			break
		}
	}
}

func (i *Invoker) NonBlockRun(args ...interface{}) {
	wg := sync.WaitGroup{}
	wg.Add(len(i.cmds))
	for _, cmd := range i.cmds {
		go func(cmd ICommand) {
			defer func() {
				wg.Done()
			}()
			_ = cmd.Exec(args...)
		}(cmd)
	}
	wg.Wait()
}

func (i *Invoker) AddCommand(cmds ...ICommand) {
	i.cmds = append(i.cmds, cmds...)
}

// ProduceMessage ...
func ProduceMessage(exchange, routingKey string, message []byte) {
	mqd := dao.NewMQDao()
	mqd.ProduceMessage(exchange, routingKey, message)
}
