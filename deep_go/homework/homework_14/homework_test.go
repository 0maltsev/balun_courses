package main

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type WorkerPool struct {
	workers     int
	taskChan    chan func()
	wg          sync.WaitGroup
	shutdownCh  chan struct{}
	isShutdown  atomic.Bool
}

func NewWorkerPool(workersNumber int) *WorkerPool {
	wp := &WorkerPool{
		workers:    workersNumber,
		taskChan:   make(chan func(), workersNumber*2),
		shutdownCh: make(chan struct{}),
	}
	for range workersNumber {
		wp.wg.Add(1)
		go wp.worker()
	}
	
	return wp
}

func (wp *WorkerPool) worker() {
	defer wp.wg.Done()
	
	for {
		select {
		case task, ok := <-wp.taskChan:
			if !ok {
				return
			}
			task()
		case <-wp.shutdownCh:
			for task := range wp.taskChan {
				task()
			}
			return
		}
	}
}

func (wp *WorkerPool) AddTask(task func()) error {
	if wp.isShutdown.Load() {
		return errors.New("worker pool is shutting down")
	}

	select {
	case wp.taskChan <- task:
		return nil
	default:
		return errors.New("task pool is full")
	}
}

func (wp *WorkerPool) Shutdown() {
	wp.isShutdown.Store(true)
	close(wp.shutdownCh)
	close(wp.taskChan)
	wp.wg.Wait()
}

func TestWorkerPool(t *testing.T) {
	var counter atomic.Int32
	task := func() {
		time.Sleep(time.Millisecond * 500)
		counter.Add(1)
	}

	pool := NewWorkerPool(2)
	_ = pool.AddTask(task)
	_ = pool.AddTask(task)
	_ = pool.AddTask(task)

	time.Sleep(time.Millisecond * 600)
	assert.Equal(t, int32(2), counter.Load())

	time.Sleep(time.Millisecond * 600)
	assert.Equal(t, int32(3), counter.Load())

	_ = pool.AddTask(task)
	_ = pool.AddTask(task)
	_ = pool.AddTask(task)
	pool.Shutdown() // wait tasks

	assert.Equal(t, int32(6), counter.Load())
}