package safepool

import (
	"bytes"
	"reflect"
	"testing"
)

func assertEqual[T any](t *testing.T, expected, actual T) {
	if !reflect.DeepEqual(expected, actual) {
		t.Logf("Expected %v, actual %v", expected, actual)
		t.Fail()
	}
}

func TestPool(t *testing.T) {
	pool := NewPool[bytes.Buffer]()
	buf := pool.Get()
	buf.Reset()
	defer pool.Put(buf)
	buf.WriteString("foo")

	buf2 := pool.Get()
	buf2.Reset()
	buf2.WriteString("bar")
	defer pool.Put(buf2)

	assertEqual(t, "foo", buf.String())
	assertEqual(t, "bar", buf2.String())
}

func TestPoolWithConstructor(t *testing.T) {
	pool := NewPoolWithConstructor[bytes.Buffer](func() interface{} {
		buf := new(bytes.Buffer)
		buf.Grow(1024)
		return buf
	})

	buf := pool.Get()
	defer pool.Put(buf)
	assertEqual(t, 1024, buf.Cap())

	buf2 := pool.Get()
	defer pool.Put(buf2)
	assertEqual(t, 1024, buf.Cap())
}

type testStruct struct {
	A int
}

func TestList(t *testing.T) {
	pool := NewPool[List[testStruct]]()
	list := pool.Get()
	defer pool.Put(list)

	e1 := list.AppendNewDirtyElem()
	e1.A = 1

	e2 := list.AppendNewDirtyElem()
	e2.A = 2

	assertEqual(t, 2, len(list.Elems))
	assertEqual(t, 1, list.Elems[0].A)
	assertEqual(t, 2, list.Elems[1].A)
}

func TestListVal(t *testing.T) {
	pool := NewPool[ListVal[testStruct]]()
	list := pool.Get()
	defer pool.Put(list)

	list.HardResetAndResizeTo(2)
	list.Elems[0].A = 1
	list.Elems[1].A = 2

	assertEqual(t, 2, len(list.Elems))
	assertEqual(t, 1, list.Elems[0].A)
	assertEqual(t, 2, list.Elems[1].A)
}
