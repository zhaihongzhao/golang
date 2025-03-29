package rbtree

import (
	"math/rand"
	"testing"
)

type myItem struct {
	key   int
	value int
}

func (x *myItem) Less(than Item) bool {
	return x.key < than.(*myItem).key
}

func TestRbt(t *testing.T) {
	rbt := New()
	// 插入节点
	const n = 1000
	for i := 0; i < n; i++ {
		item := &myItem{rand.Intn(n), rand.Intn(n)}
		rbt.Insert(item)
	}
	if !isRbt(rbt) {
		t.Fail()
	}
	// 查找最小、最大节点
	minItem := rbt.Min().(*myItem)
	maxItem := rbt.Max().(*myItem)
	t.Logf("min: %d, max: %d\n", minItem.key, maxItem.key)
	// 升序遍历
	var pre Item
	rbt.Ascend(func(x Item) {
		item, ok := x.(*myItem)
		if !ok {
			t.Fail()
		}
		t.Logf("key: %d, value: %d\n", item.key, item.value)
		if pre != nil && !pre.Less(x) {
			t.Fail()
		}
		pre = x
	})
	// 降序遍历
	pre = nil
	rbt.Descend(func(x Item) {
		if pre != nil && !x.Less(pre) {
			t.Fail()
		}
		pre = x
	})
	// 查找、删除节点
	for i := 0; i < n; i++ {
		x := &myItem{key: i}
		item := rbt.Get(x)
		if item != nil {
			rbt.Delete(x)
			item = rbt.Get(x)
			if item != nil {
				t.Fail()
			}
			if !isRbt(rbt) {
				t.Fail()
			}
		}
	}
	if rbt.Len() != 0 {
		t.Fail()
	}
}

// 判断是不是红黑树
func isRbt(t *Rbtree) bool {
	if t.root == t.Nil {
		return true
	}
	if t.root.color != BLACK {
		return false
	}
	if !check4(t) || !check5(t) {
		return false
	}
	return true
}

// 检查性质4：红色节点的子节点都是黑色
func check4(t *Rbtree) bool {
	return check4_(t, t.root)
}

func check4_(t *Rbtree, x *Node) bool {
	if x == t.Nil {
		return true
	}
	if x.color == RED {
		if x.left.color != BLACK || x.right.color != BLACK {
			return false
		}
	}
	return check4_(t, x.left) && check4_(t, x.right)
}

// 检查性质5：任一节点到其所有叶子节点的路径都包含相同数量的黑色节点
func check5(t *Rbtree) bool {
	num := 0
	for x := t.root; x != t.Nil; x = x.left {
		if x.color == BLACK {
			num++
		}
	}
	return check5_(t, t.root, num)
}

func check5_(t *Rbtree, x *Node, num int) bool {
	if x == t.Nil {
		return num == 0
	}
	if x.color == BLACK {
		num--
	}
	return check5_(t, x.left, num) && check5_(t, x.right, num)
}
