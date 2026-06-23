package job_runner

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestScheduleAfter(t *testing.T) {
	s := NewScheduler()
	defer s.Stop()

	var executed atomic.Int32

	s.ScheduleAfter(50*time.Millisecond, func() {
		executed.Store(1)
	})

	time.Sleep(100 * time.Millisecond)

	if executed.Load() != 1 {
		t.Fatalf("job not executed")
	}
}

func TestScheduleAt(t *testing.T) {
	s := NewScheduler()
	defer s.Stop()

	var executed int32

	s.ScheduleAt(time.Now().Add(50*time.Millisecond), func() {
		atomic.StoreInt32(&executed, 1)
	})

	time.Sleep(100 * time.Millisecond)

	if atomic.LoadInt32(&executed) != 1 {
		t.Fatalf("job not executed")
	}
}

func TestScheduleEvery(t *testing.T) {
	s := NewScheduler()
	defer s.Stop()

	var count atomic.Int32

	s.ScheduleEvery(30*time.Millisecond, func() {
		count.Add(1)
	})

	time.Sleep(120 * time.Millisecond)

	c := count.Load()

	if c < 3 {
		t.Fatalf("expected at least 3 executions, got %d", c)
	}
}

func TestScheduler_Stop(t *testing.T) {
	s := NewScheduler()

	var count atomic.Int32

	s.ScheduleEvery(20*time.Millisecond, func() {
		count.Add(1)
	})

	time.Sleep(50 * time.Millisecond)
	s.Stop()

	c1 := count.Load()
	time.Sleep(50 * time.Millisecond)
	c2 := count.Load()

	if c2 > c1 {
		t.Fatalf("jobs should stop after Stop()")
	}
}
