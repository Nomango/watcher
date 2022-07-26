# Watcher

简单的监听和加载工具，帮助对依赖复杂的模块进行解耦。

## Quick Start

```golang
// a.go
// A模块每隔一段时间拉取一次值
updateInterval := time.Minute // 更新间隔
notifier := watcher.NewTickNotifier(updateInterval) // tick触发器
loaderA := watcher.NewLoader(ctx, notifier, watcher.WithTransformer(func(ctx context.Context, v interface{}) interface{} {
    return serviceA.GetSomething() // 执行拉取函数
}))
loaderA.Start(ctx) // 启动加载器

// b.go
// B模块不关心A模块的值什么时候更新，直接使用即可
func AnyWhere() {
    v := loaderA.Get() // 取到A模块的值
    // do something with v
}

// c.go
// C模块监听A模块的值更新
watcher.WatchLoader(loaderA, func(ctx context.Context, v interface{}) {
    fmt.Println(v) // 每次A模块更新值时进行打印
})
```

## Usage

### Watcher

每隔一段时间执行一次指定函数：

```golang
// 每秒打印一次
notifier := watcher.NewTickNotifier(time.Second)
watcher.Watch(ctx, notifier, func(ctx context.Context, v interface{}) {
    fmt.Println("hello world")
})

// Output:
// hello world
```

自定义Notifier

```golang
ch := make(chan int)
watcher.Watch(ctx, watcher.NewNotifier(ch), func(ctx context.Context, v interface{}) {
    fmt.Println(v)
})
ch <- 1
ch <- 2

// Output:
// 1
// 2
```

控制启动停止

```golang
notifier := watcher.NewTickNotifier(time.Second)
watcher.NewWatcher(notifier, func(ctx context.Context, v interface{}) {
    fmt.Println("hello world")
})

watcher.Start(ctx)
watcher.Stop()
```

### Loader

自动加载器：

```golang
ch := make(chan int)
v := watcher.AutoLoad(ctx, watcher.NewNotifier(ch))

fmt.Println(v.Load())

ch <- 1
fmt.Println(v.Load()) // 更新延迟，也许打印出的并不是最新值

// Output:
// nil
// 1
```

加载器可以添加转换函数：

```golang
notifier := watcher.NewTickNotifier(time.Second)
v := watcher.AutoLoad(ctx, notifier, watcher.WithTransformer(func(ctx context.Context, v interface{}) interface{} {
    // TickNotifier 每次触发会送一个 time.Time 过来
    return v.(time.Time).Unix()
}))

time.Sleep(time.Second)
fmt.Println(v.Load())

time.Sleep(time.Second)
fmt.Println(v.Load())
```

可以监听 Loader 更新：

```golang
ch := make(chan int)
loader := watcher.NewLoader(watcher.NewNotifier(ch))
watcher.WatchLoader(loader, func(ctx context.Context, v interface{}) {
    fmt.Println(v)
})

ch <- 1
ch <- 2

// Output:
// 1
// 2
```
