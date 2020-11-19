package hw09_struct_validator //nolint:golint,stylecheck

import (
	"reflect"
	"sync"
)

type structCache struct {
	sync.RWMutex
	values map[reflect.Type]structRules
}

func newStructCache() *structCache {
	result := &structCache{}
	result.values = make(map[reflect.Type]structRules)
	return result
}

func (sc *structCache) lookup(value reflect.Type) (structRules, bool) {
	sc.RLock()
	defer sc.RUnlock()

	v, ok := sc.values[value]
	return v, ok
}

func (sc *structCache) add(value reflect.Type, rules structRules) {
	sc.Lock()
	defer sc.Unlock()

	sc.values[value] = rules
}
