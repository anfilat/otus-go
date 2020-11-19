package hw09_struct_validator //nolint:golint,stylecheck

import (
	"regexp"
	"sync"
)

type regexpCache struct {
	sync.RWMutex
	values map[string]*regexp.Regexp
}

func newRegexpCache() *regexpCache {
	result := &regexpCache{}
	result.values = make(map[string]*regexp.Regexp)
	return result
}

func (rc *regexpCache) get(pattern string) (*regexp.Regexp, error) {
	rg, ok := rc._lookup(pattern)
	if ok {
		return rg, nil
	}

	rg, err := rc._add(pattern)
	if err != nil {
		return nil, err
	}
	return rg, nil
}

func (rc *regexpCache) _lookup(pattern string) (*regexp.Regexp, bool) {
	rc.RLock()
	defer rc.RUnlock()

	rg, ok := rc.values[pattern]
	return rg, ok
}

func (rc *regexpCache) _add(pattern string) (*regexp.Regexp, error) {
	rg, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	rc.Lock()
	defer rc.Unlock()

	rc.values[pattern] = rg
	return rg, nil
}
