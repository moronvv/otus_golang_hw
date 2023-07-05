package hw04lrucache

import (
	"sync"
)

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	queue    List
	items    map[Key]*ListItem
	capacity int
	mutex    sync.Mutex
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

// Item structure to store in queue.
type queueItem struct {
	Value interface{}
	Key   Key
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	item, found := c.items[key]
	if found {
		item.Value = queueItem{
			Key:   key,
			Value: value,
		}
		c.queue.MoveToFront(item)
		c.items[key] = item
	} else {
		newItem := c.queue.PushFront(queueItem{
			Key:   key,
			Value: value,
		})
		c.items[key] = newItem

		if c.queue.Len() > c.capacity {
			itemToDiscard := c.queue.Back()
			c.queue.Remove(itemToDiscard)

			itemToDiscardKey := itemToDiscard.Value.(queueItem).Key
			delete(c.items, itemToDiscardKey)
		}
	}

	return found
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	var val interface{}
	item, found := c.items[key]
	if found {
		c.queue.MoveToFront(item)
		val = item.Value.(queueItem).Value
	} else {
		val = nil
	}

	return val, found
}

func (c *lruCache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.queue = NewList()
	c.items = make(map[Key]*ListItem)
}
