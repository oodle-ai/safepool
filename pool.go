package safepool

import "sync"

// Pool is a type safe version of sync.Pool.
// Usage:
//
// Example 1.
//
//	pool := safepool.NewPool[bytes.Buffer]()
//
//	buf := pool.Get()
//	buf.Reset()
//	defer pool.Put(buf)
//	... use buf ...
//
// Example 2.
//
//	pool := safepool.NewPool[safepool.List[CustomStruct]]()
//
//	list := pool.Get()
//	list.Reset()
//	defer pool.Put(list)
//	newElem := list.AppendNewDirtyElem()
//	... Reset newElem
//	... Populate fields of newElem.
type Pool[T any] struct {
	pool *sync.Pool
}

// Get an item from the pool. The item may be dirty.
// The caller should reset the item before using it.
func (p *Pool[T]) Get() *T {
	return p.pool.Get().(*T)
}

// Put an item back to the pool.
func (p *Pool[T]) Put(elem *T) {
	p.pool.Put(elem)
}

// NewPool creates a new pool to create objects of type T.
func NewPool[T any]() *Pool[T] {
	return &Pool[T]{
		pool: &sync.Pool{
			New: func() interface{} {
				return new(T)
			},
		},
	}
}

// NewPoolWithConstructor creates a new pool to create objects of type T
// with a custom constructor.
func NewPoolWithConstructor[T any](constructor func() interface{}) *Pool[T] {
	return &Pool[T]{
		pool: &sync.Pool{
			New: constructor,
		},
	}
}
