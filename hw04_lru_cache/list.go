package hw04_lru_cache //nolint:golint,stylecheck

type List interface {
	Len() int                          // длина списка
	Front() *listItem                  // первый элемент списка
	Back() *listItem                   // последний элемент списка
	PushFront(v interface{}) *listItem // добавить значение в начало
	PushBack(v interface{}) *listItem  // добавить значение в конец
	Remove(i *listItem)                // удалить элемент
	MoveToFront(i *listItem)           // переместить элемент в начало
}

type listItem struct {
	Value interface{} // значение
	Prev  *listItem   // предыдущий элемент
	Next  *listItem   // следующий элемент
}

type list struct {
	len   int
	front *listItem
	back  *listItem
}

// случай l == nil проверять нигде не будем, структура не универсальная.

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *listItem {
	return l.front
}

func (l *list) Back() *listItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *listItem {
	item := &listItem{Value: v}

	if l.len == 0 {
		l.back = item
	} else {
		item.Next = l.front
		l.front.Prev = item
	}
	l.front = item

	l.len++
	return item
}

func (l *list) PushBack(v interface{}) *listItem {
	item := &listItem{Value: v}

	if l.len == 0 {
		l.front = item
	} else {
		item.Prev = l.back
		l.back.Next = item
	}
	l.back = item

	l.len++
	return item
}

func (l *list) Remove(i *listItem) {
	if i.Prev == nil {
		l.front = i.Next
	} else {
		i.Prev.Next = i.Next
	}

	if i.Next == nil {
		l.back = i.Prev
	} else {
		i.Next.Prev = i.Prev
	}

	l.len--
}

func (l *list) MoveToFront(i *listItem) {
	if i.Prev == nil {
		return
	}

	l.Remove(i)

	i.Prev = nil
	i.Next = l.front
	l.front.Prev = i
	l.front = i
	l.len++
}

func NewList() List {
	return &list{}
}
