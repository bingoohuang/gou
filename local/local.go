package local

import (
	"context"
	"sync"
)

// nolint gochecknoglobals
var locals = &struct {
	sync.RWMutex
	ctx map[uint64]context.Context
}{
	ctx: make(map[uint64]context.Context),
}

func get(gid uint64) context.Context {
	locals.RLock()
	ctx := locals.ctx[gid]
	locals.RUnlock()

	if ctx == nil {
		ctx = context.Background()
	}

	return ctx
}

// nolint golint
func set(gid uint64, ctx context.Context) {
	locals.Lock()
	locals.ctx[gid] = ctx
	locals.Unlock()
}

func temp(gid uint64, key, val interface{}) context.Context {
	ctx := context.WithValue(get(gid), key, val)
	set(gid, ctx)

	return ctx
}

func clear(gid uint64) context.Context {
	locals.Lock()
	ctx := locals.ctx[gid]
	delete(locals.ctx, gid)
	locals.Unlock()

	return ctx
}

// Get ...
func Get() context.Context {
	return get(Goid())
}

// Set ...
func Set(ctx context.Context) {
	set(Goid(), ctx)
}

// ...
func Clear() context.Context {
	return clear(Goid())
}

// Value ...
func Temp(key, val interface{}) context.Context {
	return temp(Goid(), key, val)
}

// Value ...
func Value(key interface{}) interface{} {
	ctx := Get()
	if ctx == nil {
		return nil
	}

	return ctx.Value(key)
}

// Go ...
func Go(fn func()) {
	GoContext(Get(), fn)
}

// GoContext ...
func GoContext(ctx context.Context, fn func()) {
	go func() {
		Set(ctx)

		defer Clear()

		fn()
	}()
}
