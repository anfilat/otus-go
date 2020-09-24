package hw04_lru_cache //nolint:golint,stylecheck

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool // Добавить значение в кэш по ключу
	Get(key Key) (interface{}, bool)     // Получить значение из кэша по ключу
	Clear()                              // Очистить кэш
}

// вся работа с указателями в одном месте

type itemsMap map[Key]*listItem

func listItemToCacheItem(item *listItem) *cacheItem {
	return item.Value.(*cacheItem)
}

func newCacheItem(key Key, value interface{}) *cacheItem {
	return &cacheItem{key, value}
}

// end работа с указателями

type lruCache struct {
	sync.Mutex
	capacity int
	queue    List
	items    itemsMap
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	l.Lock()
	defer l.Unlock()

	item, ok := l.items[key]

	if ok {
		l.queue.MoveToFront(item)
		listItemToCacheItem(item).value = value
		return true
	}

	if l.queue.Len() == l.capacity {
		removable := l.queue.Back()
		l.queue.Remove(removable)
		delete(l.items, listItemToCacheItem(removable).key)
	}

	item = l.queue.PushFront(newCacheItem(key, value))
	l.items[key] = item

	return false
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	l.Lock()
	defer l.Unlock()

	item, ok := l.items[key]

	if ok {
		l.queue.MoveToFront(item)
		return listItemToCacheItem(item).value, true
	}

	return nil, false
}

func (l *lruCache) Clear() {
	l.Lock()
	defer l.Unlock()

	for l.queue.Len() > 0 {
		l.queue.Remove(l.queue.Back())
	}
	l.items = make(itemsMap, l.capacity)
}

type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(itemsMap, capacity),
	}
}
