package closer

import (
	"context"
	"errors"
	"sync"
)

var (
	ErrShutdownClose    = errors.New("shutdown finished with error(s)")
	ErrShutdownCanceled = errors.New("shutdown cancelled")
)

// CloseFunc функция для завершения сервиса
type CloseFunc func(ctx context.Context) bool

type Closer struct {
	mu        sync.Mutex
	closeFunc []CloseFunc
}

func (c *Closer) Register(cfn CloseFunc) {
	if cfn == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	c.closeFunc = append(c.closeFunc, cfn)
}

func (c *Closer) Close(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var (
		hasError bool
		complete = make(chan struct{}, 1)
	)
	go func() {
		for _, closer := range c.closeFunc {
			if ok := closer(ctx); !ok {
				hasError = true
			}
		}
		complete <- struct{}{}
	}()

	select {
	case <-complete:
		break
	case <-ctx.Done():
		return ErrShutdownCanceled
	}

	if hasError {
		return ErrShutdownClose
	}

	return nil
}
