package spider

import "sync"

type task struct {
	total   int32
	success int32
	fail    int32
	done    int32
	wg      *sync.WaitGroup
	ch1     chan struct{}
	ch2     chan struct{}
}

func newTask(total int32, wg *sync.WaitGroup, ch1, ch2 chan struct{}) *task {
	return &task{total: total, wg: wg, ch1: ch1, ch2: ch2}
}
