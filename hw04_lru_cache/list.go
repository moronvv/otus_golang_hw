package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	front  *ListItem
	back   *ListItem
	length int
}

func NewList() List {
	return new(list)
}

func (l *list) Len() int {
	return l.length
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	newFront := ListItem{Value: v}

	oldFront := l.front
	if oldFront != nil {
		oldFront.Prev = &newFront
		newFront.Next = oldFront
		// set new front
		l.front = &newFront
	} else {
		// first element init
		l.front = &newFront
		l.back = &newFront
	}

	// inc length
	l.length++

	return l.front
}

func (l *list) PushBack(v interface{}) *ListItem {
	newBack := ListItem{Value: v}

	oldBack := l.back
	if oldBack != nil {
		oldBack.Next = &newBack
		newBack.Prev = oldBack
		// set new back
		l.back = &newBack
	} else {
		// first element init
		l.front = &newBack
		l.back = &newBack
	}

	// inc length
	l.length++

	return l.back
}

func (l *list) Remove(i *ListItem) {
	// corner cases
	if i == l.front {
		l.front = i.Next
	}
	if i == l.back {
		l.back = i.Prev
	}

	// swap
	if i.Prev != nil {
		i.Prev.Next = i.Next
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
	}

	// dec length
	l.length--
}

func (l *list) MoveToFront(i *ListItem) {
	l.Remove(i)
	l.PushFront(i.Value)
}
