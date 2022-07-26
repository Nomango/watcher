package watcher

import (
	"reflect"
	"time"
)

type Notifier <-chan interface{}

// NewNotifier takes any readable channel type (chan or <-chan but not chan<-) and
// exposes it as a Notifier
func NewNotifier(ch interface{}) Notifier {
	return Notifier(wrapChannel(ch))
}

func NewTickNotifier(interval time.Duration) Notifier {
	t := time.NewTicker(interval)
	return NewNotifier(t.C)
}

func wrapChannel(ch interface{}) <-chan interface{} {
	t := reflect.TypeOf(ch)
	if t.Kind() != reflect.Chan || t.ChanDir()&reflect.RecvDir == 0 {
		panic("channels: input to Wrap must be readable channel")
	}
	realChan := make(chan interface{})

	go func() {
		v := reflect.ValueOf(ch)
		for {
			x, ok := v.Recv()
			if !ok {
				close(realChan)
				return
			}
			realChan <- x.Interface()
		}
	}()
	return realChan
}
