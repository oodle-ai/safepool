# `safepool` simplified usage of sync.Pool

# Pool Usage

Pool is a type safe wrapper around `sync.Pool`

Examples.
```go
import "github.com/oodle-ai/safepool"

// Example 1: Byte buffer.
bufPool := safepool.NewPool[bytes.Buffer]()
buf := bufPool.Get()
defer bufPool.Put(buf)
...
	
// Example 2: List of structs.
listPool := safepool.NewPool[safepool.List[TestStruct]]()
list := listPool.Get()
defer listPool.Put(list)
elem1 := list.AppendNewDirtyElem()
elem1.A = 1
elem2 := list.AppendNewDirtyElem()
elem2.A = 2
...

// Example 3: List of integers.
intListPool := safepool.NewPool[safepool.ListVal[int]]()
intList := intListPool.Get()
defer intListPool.Put(intList)
intList.HardResetAndResizeTo(2)
intList.Elems[0] = 1
intList.Elems[1] = 2

```

# Pool Manager Usage

Pool manager is used to manage lifecycle of objects retrieved from a pool.
The pool manager can be passed to other functions which can get objects from the 
pool.

When the pool manager is released, all objects retrieved from the pool manager are
returned to the pool.

```go
import "github.com/oodle-ai/safepool"

anyListPool := safepool.NewPool[safepool.ListVal[any]]()
pool := safepool.NewPool[bytes.Buffer]()
poolManager := safepool.NewPoolManager[bytes.Buffer](anyListPool, pool)
defer poolManager.ReturnToPool() // Eventually returns all buffers back to the pool.

buf1 := poolManager.Get()
buf1.Reset()

buf2 := poolManager.Get()
buf2.Reset()

buf3 := poolManager.Get()
buf3.Reset()

someFunctionCall(poolManager) // this function can do more poolManager.Get()
```
