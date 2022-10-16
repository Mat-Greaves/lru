// see https://cs.opensource.google/go/go/+/master:src/container/list/list_test.go

package list

import (
	"testing"
)

func checkListLen[T any](t *testing.T, l *List[T], len int) bool {
	if n := l.Len(); n != len {
		t.Errorf("l.Len() = %d, want %d", n, len)
		return false
	}
	return true
}

func checkListPointers[T any](t *testing.T, l *List[T], es []*Element[T]) {
	root := &l.root

	if !checkListLen(t, l, len(es)) {
		return
	}

	if len(es) == 0 {
		// zero length lists musts be the zero value or a propperly initialised list
		if l.root.next != nil && l.root.next != root || l.root.prev != nil && l.root.prev != root {
			t.Errorf("l.root.next = %p, l.root.prev = %p; both should be nil or %p", l.root.next, l.root.prev, root)
		}
		return
	}
	// len(es) > 0

	// check internal and external prev/next connections
	for i, e := range es {
		prev := root
		Prev := (*Element[T])(nil)
		if i > 0 {
			prev = es[i-1]
			Prev = prev
		}

		if p := e.prev; p != prev {
			t.Errorf("elt[%d](%p).prev = %p, want %p", i, e, p, prev)
		}
		if p := e.Prev(); p != Prev {
			t.Errorf("elt[%d](%p).Prev = %p, want %p", i, e, p, Prev)
		}

		next := root
		Next := (*Element[T])(nil)
		if i < len(es)-1 {
			next = es[i+1]
			Next = next
		}

		if n := e.next; n != next {
			t.Errorf("elt[%d](%p).next = %p, want %p", i, e, n, next)
		}
		if n := e.Next(); n != Next {
			t.Errorf("elt[%d](%p).Next = %p, want %p", i, e, n, Next)
		}
	}
	// loop over each element of
}

func Test_PushFront(t *testing.T) {
	t.Run("push into empty list", func(t *testing.T) {
		l := New[int]()
		e := l.PushFront(1)
		checkListPointers(t, l, []*Element[int]{e})
	})

	t.Run("push into empty unitialised list", func(t *testing.T) {
		l := List[int]{}
		e := l.PushFront(1)
		checkListPointers(t, &l, []*Element[int]{e})
	})

	t.Run("push multiple entries", func(t *testing.T) {
		l := New[int]()
		e1 := l.PushFront(0)
		e2 := l.PushFront(0)
		checkListPointers(t, l, []*Element[int]{e2, e1})
	})
}

func Test_Front(t *testing.T) {
	t.Run("Front of empty list", func(t *testing.T) {
		l := New[int]()
		if f := l.Front(); f != nil {
			t.Errorf("l.Front() = %p, want %p", f, (*Element[int])(nil))
		}
	})

	t.Run("Front of non-empty list", func(t *testing.T) {
		l := New[int]()
		e := l.PushFront(0)
		if f := l.Front(); f != e {
			t.Errorf("l.Front() = %p, want %p", f, e)
		}
	})
}

func Test_Back(t *testing.T) {
	t.Run("Back of empty list", func(t *testing.T) {
		l := New[int]()
		if f := l.Back(); f != nil {
			t.Errorf("l.Back() = %p, want %p", f, (*Element[int])(nil))
		}
	})

	t.Run("Back of non-empty list", func(t *testing.T) {
		l := New[int]()
		e := l.PushBack(0)
		if f := l.Back(); f != e {
			t.Errorf("l.Back() = %p, want %p", f, e)
		}
	})
}

func Test_PushBack(t *testing.T) {
	t.Run("push into empty list", func(t *testing.T) {
		l := New[int]()
		e := l.PushBack(0)
		checkListPointers(t, l, []*Element[int]{e})
	})

	t.Run("push into empty unitialised list", func(t *testing.T) {
		l := List[int]{}
		e := l.PushBack(0)
		checkListPointers(t, &l, []*Element[int]{e})
	})

	t.Run("push multiple", func(t *testing.T) {
		l := New[int]()
		e := l.PushBack(0)
		e2 := l.PushBack(0)
		checkListPointers(t, l, []*Element[int]{e, e2})
	})
}

func Test_Remove(t *testing.T) {
	t.Run("remove element not in list", func(t *testing.T) {
		l := New[int]()
		e2 := &Element[int]{}
		e := l.PushFront(0)
		if v := l.Remove(e2); v != e2.Value {
			t.Errorf("l.Remove(e2) = %d, want %d", v, e2.Value)
		}
		checkListPointers(t, l, []*Element[int]{e})
	})

	t.Run("remove element from list", func(t *testing.T) {
		l := New[int]()
		e := l.PushFront(0)
		if v := l.Remove(e); v != e.Value {
			t.Errorf("l.Remove(e) = %d, want %d", v, e.Value)
		}
		checkListPointers(t, l, []*Element[int]{})
	})
}

func Test_MoveToFront(t *testing.T) {
	t.Run("move element not part of list", func(t *testing.T) {
		l := New[int]()
		e := &Element[int]{}
		l.MoveToFront(e)
		checkListPointers(t, l, []*Element[int]{})
	})

	t.Run("move front element; len = 0", func(t *testing.T) {
		l := New[int]()
		e := l.PushFront(0)
		l.MoveToFront(e)
		checkListPointers(t, l, []*Element[int]{e})
	})

	t.Run("move front element; len > 0", func(t *testing.T) {
		l := New[int]()
		e := l.PushFront(0)
		e2 := l.PushFront(0)
		l.MoveToFront(e2)
		checkListPointers(t, l, []*Element[int]{e2, e})
	})

	t.Run("move back to front", func(t *testing.T) {
		l := New[int]()
		e := l.PushFront(0)
		e2 := l.PushFront(0)
		l.MoveToFront(e)
		checkListPointers(t, l, []*Element[int]{e, e2})
	})

	t.Run("move middle to front", func(t *testing.T) {
		l := New[int]()
		e := l.PushFront(0)
		e2 := l.PushFront(0)
		e3 := l.PushFront(0)
		l.MoveToFront(e2)
		checkListPointers(t, l, []*Element[int]{e2, e3, e})
	})
}

func Test_MoveToBack(t *testing.T) {
	t.Run("move element not part of list", func(t *testing.T) {
		l := New[int]()
		e := &Element[int]{}
		l.MoveToBack(e)
		checkListPointers(t, l, []*Element[int]{})
	})

	t.Run("move back element; len = 0", func(t *testing.T) {
		l := New[int]()
		e := l.PushFront(0)
		l.MoveToBack(e)
		checkListPointers(t, l, []*Element[int]{e})
	})

	t.Run("move back element; len > 0", func(t *testing.T) {
		l := New[int]()
		e := l.PushFront(0)
		e2 := l.PushFront(0)
		l.MoveToBack(e)
		checkListPointers(t, l, []*Element[int]{e2, e})
	})

	t.Run("move back to back", func(t *testing.T) {
		l := New[int]()
		e := l.PushFront(0)
		e2 := l.PushFront(0)
		l.MoveToBack(e)
		checkListPointers(t, l, []*Element[int]{e2, e})
	})

	t.Run("move middle to back", func(t *testing.T) {
		l := New[int]()
		e := l.PushFront(0)
		e2 := l.PushFront(0)
		e3 := l.PushFront(0)
		l.MoveToBack(e2)
		checkListPointers(t, l, []*Element[int]{e3, e, e2})
	})
}
