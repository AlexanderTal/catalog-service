package processor

import (
	"context"
	"io"
	"sync"
	"time"
)

func Wrap(ctx context.Context, wg *sync.WaitGroup, cb func(context.Context)) {
	if wg != nil {
		wg.Add(1)
	}
	go func() {
		if wg != nil {
			defer wg.Done()
		}
		select {
		case <-ctx.Done():
			return
		default:
			cb(ctx)
		}
	}()
}

type (
	CloserFunc        func() error
	CloserContextFunc = func(ctx context.Context) error
)

func (f CloserFunc) Close() error {
	return f()
}

func NewCloserContextFunc(
	f CloserContextFunc,
	ctx context.Context, timeout time.Duration,
) CloserFunc {
	return func() error {
		if timeout > 0 {
			var cancelFunc context.CancelFunc
			ctx, cancelFunc = context.WithTimeout(ctx, timeout)
			defer cancelFunc()
		}
		return f(ctx)
	}
}

func WatchForShutdown(ctx context.Context, wg *sync.WaitGroup, closer io.Closer) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		_ = closer.Close()
	}()
}
