package watcher

import (
	"context"

	"github.com/Nomango/ark/logs"
)

// Watch executes f every time n triggers
func Watch(ctx context.Context, n Notifier, f Executer) {
	NewWatcher(n, f).Start(ctx)
}

type Executer = func(context.Context, interface{})

type Watcher struct {
	n    Notifier
	f    Executer
	stop chan struct{}
}

func NewWatcher(n Notifier, f Executer) *Watcher {
	return &Watcher{
		n:    n,
		f:    f,
		stop: make(chan struct{}),
	}
}

func (w *Watcher) Start(ctx context.Context) {
	ctx = logs.CtxWithKVs(ctx, logs.KV("from", "watcher"))
	go func() {
		for {
			select {
			case v, ok := <-w.n:
				if !ok {
					logs.CtxNoticef(ctx, "notifier is closed")
					return
				}
				w.Execute(ctx, v)
			case <-ctx.Done():
				logs.CtxNoticef(ctx, "context is done, err=%v", ctx.Err())
				return
			case <-w.stop:
				logs.CtxNoticef(ctx, "watcher is stoped")
				return
			}
		}
	}()
}

func (w *Watcher) Stop() {
	w.stop <- struct{}{}
}

func (w *Watcher) GetNotifier() Notifier {
	return w.n
}

func (w *Watcher) Execute(ctx context.Context, v interface{}) {
	defer func() {
		if e := recover(); e != nil {
			logs.CtxErrorf(ctx, "PANIC occurred!!! msg=%v", e)
		}
	}()
	w.f(ctx, v)
}
