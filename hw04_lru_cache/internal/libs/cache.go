package libs

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	Cache
	mutex    sync.Mutex
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if value == nil {
		return false
	}
	item := NewCacheItem(value, key)
	if cur, isExist := l.items[key]; isExist {
		cur.Value = item
		l.queue.MoveToFront(cur)
		return true
	}
	saveQueue := l.queue.PushFront(item)
	l.items[key] = saveQueue
	if l.queue.Len() == l.capacity {
		lastRecently := l.queue.Back()
		if lastRecently != nil {
			l.queue.Remove(lastRecently)
			delete(l.items, lastRecently.Value.(*CacheItem).Key)
		}
	}
	return false
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if val, isExist := l.items[key]; isExist {
		l.queue.MoveToFront(val)
		return val.Value, true
	}
	return nil, false
}

func (l *lruCache) Clear() {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.items = make(map[Key]*ListItem, l.capacity)
	l.queue = NewList()
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
