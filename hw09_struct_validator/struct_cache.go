package hw09_struct_validator //nolint:golint,stylecheck

import (
	"reflect"
	"sync/atomic"
)

type structCache struct {
	values atomic.Value // map[reflect.Type]structRules
}

func newStructCache() *structCache {
	values := make(map[reflect.Type]structRules)
	cache := &structCache{}
	cache.values.Store(values)
	return cache
}

func (sc *structCache) lookup(value reflect.Type) (structRules, bool) {
	v, ok := sc.values.Load().(map[reflect.Type]structRules)[value]
	return v, ok
}

func (sc *structCache) add(value reflect.Type, rules structRules) {
	values := sc.values.Load().(map[reflect.Type]structRules)
	newValues := make(map[reflect.Type]structRules, len(values)+1)
	for k, v := range values {
		newValues[k] = v
	}
	newValues[value] = rules
	sc.values.Store(newValues)
}
