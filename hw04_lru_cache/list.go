package hw04lrucache

// List nil <- (prev) front <-> ... <-> elem <-> ... <-> back (next) -> nil.
type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

// ListItem linked list item.
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

func (l list) Len() int {
	return l.length
}

func (l list) Front() *ListItem {
	return l.front
}

func (l list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	curFront := l.front
	item := &ListItem{
		Value: v,
		Next:  curFront,
	}
	l.front = item
	if curFront == nil {
		l.back = item
	} else {
		curFront.Prev = item
	}
	l.length++
	return item
}

func (l *list) PushBack(v interface{}) *ListItem {
	curBack := l.back
	item := &ListItem{
		Value: v,
		Prev:  curBack,
	}
	l.back = item
	if curBack == nil {
		l.front = item
	} else {
		curBack.Next = item
	}
	l.length++
	return item
}

func (l *list) Remove(i *ListItem) {
	if i.Prev == nil && i.Next == nil {
		// единственный элемент
		l.front = nil
		l.back = nil
	} else if i.Prev == nil {
		// это front
		i.Next.Prev = nil
		l.front = i.Next
	} else if i.Next == nil {
		// это back
		i.Prev.Next = nil
		l.back = i.Prev
	} else {
		i.Prev.Next, i.Next.Prev = i.Next, i.Prev
	}
	// на всякий случай
	i.Prev = nil
	i.Next = nil
	l.length--
}

func (l *list) MoveToFront(i *ListItem) {
	if i == l.Front() {
		return
	}
	l.Remove(i)
	l.PushFront(i.Value)
}

func NewList() List {
	return new(list)
}
