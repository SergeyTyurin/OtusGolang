package hw04lrucache

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
}

type cacheValue struct {
	key   Key
	value interface{}
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	innerValue := cacheValue{key, value}
	if q, ok := l.items[key]; ok {
		q.Value = innerValue
		l.queue.MoveToFront(q)
		return true
	}

	newItem := l.queue.PushFront(innerValue)
	l.items[key] = newItem

	if l.queue.Len() > l.capacity {
		removedItem := l.queue.Back()
		l.queue.Remove(removedItem)
		removedData := removedItem.Value.(cacheValue)
		delete(l.items, removedData.key)
		removedItem.Value = nil
	}
	return false
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	if q, ok := l.items[key]; ok {
		l.queue.MoveToFront(q)
		return q.Value.(cacheValue).value, true
	}

	return nil, false
}

func (l *lruCache) Clear() {
	for key, value := range l.items {
		l.queue.Remove(value)
		delete(l.items, key)
	}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
