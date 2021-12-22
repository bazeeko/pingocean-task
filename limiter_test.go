package main

import (
	"testing"
	"time"
)

func TestAtomicBool(t *testing.T) {
	testCases := []struct {
		input  bool
		output bool
	}{
		{
			input:  true,
			output: true,
		},
		{
			input:  false,
			output: false,
		},
	}

	b := AtomicBool{}
	for _, tc := range testCases {
		b.Set(tc.input)
		out := b.Get()
		if out != tc.output {
			t.Errorf("TestAtomicBool: got %v, want %v", out, tc.output)
		}
	}
}

func TestLimiter(t *testing.T) {
	var maxGoroutines = 2
	var runningGorountines int

	limiter := NewLimiter(maxGoroutines)

	for i := 0; i < 100; i++ {
		limiter.Add()
		go func(i int) {
			if limiter.RunningGorountinesCount() > runningGorountines {
				runningGorountines = limiter.RunningGorountinesCount()
			}
			time.Sleep(100 * time.Millisecond)
			limiter.Done()
		}(i)
	}

	limiter.Wait()

	if runningGorountines > maxGoroutines {
		t.Errorf("TestLimiter: number of concurrent goroutines was %d, want %d", runningGorountines, maxGoroutines)
	}
}
