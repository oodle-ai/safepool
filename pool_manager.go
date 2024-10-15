package safepool

import "sync"

// PoolManager is used to manage life cycle of pool objects.
// Usage:
//
//	anyListPool := safepool.NewPool[safepool.ListVal[any]]()
//	pool := safepool.NewPool[bytes.Buffer]()
//	poolManager := safepool.NewPoolManager[bytes.Buffer](anyListPool, pool)
//	defer poolManager.ReturnToPool() // Returns all buffers back to the pool.
//
//	buf1 := poolManager.Get()
//	buf1.Reset()
//
//	buf2 := poolManager.Get()
//	buf2.Reset()
//
//	buf3 := poolManager.Get()
//	buf3.Reset()
//
//	someFunctionCall(poolManager); // this function can do more poolManager.Get()
type PoolManager[T any] struct {
	sync.Mutex
	anyList     *ListVal[any]
	pool        *Pool[T]
	anyListPool *Pool[ListVal[any]]
}

// NewPoolManager creates a new pool manager to manage objects of type T.
func NewPoolManager[T any](anyListPool *Pool[ListVal[any]], pool *Pool[T]) *PoolManager[T] {
	anyList := anyListPool.Get()
	anyList.Reset()
	return &PoolManager[T]{
		anyList:     anyList,
		pool:        pool,
		anyListPool: anyListPool,
	}
}

// Get an item from the pool. The caller should reset the
// item before using it.
func (p *PoolManager[T]) Get() *T {
	elem := p.pool.Get()
	if elem != nil {
		p.Lock()
		defer p.Unlock()
		p.anyList.Elems = append(p.anyList.Elems, elem)
	}
	return elem
}

// ElemsUnsafe returns the list of elements fetched using Get.
// NOTE: This should not be called concurrently with Get.
// The returned list is not safe to iterate over if Get
// is being called concurrently.
func (p *PoolManager[T]) ElemsUnsafe() *ListVal[any] {
	p.Lock()
	defer p.Unlock()
	return p.anyList
}

// ReturnToPool returns elements fetched using Get to the pool.
// NOTE: This should not be called concurrently with Get.
func (p *PoolManager[T]) ReturnToPool() {
	for _, elem := range p.anyList.Elems {
		p.pool.Put(elem.(*T))
	}

	p.anyListPool.Put(p.anyList)
}
