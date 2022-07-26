package watcher

import (
	"context"
	"log"
)

// Watch executes f every time n triggers
func Watch(ctx context.Context, n Notifier, f Executer) {
	NewWatcher(n, f).Start(ctx)
}

type Executer = func(context.Context, interface{})

type Watcher struct {
	n      Notifier
	f      Executer
	stop   chan struct{}
	logger Logger
}

func NewWatcher(n Notifier, f Executer, opts ...WatcherOption) *Watcher {
	w := &Watcher{
		n:    n,
		f:    f,
		stop: make(chan struct{}),
	}
	for _, opt := range opts {
		opt(w)
	}
	if w.logger == nil {
		w.logger = &stdLogger{}
	}
	return w
}

func (w *Watcher) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case v, ok := <-w.n:
				if !ok {
					w.logger.Info("[watcher] notifier is closed")
					return
				}
				w.Execute(ctx, v)
			case <-ctx.Done():
				w.logger.Info("[watcher] context is done, err=%v", ctx.Err())
				return
			case <-w.stop:
				w.logger.Info("[watcher] stoped")
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
			w.logger.Error("[watcher] PANIC occurred!!! msg=%v", e)
		}
	}()
	w.f(ctx, v)
}

type WatcherOption func(*Watcher)

type Logger interface {
	Info(format string, args ...interface{})
	Error(format string, args ...interface{})
}

func WithLogger(logger Logger) WatcherOption {
	return func(w *Watcher) {
		w.logger = logger
	}
}

type stdLogger struct{}

func (l *stdLogger) Info(format string, args ...interface{}) {
	log.Printf(format, args...)
}

func (l *stdLogger) Error(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}
