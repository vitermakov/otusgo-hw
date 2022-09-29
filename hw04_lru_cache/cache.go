package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
	mx       sync.Mutex
}

// public
func (lc *lruCache) Set(key Key, value interface{}) (exists bool) {
	lc.safeExec(func() {
		exists = lc.set(key, value)
	})
	return
}
func (lc *lruCache) Get(key Key) (value interface{}, exists bool) {
	lc.safeExec(func() {
		value, exists = lc.get(key)
	})
	return
}
func (lc *lruCache) Clear() {
	lc.safeExec(lc.clear)
}

// private
func (lc *lruCache) safeExec(unsafeFunc func()) {
	lc.mx.Lock()
	unsafeFunc()
	lc.mx.Unlock()
}
func (lc *lruCache) set(key Key, value interface{}) bool {
	item, exists := lc.items[key]
	val := cacheItem{key: key, value: value}
	if exists {
		item.Value = val
		lc.queue.MoveToFront(item)
	} else {
		item = lc.queue.PushFront(val)
		lc.items[key] = item
		for lc.queue.Len() > lc.capacity {
			back := lc.queue.Back()
			lc.queue.Remove(back)
			delete(lc.items, back.Value.(cacheItem).key)
		}
	}
	return exists
}
func (lc *lruCache) get(key Key) (interface{}, bool) {
	item, exists := lc.items[key]
	if exists {
		lc.queue.MoveToFront(item)
		return item.Value.(cacheItem).value, true
	}
	return nil, false
}
func (lc *lruCache) clear() {
	lc.queue = NewList()
	lc.items = make(map[Key]*ListItem, lc.capacity)
}

type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
