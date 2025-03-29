package heap

/*
堆是一种特殊的完全二叉树结构，它满足堆属性：
即对于任意节点，它的值都大于等于其子节点的值（在最大堆中），或小于等于其子节点的值（在最小堆中）。
相比标准库的"container/heap"包，此版本无需自定义存储结构（默认用切片），只需让元素实现Less方法，更简洁易用。
*/

// 用户数据接口
type Item interface {
    Less(Item) bool
}

// 堆
type Heap struct {
    arr []Item
}

// 创建堆
func New() *Heap {
    return new(Heap)
}

// 获取元素数量
func (h *Heap) Len() int {
    return len(h.arr)
}

// 入堆
func (h *Heap) Push(x Item) {
    h.arr = append(h.arr, x)
    h.up(len(h.arr) - 1)
}

// 出堆
func (h *Heap) Pop() Item {
    if len(h.arr) == 0 {
        return nil
    }
    x := h.arr[0]
    h.arr[0] = h.arr[len(h.arr)-1]
    h.arr = h.arr[:len(h.arr)-1]
    h.down(0)
    return x
}

// 向上调整
func (h *Heap) up(i int) {
    for i > 0 {
        p := (i - 1) / 2 // 父节点
        if !h.arr[i].Less(h.arr[p]) {
            break
        }
        h.arr[i], h.arr[p] = h.arr[p], h.arr[i]
        i = p
    }
}

// 向下调整
func (h *Heap) down(i int) {
    for {
        min := i
        left, right := i*2+1, i*2+2
        if left < len(h.arr) && h.arr[left].Less(h.arr[min]) {
            min = left
        }
        if right < len(h.arr) && h.arr[right].Less(h.arr[min]) {
            min = right
        }
        if min == i {
            break
        }
        // 把最小节点交换上来，迭代向下调整
        h.arr[i], h.arr[min] = h.arr[min], h.arr[i]
        i = min
    }
}
