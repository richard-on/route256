// Package linkedlist provides a doubly linked list implementation
// identical to stdlib's container/list.
// This implementation, unlike stdlib's, uses generics.
package linkedlist

// Node is an element of a linked list.
type Node[T any] struct {
	next, prev *Node[T]
	list       *LinkedList[T]
	Payload    T
}

// Next returns the next list node or nil.
func (n *Node[T]) Next() *Node[T] {
	if p := n.next; n.list != nil && p != &n.list.root {
		return p
	}
	return nil
}

// Prev returns the previous list node or nil.
func (n *Node[T]) Prev() *Node[T] {
	if p := n.prev; n.list != nil && p != &n.list.root {
		return p
	}
	return nil
}

// LinkedList represents a doubly linked list.
type LinkedList[T any] struct {
	root Node[T]
	len  int
}

// New returns an initialized list.
func New[T any]() *LinkedList[T] {
	return new(LinkedList[T]).Init()
}

// Len returns the number of elements of list l.
func (l *LinkedList[T]) Len() int {
	return l.len
}

// Init initializes or clears list l.
func (l *LinkedList[T]) Init() *LinkedList[T] {
	l.root.next = &l.root
	l.root.prev = &l.root
	l.len = 0
	return l
}

// Front returns the first element of list l or nil if the list is empty.
func (l *LinkedList[T]) Front() *Node[T] {
	if l.len == 0 {
		return nil
	}
	return l.root.next
}

// Back returns the last element of list l or nil if the list is empty.
func (l *LinkedList[T]) Back() *Node[T] {
	if l.len == 0 {
		return nil
	}
	return l.root.prev
}

// Remove removes e from l if e is an element of list l.
// It returns the element value e.Payload.
// The element must not be nil.
func (l *LinkedList[T]) Remove(n *Node[T]) T {
	if n.list == l {
		// if n.list == l, l must have been initialized when n was inserted
		// in l or l == nil (n is a zero Element) and l.remove will crash
		l.remove(n)
	}
	return n.Payload
}

// PushFront inserts a new element e with value v at the front of list l and returns e.
func (l *LinkedList[T]) PushFront(v T) *Node[T] {
	l.initNext()
	return l.insertValue(v, &l.root)
}

// PushBack inserts a new element e with value v at the back of list l and returns e.
func (l *LinkedList[T]) PushBack(v T) *Node[T] {
	l.initNext()
	return l.insertValue(v, l.root.prev)
}

// InsertBefore inserts a new element e with value v immediately before mark and returns e.
// If mark is not an element of l, the list is not modified.
// The mark must not be nil.
func (l *LinkedList[T]) InsertBefore(v T, mark *Node[T]) *Node[T] {
	if mark.list != l {
		return nil
	}
	// see comment in List.Remove about initialization of l
	return l.insertValue(v, mark.prev)
}

// InsertAfter inserts a new element e with value v immediately after mark and returns e.
// If mark is not an element of l, the list is not modified.
// The mark must not be nil.
func (l *LinkedList[T]) InsertAfter(v T, mark *Node[T]) *Node[T] {
	if mark.list != l {
		return nil
	}
	// see comment in List.Remove about initialization of l
	return l.insertValue(v, mark)
}

// MoveToFront moves element e to the front of list l.
// If e is not an element of l, the list is not modified.
// The element must not be nil.
func (l *LinkedList[T]) MoveToFront(n *Node[T]) {
	if n.list != l || l.root.next == n {
		return
	}
	// see comment in List.Remove about initialization of l
	l.move(n, &l.root)
}

// MoveToBack moves element e to the back of list l.
// If e is not an element of l, the list is not modified.
// The element must not be nil.
func (l *LinkedList[T]) MoveToBack(n *Node[T]) {
	if n.list != l || l.root.prev == n {
		return
	}
	// see comment in List.Remove about initialization of l
	l.move(n, l.root.prev)
}

// MoveBefore moves element e to its new position before mark.
// If e or mark is not an element of l, or e == mark, the list is not modified.
// The element and mark must not be nil.
func (l *LinkedList[T]) MoveBefore(n, mark *Node[T]) {
	if n.list != l || n == mark || mark.list != l {
		return
	}
	l.move(n, mark.prev)
}

// MoveAfter moves element e to its new position after mark.
// If e or mark is not an element of l, or e == mark, the list is not modified.
// The element and mark must not be nil.
func (l *LinkedList[T]) MoveAfter(n, mark *Node[T]) {
	if n.list != l || n == mark || mark.list != l {
		return
	}
	l.move(n, mark)
}

// PushBackList inserts a copy of another list at the back of list l.
// The lists l and other may be the same. They must not be nil.
func (l *LinkedList[T]) PushBackList(other *LinkedList[T]) {
	l.initNext()
	for i, e := other.Len(), other.Front(); i > 0; i, e = i-1, e.Next() {
		l.insertValue(e.Payload, l.root.prev)
	}
}

// PushFrontList inserts a copy of another list at the front of list l.
// The lists l and other may be the same. They must not be nil.
func (l *LinkedList[T]) PushFrontList(other *LinkedList[T]) {
	l.initNext()
	for i, e := other.Len(), other.Back(); i > 0; i, e = i-1, e.Prev() {
		l.insertValue(e.Payload, &l.root)
	}
}

func (l *LinkedList[T]) initNext() {
	if l.root.next == nil {
		l.Init()
	}
}

// insert inserts e after at, increments l.len, and returns e.
func (l *LinkedList[T]) insert(node, at *Node[T]) *Node[T] {
	node.prev = at
	node.next = at.next
	node.prev.next = node
	node.next.prev = node
	node.list = l
	l.len++
	return node
}

// insertValue is a convenience wrapper for insert(&Element{Payload: v}, at).
func (l *LinkedList[T]) insertValue(v T, at *Node[T]) *Node[T] {
	return l.insert(&Node[T]{Payload: v}, at)
}

// remove removes e from its list, decrements l.len
func (l *LinkedList[T]) remove(e *Node[T]) {
	e.prev.next = e.next
	e.next.prev = e.prev
	e.next = nil // avoid memory leaks
	e.prev = nil // avoid memory leaks
	e.list = nil
	l.len--
}

// move moves e to next to at.
func (l *LinkedList[T]) move(e, at *Node[T]) {
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
