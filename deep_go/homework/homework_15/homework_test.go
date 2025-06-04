package main

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type Group struct {
	wg     sync.WaitGroup
	cancel context.CancelFunc
	errCh  chan error
}

func NewErrGroup(ctx context.Context) (*Group, context.Context) {
	childCtx, cancel := context.WithCancel(ctx)
	g := &Group{
		cancel: cancel,
		errCh:  make(chan error, 1),
	}
	return g, childCtx
}

func (g *Group) Go(action func() error) {
	if action == nil {
		return
	}
	
	g.wg.Add(1)
	
	go func() {
		defer g.wg.Done()
		
		if err := action(); err != nil {
			// Пытаемся отправить ошибку в канал (неблокирующая операция)
			select {
			case g.errCh <- err:
				// Ошибка отправлена, отменяем контекст
				g.cancel()
			default:
				// Канал уже заполнен (уже есть ошибка), игнорируем
			}
		}
	}()
}

func (g *Group) Wait() error {
	done := make(chan struct{})
	go func() {
		g.wg.Wait()
		close(done)
	}()
	select {
	case err := <-g.errCh:
		g.cancel()
		g.wg.Wait()
		return err
	case <-done:
		g.cancel()
		return nil
	}
}

func TestErrGroupWithoutError(t *testing.T) {
	var counter atomic.Int32
	group, _ := NewErrGroup(context.Background())

	for i := 0; i < 5; i++ {
		group.Go(func() error {
			time.Sleep(time.Second)
			counter.Add(1)
			return nil
		})
	}

	err := group.Wait()
	assert.Equal(t, int32(5), counter.Load())
	assert.NoError(t, err)
}

func TestErrGroupWithError(t *testing.T) {
	var counter atomic.Int32
	group, ctx := NewErrGroup(context.Background())

	for i := 0; i < 5; i++ {
		group.Go(func() error {
			timer := time.NewTimer(time.Second)
			defer timer.Stop()

			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-timer.C:
				counter.Add(1)
				return nil
			}
		})
	}

	group.Go(func() error {
		return errors.New("error")
	})

	err := group.Wait()
	assert.Equal(t, int32(0), counter.Load())
	assert.Error(t, err)
}