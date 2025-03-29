package heap

import (
	"math/rand"
	"testing"
)

type myItem struct {
	priority int
	value    int
}

func (x *myItem) Less(than Item) bool {
	return x.priority > than.(*myItem).priority
}

func TestHeap(t *testing.T) {
	h := New()
	const n = 1000
	// 入堆
	for i := 0; i < n; i++ {
		x := &myItem{rand.Intn(n), rand.Intn(n)}
		h.Push(x)
	}
	// 出堆
	res := make([]*myItem, 0, n)
	for h.Len() > 0 {
		x, ok := h.Pop().(*myItem)
		if !ok {
			t.Fail()
		}
		res = append(res, x)
	}
	// 检验结果
	if len(res) != n {
		t.Fail()
	}
	for i := 1; i < n; i++ {
		if res[i].Less(res[i-1]) {
			t.Fail()
		}
	}
}
