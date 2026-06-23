package job_runner

import (
	"context"
	"sync"
)

type Task func(ctx context.Context)

type Pool struct {
	workers int
	tasks   chan Task
	wg      sync.WaitGroup
	ctx     context.Context
	cancel  context.CancelFunc
}

func NewPool(workers int, queueSize int) *Pool {
	ctx, cancel := context.WithCancel(context.Background())

	p := &Pool{
		workers: workers,
		tasks:   make(chan Task, queueSize),
		ctx:     ctx,
		cancel:  cancel,
	}

	p.start()
	return p
}

func (p *Pool) start() {
	for i := 0; i < p.workers; i++ {
		go p.worker()
	}
}

func (p *Pool) worker() {
	for task := range p.tasks {
		if task == nil {
			continue
		}
		task(p.ctx)
		p.wg.Done()
	}
}

func (p *Pool) Submit(task Task) {
	select {
	case <-p.ctx.Done():
		return
	default:
	}

	p.wg.Add(1)

	select {
	case p.tasks <- task:
	case <-p.ctx.Done():
		// rollback if not submitted
		p.wg.Done()
	}
}

func (p *Pool) Stop() {
	p.cancel()
	close(p.tasks)
	p.wg.Wait()
}
