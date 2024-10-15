package safepool

import (
	"bytes"
	"sync"
	"testing"
)

func TestPoolManager(t *testing.T) {
	pool := NewPool[bytes.Buffer]()
	anyListPool := NewPool[ListVal[any]]()

	for i := 0; i < 10; i++ {
		poolManager := NewPoolManager[bytes.Buffer](anyListPool, pool)
		defer poolManager.ReturnToPool()

		buf1 := poolManager.Get()
		buf1.Reset()
		buf1.WriteString("foo")

		buf2 := poolManager.Get()
		buf2.Reset()
		buf2.WriteString("bar")

		buf3 := poolManager.Get()
		buf3.Reset()
		buf3.WriteString("baz")

		done := sync.WaitGroup{}
		done.Add(10)
		for j := 0; j < 10; j++ {
			go func() {
				buf := poolManager.Get()
				buf.Reset()
				done.Done()
			}()
		}

		done.Wait()

		assertEqual(t, "foo", buf1.String())
		assertEqual(t, "bar", buf2.String())
		assertEqual(t, "baz", buf3.String())
	}
}
