package libs

import "sync"

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

const (
	ZERO = iota
)

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	List
	front *ListItem
	back  *ListItem
	size  int
	mutex sync.Mutex
}

func (l *list) MoveToFront(i *ListItem) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if i == nil || l.isEmpty() || l.front == i {
		return
	}
	if i.Prev != nil {
		i.Prev.Next = i.Next
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
	}
	if l.back == i {
		l.back = i.Prev
	}
	i.Prev = nil
	i.Next = l.front
	if l.front != nil {
		l.front.Prev = i
	}
	l.front = i
	if l.back == nil {
		l.back = i
	}
}

func (l *list) Remove(i *ListItem) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if i == nil || l.isEmpty() {
		return
	}
	if i.Prev != nil {
		i.Prev.Next = i.Next
	} else {
		l.front = i.Next
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.back = i.Prev
	}
	i.Prev = nil
	i.Next = nil
	l.size--
}

func (l *list) PushFront(v interface{}) *ListItem {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	elem := &ListItem{
		Value: v,
	}
	if l.isEmpty() {
		l.init(elem)
	} else {
		l.front.Prev = elem
		elem.Next = l.front
		l.front = elem
	}
	l.size++
	return elem
}

func (l *list) PushBack(v interface{}) *ListItem {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	elem := &ListItem{
		Value: v,
	}
	if l.isEmpty() {
		l.init(elem)
	} else {
		l.back.Next = elem
		elem.Prev = l.back
		l.back = elem
	}
	l.size++
	return elem
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) Len() int {
	return l.size
}

func (l *list) isEmpty() bool {
	return l.Len() == ZERO
}

func (l *list) init(elem *ListItem) {
	l.front = elem
	l.back = elem
}

func NewList() List {
	return new(list)
}
