package watcher_test

import (
	"context"
	"reflect"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Nomango/watcher"
)

func TestWatcher(t *testing.T) {
	v := int32(1)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	n := watcher.NewTickNotifier(time.Millisecond * 50)
	f := func(context.Context, interface{}) { atomic.AddInt32(&v, 1) }
	watcher.Watch(ctx, n, f)

	time.Sleep(time.Millisecond * 70)
	AssertEqual(t, int32(2), atomic.LoadInt32(&v))
	time.Sleep(time.Millisecond * 50)
	AssertEqual(t, int32(3), atomic.LoadInt32(&v))

	cancel()
	time.Sleep(time.Millisecond * 50)
	AssertEqual(t, int32(3), atomic.LoadInt32(&v))
}

func AssertEqual(t *testing.T, expect, actual interface{}) {
	if !reflect.DeepEqual(expect, actual) {
		t.Fatalf("values are not equal\nexpected=%v\ngot=%v", expect, actual)
	}
}
