package hashtable

import (
	"math/rand"
	"testing"
)

type myItem struct {
	key   int
	value int
}

func (x *myItem) Equal(to Item) bool {
	return x.key == to.(*myItem).key
}

func (x *myItem) Hash() uint {
	return uint(x.key)
}

func TestHashTable(t *testing.T) {
	ht := New()
	// 插入
	const n = 1000
	for i := 0; i < n; i++ {
		x := &myItem{rand.Intn(n), rand.Intn(n)}
		ht.Set(x)
	}
	t.Logf("size: %d, slot: %d\n", ht.size, len(ht.arr))
	if float64(ht.size)/float64(len(ht.arr)) > threshold {
		t.Fail()
	}
	// 遍历
	ht.ForEach(func(x Item) {
		item, ok := x.(*myItem)
		if !ok {
			t.Fail()
		}
		t.Logf("key: %d, value: %d\n", item.key, item.value)
	})
	// 查找，删除
	for i := 0; i < n; i++ {
		x := &myItem{key: i}
		item := ht.Get(x)
		if item != nil {
			ht.Delete(x)
			item = ht.Get(x)
			if item != nil {
				t.Fail()
			}
		}
	}
	if ht.size != 0 {
		t.Fail()
	}
}
