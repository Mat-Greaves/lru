// Package list implents a generic doubly linked list
// list is a generic port of containers/list
package list

type List[T any] struct {
	root Element[T] // sentinal element, only used for next and prev
	len  int
}

func New[T any]() *List[T] {
	return new(List[T]).Init()
}

// Init initialises or clears a list
func (l *List[T]) Init() *List[T] {
	l.root.prev = &l.root
	l.root.next = &l.root
	l.len = 0
	return l
}

// lazyInit lazily initialises a zero List value
func (l *List[T]) lazyInit() {
	if l.root.next == nil {
		l.Init()
	}
}

func (l *List[T]) Len() int {
	return l.len
}

// Front returns the first element of l or nil if the list is empty
func (l *List[T]) Front() *Element[T] {
	if l.len == 0 {
		return nil
	}
	return l.root.next
}

// Back returns the last element of l or nil if the list is empty
func (l *List[T]) Back() *Element[T] {
	if l.len == 0 {
		return nil
	}
	return l.root.prev
}

// Remove removes e from l if e is an element of list l.
// It returns the element value e.Value
// The elemenent must not be nil
func (l *List[T]) Remove(e *Element[T]) T {
	if e.list == l {
		l.remove(e)
	}
	return e.Value
}

func (l *List[T]) PushFront(v T) *Element[T] {
	l.lazyInit()
	return l.insertValue(v, &l.root)
}

func (l *List[T]) PushBack(v T) *Element[T] {
	l.lazyInit()
	return l.insertValue(v, l.root.prev)
}

// insert inserts e after at, incrementing l.len and returning e
func (l *List[T]) insert(e, at *Element[T]) *Element[T] {
	e.prev = at
	e.next = at.next
	e.prev.next = e
	e.next.prev = e
	e.list = l
	l.len++
	return e
}

// insertValue is a convinience wrapper around insert(&Element[T]{Value: v})
func (l *List[T]) insertValue(v T, at *Element[T]) *Element[T] {
	return l.insert(&Element[T]{Value: v}, at)
}

func (l *List[T]) remove(e *Element[T]) {
	e.prev.next = e.next
	e.next.prev = e.prev
	e.prev = nil // get rid of references to prevent memory leaks
	e.next = nil // TODO: find link explaining this
	e.list = nil
	l.len--
}

// move moves e next to at

type Element[T any] struct {
	Value T
	list  *List[T]
	prev  *Element[T]
	next  *Element[T]
}

// Next returns the next list element or nil
func (e *Element[T]) Next() *Element[T] {
	if n := e.next; e.list != nil && n != &e.list.root {
		return n
	}
	return nil
}

// Prev returns the previous list element or nil
func (e *Element[T]) Prev() *Element[T] {
	if p := e.prev; e.list != nil && p != &e.list.root {
		return p
	}
	return nil
}

// MoveToFront moves element e to the front of list l.
// If e is not an element of l, the list is not modified.
// e must not be nil.
func (l *List[T]) MoveToFront(e *Element[T]) {
	if e.list != l || l.root.next == e {
		return
	}
	move(e, &e.list.root)
}

// MoveToBack moves element e to the back of List l.
// If e is not a member of list l, the list is not modified.
// e must not be nil.
func (l *List[T]) MoveToBack(e *Element[T]) {
	if e.list != l || l.root.prev == e {
		return
	}
	move(e, e.list.root.prev)
}

func move[T any](e, at *Element[T]) {
	if e == at {
		return
	}
	e.prev.next = e.next
	e.next.prev = e.prev

	e.prev = at
	e.next = at.next
	e.prev.next = e
	e.next.prev = e
}
