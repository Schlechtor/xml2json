package transformations

import (
	"sync"
)

type Transformer func(data map[string]interface{}, param any) interface{}

var (
	transformerRegistry = make(map[string]Transformer)
	mu                  sync.RWMutex
)

func Register(name string, transformer Transformer) {
	mu.Lock()
	defer mu.Unlock()
	transformerRegistry[name] = transformer
}

func GetTransformer(name string) (Transformer, bool) {
	mu.RLock()
	defer mu.RUnlock()
	transformer, exists := transformerRegistry[name]
	return transformer, exists
}
