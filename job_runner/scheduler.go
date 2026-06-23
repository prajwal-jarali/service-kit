package job_runner

import (
	"context"
	"time"
)

type Job func()

type Scheduler struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func NewScheduler() *Scheduler {
	ctx, cancel := context.WithCancel(context.Background())
	return &Scheduler{ctx: ctx, cancel: cancel}
}

func (s *Scheduler) Stop() {
	s.cancel()
}

// Run once at specific time
func (s *Scheduler) ScheduleAt(t time.Time, job Job) {
	go func() {
		delay := max(time.Until(t), 0)

		select {
		case <-time.After(delay):
			job()
		case <-s.ctx.Done():
		}
	}()
}

// Run once after delay
func (s *Scheduler) ScheduleAfter(delay time.Duration, job Job) {
	go func() {
		select {
		case <-time.After(delay):
			job()
		case <-s.ctx.Done():
		}
	}()
}

// Fixed delay loop (drift-resistant)
func (s *Scheduler) ScheduleEvery(interval time.Duration, job Job) {
	go func() {
		start := time.Now()

		for i := 0; ; i++ {
			next := start.Add(time.Duration(i+1) * interval)
			sleep := max(time.Until(next), 0)

			select {
			case <-time.After(sleep):
				job()
			case <-s.ctx.Done():
				return
			}
		}
	}()
}
