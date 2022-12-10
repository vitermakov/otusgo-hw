package hw04lrucache

import "sync"

type Key string

// Cache common simple interface.
// Add Unset method. Set nil is not unset.
type Cache interface {
	Set(key Key, value interface{}) bool
	Unset(key Key) (interface{}, bool)
	Get(key Key) (interface{}, bool)
	Clear()
}

// LruCache LRU type caching.
type LruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
	sync.Mutex
}

// public

// Set thread-safe setting value by key.
// Returns existing-flag before setting, value can be nil.
func (lc *LruCache) Set(key Key, value interface{}) (exists bool) {
	lc.safeExec(func() {
		exists = lc.set(key, value)
	})
	return
}

// Get thread-safe getting value by key.
func (lc *LruCache) Get(key Key) (value interface{}, exists bool) {
	lc.safeExec(func() {
		value, exists = lc.get(key)
	})
	return
}

// Unset thread-safe unset value by key.
// Returns value for possible using.
func (lc *LruCache) Unset(key Key) (value interface{}, exists bool) {
	lc.safeExec(func() {
		value, exists = lc.unset(key)
	})
	return
}

// Clear thread-safe clear all cache.
func (lc *LruCache) Clear() {
	lc.safeExec(lc.clear)
}

// private

func (lc *LruCache) safeExec(unsafeFunc func()) {
	lc.Lock()
	defer lc.Unlock()

	unsafeFunc()
}

func (lc *LruCache) set(key Key, value interface{}) bool {
	item, exists := lc.items[key]
	val := cacheItem{key: key, value: value}
	if exists {
		item.Value = val
		lc.queue.MoveToFront(item)
	} else {
		if lc.queue.Len() == lc.capacity {
			back := lc.queue.Back()
			lc.queue.Remove(back)
			delete(lc.items, back.Value.(cacheItem).key)
		}
		item = lc.queue.PushFront(val)
		lc.items[key] = item
	}
	return exists
}

func (lc *LruCache) get(key Key) (interface{}, bool) {
	item, exists := lc.items[key]
	if exists {
		lc.queue.MoveToFront(item)
		return item.Value.(cacheItem).value, true
	}
	return nil, false
}

func (lc *LruCache) unset(key Key) (interface{}, bool) {
	value, exists := lc.get(key)
	if exists {
		item := lc.items[key]
		lc.queue.Remove(item)
		delete(lc.items, key)
	}
	return value, exists
}

func (lc *LruCache) clear() {
	lc.queue = NewList()
	lc.items = make(map[Key]*ListItem, lc.capacity)
}

type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	return &LruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
