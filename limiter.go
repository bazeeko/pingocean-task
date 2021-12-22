package main

import "sync/atomic"

type AtomicBool struct {
	value int32
}

func (b *AtomicBool) Set(v bool) {
	if v {
		atomic.StoreInt32(&b.value, 1)
	} else {
		atomic.StoreInt32(&b.value, 0)
	}
}

func (b *AtomicBool) Get() bool {
	return atomic.LoadInt32(&b.value) == 1
}

type Limiter struct {
	runningGorountinesCount int32
	pool                    chan struct{}
	goroutineDone           chan struct{}
	allDone                 chan struct{}
	closed                  AtomicBool
}

func NewLimiter(limit int) *Limiter {
	l := &Limiter{
		runningGorountinesCount: 0,
		pool:                    make(chan struct{}, limit),
		goroutineDone:           make(chan struct{}),
		allDone:                 make(chan struct{}),
		closed:                  AtomicBool{},
	}

	go l.run()

	return l
}

func (l *Limiter) Add() {
	l.pool <- struct{}{}

	atomic.AddInt32(&l.runningGorountinesCount, 1)
}

func (l *Limiter) Done() {
	atomic.AddInt32(&l.runningGorountinesCount, -1)

	l.goroutineDone <- struct{}{}
}

func (l *Limiter) Wait() {
	l.closed.Set(true)

	<-l.allDone

	close(l.goroutineDone)
	close(l.pool)
	close(l.allDone)
}

func (l *Limiter) RunningGorountinesCount() int {
	count := atomic.LoadInt32(&l.runningGorountinesCount)
	return int(count)
}

func (l *Limiter) run() {
LOOP:
	for {
		<-l.goroutineDone
		<-l.pool

		if l.closed.Get() && atomic.LoadInt32(&l.runningGorountinesCount) == 0 {
			break LOOP
		}
	}

	l.allDone <- struct{}{}
}
