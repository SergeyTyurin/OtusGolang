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
	len  int
	head *ListItem
	tail *ListItem
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.head
}

func (l *list) Back() *ListItem {
	return l.tail
}

func (l *list) PushFront(v interface{}) (newItem *ListItem) {
	newItem = new(ListItem)
	newItem.Value = v
	if l.head != nil {
		newItem.Next, l.head.Prev = l.head, newItem
		l.head = newItem
	} else {
		l.head = newItem
		l.tail = l.head
	}
	l.len++
	return
}

func (l *list) PushBack(v interface{}) (newItem *ListItem) {
	newItem = new(ListItem)
	newItem.Value = v
	if l.tail != nil {
		newItem.Prev, l.tail.Next = l.tail, newItem
		l.tail = newItem
	} else {
		l.tail = newItem
		l.head = l.tail
	}
	l.len++
	return
}

func (l *list) Remove(i *ListItem) {
	defer func() {
		switch l.len {
		case 0:
			return
		default:
			l.len--
		}
	}()
	if l.len < 2 {
		l.head, l.tail = nil, nil
		return
	}
	switch i {
	case l.head:
		l.head = i.Next
		l.head.Prev = nil
	case l.tail:
		l.tail = i.Prev
		l.tail.Next = nil
	default:
		i.Prev.Next, i.Next.Prev = i.Next, i.Prev
	}
}

func (l *list) MoveToFront(i *ListItem) {
	if l.len == 0 || i == l.head {
		return
	}

	if l.head != nil {
		l.Remove(i)
		l.len++
		i.Prev, i.Next, l.head.Prev = nil, l.head, i
		l.head = i
	}
}

func NewList() List {
	return new(list)
}
