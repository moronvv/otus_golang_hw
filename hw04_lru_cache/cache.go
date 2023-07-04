package hw04lrucache

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
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	item, found := c.items[key]
	if found {
		item.Value = value
		c.queue.MoveToFront(item)
		c.items[key] = item
	} else {
		newItem := c.queue.PushFront(value)
		c.items[key] = newItem

		if c.queue.Len() > c.capacity {
			c.queue.Remove(c.queue.Back())
			delete(c.items, key)
		}
	}

	return found
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	var val interface{}

	item, found := c.items[key]
	if found {
		c.queue.MoveToFront(item)
		val = item.Value
	} else {
		val = nil
	}

	return val, found
}

func (c *lruCache) Clear() {
	c.queue = NewList()
	c.items = make(map[Key]*ListItem)
}
