package booter

import (
	"context"
	"sync"
)

var defaultBootContext = NewBootContext()

func DefaultContext() *BootContext {
	return defaultBootContext
}

type BootContext struct {
	components map[string]Component // component作为创建实例的模板
	ctx        context.Context
	container  map[string]Component // container作为实例装载的容器
	m          map[string]any       // 用户自定义map

	circleCheck  map[string]struct{}
	circleRecord [][]any

	lock sync.RWMutex
}

func NewBootContext() *BootContext {
	return &BootContext{
		components: make(map[string]Component),
		ctx:        context.Background(),
		container:  make(map[string]Component),
		m:          make(map[string]any),
	}
}

func (b *BootContext) GetContext() context.Context {
	return b.ctx
}

func (b *BootContext) Set(key string, value any) {
	b.m[key] = value
}

func (b *BootContext) Get(key string, defaultValue ...any) any {
	val, ok := b.m[key]
	if ok {
		return val
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return nil
}

func (b *BootContext) GetExists(key string) (any, bool) {
	val, ok := b.m[key]
	return val, ok
}
