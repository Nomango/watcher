package watcher_test

import (
	"context"
	"testing"
	"time"

	"github.com/Nomango/watcher"
)

func TestLoader(t *testing.T) {
	ch := make(chan int)
	transform := func(ctx context.Context, i interface{}) interface{} {
		return i.(int) + 100
	}
	alv := watcher.AutoLoad(context.Background(), watcher.NewNotifier(ch), watcher.WithTransformer(transform))

	if alv.Load() != nil {
		t.Fatalf("not expected")
	}
	AssertEqual(t, nil, alv.Load())

	ch <- 1
	time.Sleep(time.Millisecond * 50)
	AssertEqual(t, 101, alv.Load())

	ch <- 2
	time.Sleep(time.Millisecond * 50)
	AssertEqual(t, 102, alv.Load())
}
