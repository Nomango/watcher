package watcher

import (
	"context"
	"sync"
	"sync/atomic"
)

// AutoLoad updates returned atomic value every time n triggers
func AutoLoad(ctx context.Context, n Notifier, opts ...LoaderOption) *atomic.Value {
	l := NewLoader(n, opts...)
	l.Start(ctx)
	return &l.v
}

// WatchLoader executes f every time l updates
func WatchLoader(l *Loader, f Executer) {
	w := NewWatcher(nil, f)
	l.apply(
		WithTransformer(func(ctx context.Context, origin interface{}) interface{} {
			go w.Execute(ctx, origin)
			return origin
		}),
	)
}

type Transformer = func(context.Context, interface{}) interface{}

type Loader struct {
	*Watcher
	v          atomic.Value
	mu         sync.Mutex
	transforms []Transformer
}

func NewLoader(n Notifier, opts ...LoaderOption) *Loader {
	l := &Loader{}
	l.Watcher = NewWatcher(n, l.receive)
	l.apply(opts...)
	return l
}

func (l *Loader) Get() interface{} {
	return l.v.Load()
}

func (l *Loader) apply(opts ...LoaderOption) {
	l.mu.Lock()
	for _, opt := range opts {
		opt(l)
	}
	l.mu.Unlock()
}

func (l *Loader) receive(ctx context.Context, v interface{}) {
	for _, t := range l.transforms {
		v = t(ctx, v)
	}
	l.v.Store(v)
}

type LoaderOption func(*Loader)

func WithTransformer(t Transformer) LoaderOption {
	return func(l *Loader) {
		l.transforms = append(l.transforms, t)
	}
}
