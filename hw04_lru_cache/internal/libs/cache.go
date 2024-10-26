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
	item := &ListItem{
		Value: value,
	}
	if val, exist := l.items[key]; exist {
		l.items[key] = item
		l.queue.MoveToFront(item)
		val.Value = value
		return true
	}
	l.items[key] = item
	l.queue.PushFront(item)
	if l.queue.Len() > l.capacity {
		lastRecently := l.queue.Back()
		l.queue.Remove(lastRecently)
		l.removeFromItems(lastRecently)
	}
	return false
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if val, exist := l.items[key]; exist {
		l.queue.MoveToFront(val)
		return val.Value, true
	}
	return nil, false
}

func (l *lruCache) Clear() {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.items = make(map[Key]*ListItem)
	l.queue = new(list)
}

func (l *lruCache) removeFromItems(elem *ListItem) {
	for k, v := range l.items {
		if elem == v {
			delete(l.items, k)
			break
		}
	}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
