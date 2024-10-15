package safepool

// List is a wrapper around a slice of pointers to T.
type List[T any] struct {
	Elems []*T
}

// Reset resets the list to 0 length.
func (l *List[T]) Reset() {
	l.Elems = l.Elems[:0]
}

// AppendNewDirtyElem appends and returns a new element. Caller has to
// reset the element before using it.
func (l *List[T]) AppendNewDirtyElem() *T {
	if len(l.Elems) == cap(l.Elems) {
		newElem := new(T)
		l.Elems = append(l.Elems, newElem)
		return newElem
	}

	l.Elems = l.Elems[:len(l.Elems)+1]
	if l.Elems[len(l.Elems)-1] == nil {
		l.Elems[len(l.Elems)-1] = new(T)
	}

	return l.Elems[len(l.Elems)-1]
}

// ListVal is a wrapper around a slice of T.
type ListVal[T any] struct {
	Elems []T
}

// Reset resets the list to 0 length.
func (l *ListVal[T]) Reset() {
	l.Elems = l.Elems[:0]
}

// HardResetAndResizeTo changes the length of the list to n.
// All elements are set the zero value of T.
func (l *ListVal[T]) HardResetAndResizeTo(finalSize int) {
	var zero T
	if finalSize > cap(l.Elems) {
		l.Elems = make([]T, finalSize)
	}

	l.Elems = l.Elems[:finalSize]
	for i := 0; i < finalSize; i++ {
		l.Elems[i] = zero
	}
}
