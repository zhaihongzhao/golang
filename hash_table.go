package hashtable

/*
哈希冲突的解决方法：
1. 链地址法。哈希表的每个槽位维护一个链表，当多个元素被哈希到同一个槽位时，它们会被放在这个槽位的链表中。
2. 开放寻址法。当发生冲突时，算法会寻找下一个空的槽位。常见的探测方式有线性探测、二次探测和双重哈希（使用另一个哈希函数计算步长）。
3. 再哈希法。当发生冲突时，使用不同的哈希函数计算新的哈希值，直到找到空的槽位。

哈希表的扩容：
装载因子 = 元素数量/槽位数量，超过某个阈值需要扩容，重新计算所有元素的哈希值并重新分配到新的槽位。
*/

// 阈值
const threshold = 1.0

// 用户数据接口
type Item interface {
    Equal(Item) bool
    Hash() uint
}

// 存储结点
type Node struct {
    Item
    next *Node
}

// 哈希表
type HashTable struct {
    arr  []*Node
    size int
}

// 创建哈希表
func New() *HashTable {
    return &HashTable{}
}

// 获取元素数量
func (h *HashTable) Len() int {
    return h.size
}

// 插入元素
func (h *HashTable) Set(item Item) {
    if len(h.arr) == 0 {
        h.arr = make([]*Node, 8)
    }
    pos := item.Hash() % uint(len(h.arr))
    node := h.arr[pos]
    for node != nil && !node.Equal(item) {
        node = node.next
    }
    if node != nil {
        // 已存在，更新数据即可
        node.Item = item
    } else {
        // 插入结点（头插法）
        h.arr[pos] = &Node{Item: item, next: h.arr[pos]}
        h.size++
        // 装载因子超过阈值需要扩容
        if float64(h.size)/float64(len(h.arr)) > threshold {
            h.reHash()
        }
    }
}

// 查找元素
func (h *HashTable) Get(item Item) Item {
    if len(h.arr) == 0 {
        return nil
    }
    pos := item.Hash() % uint(len(h.arr))
    node := h.arr[pos]
    for node != nil && !node.Equal(item) {
        node = node.next
    }
    if node == nil {
        return nil
    }
    return node.Item
}

// 删除元素
func (h *HashTable) Delete(item Item) {
    if len(h.arr) == 0 {
        return
    }
    pos := item.Hash() % uint(len(h.arr))
    node := h.arr[pos]
    if node == nil {
        return
    }
    if node.Equal(item) {
        // 删除链表的头结点
        h.arr[pos] = node.next
        h.size--
    } else {
        for ; node.next != nil; node = node.next {
            if node.next.Equal(item) {
                node.next = node.next.next
                h.size--
                break
            }
        }
    }
}

// 遍历
func (h *HashTable) ForEach(f func(Item)) {
    for _, node := range h.arr {
        for node != nil {
            f(node.Item)
            node = node.next
        }
    }
}

// 扩容
func (h *HashTable) reHash() {
    oldArr := h.arr
    h.arr = make([]*Node, len(h.arr)*2)
    for _, node := range oldArr {
        for node != nil {
            // 重新计算哈希值
            pos := node.Hash() % uint(len(h.arr))
            // 插入到新表（头插法）
            next := node.next
            node.next = h.arr[pos]
            h.arr[pos] = node
            node = next
        }
    }
}
