package rbtree

/*
红黑树的性质：
1. 节点是红色或黑色
2. 根节点是黑色
3. 叶子节点都是黑色
4. 红色节点的子节点都是黑色
5. 任一节点到其所有叶子节点的路径都包含相同数量的黑色节点

红黑树的平衡性：
任一节点到其所有叶子节点的最长路径不会超过最短路径的2倍。

红黑树的优点：
相比AVL树，红黑树的平衡性要求较弱，插入和删除节点通常需要更少的旋转操作来维持平衡，因此具有更好的性能。
*/

const (
    RED   = true
    BLACK = false
)

// 用户数据接口
type Item interface {
    Less(Item) bool
}

// 存储节点
type Node struct {
    Item
    color               bool
    left, right, parent *Node
}

// 红黑树
type Rbtree struct {
    Nil  *Node
    root *Node
    size int
}

// 创建红黑树
func New() *Rbtree {
    Nil := &Node{color: BLACK}
    return &Rbtree{
        Nil:  Nil,
        root: Nil,
        size: 0,
    }
}

// 获取节点数量
func (t *Rbtree) Len() int {
    return t.size
}

// 左旋转
func (t *Rbtree) leftRotate(x *Node) {
    if x == t.Nil {
        return
    }
    y := x.right
    if y == t.Nil {
        return
    }
    // y顶替x
    p := x.parent
    if p == t.Nil {
        t.root = y
    } else if x == p.left {
        p.left = y
    } else {
        p.right = y
    }
    y.parent = p
    // y左子树移到x右子树
    x.right = y.left
    y.left.parent = x
    // x成为y左子树
    y.left = x
    x.parent = y
}

// 右旋转
func (t *Rbtree) rightRotate(x *Node) {
    if x == t.Nil {
        return
    }
    y := x.left
    if y == t.Nil {
        return
    }
    // y顶替x
    p := x.parent
    if p == t.Nil {
        t.root = y
    } else if x == p.left {
        p.left = y
    } else {
        p.right = y
    }
    y.parent = p
    // y右子树移到x左子树
    x.left = y.right
    y.right.parent = x
    // x成为y右子树
    y.right = x
    x.parent = y
}

// 插入节点
func (t *Rbtree) Insert(item Item) {
    if item == nil {
        return
    }
    // 寻找插入位置
    p := t.root
    for p != t.Nil {
        if item.Less(p.Item) {
            if p.left != t.Nil {
                p = p.left
            } else {
                break
            }
        } else if p.Less(item) {
            if p.right != t.Nil {
                p = p.right
            } else {
                break
            }
        } else {
            // 已存在，更新数据即可
            p.Item = item
            return
        }
    }
    // 插入节点（设置为红色，可以保持性质5）
    x := &Node{item, RED, t.Nil, t.Nil, t.Nil}
    if p == t.Nil {
        t.root = x
    } else if item.Less(p.Item) {
        p.left = x
    } else {
        p.right = x
    }
    x.parent = p
    t.size++
    // 修复
    t.insertFix(x)
}

// 插入修复
func (t *Rbtree) insertFix(x *Node) {
    // 如果是根节点，变成黑色即可
    // 如果父节点是黑色，无需处理
    // 如果父节点是红色，破坏了性质4，分多种情况处理
    for x.parent.color == RED {
        // p, g, u 分别是x的父节点，祖父节点，叔叔节点
        p := x.parent
        g := p.parent
        if p == g.left {
            u := g.right
            if u.color == RED {
                // 叔叔节点是红色。变色+迭代向上处理
                p.color = BLACK
                u.color = BLACK
                g.color = RED
                x = g
            } else {
                if x == p.right {
                    // 叔叔节点是黑色，父子异侧。先旋转到同侧
                    t.leftRotate(p)
                    p = x
                }
                // 叔叔节点是黑色，父子同侧。变色+旋转
                p.color = BLACK
                g.color = RED
                t.rightRotate(g)
                break
            }
        } else {
            // 跟上面情况对称
            u := g.left
            if u.color == RED {
                p.color = BLACK
                u.color = BLACK
                g.color = RED
                x = g
            } else {
                if x == p.left {
                    t.rightRotate(p)
                    p = x
                }
                p.color = BLACK
                g.color = RED
                t.leftRotate(g)
                break
            }
        }
    }
    // 根结点是黑色
    t.root.color = BLACK
}

// 删除节点
func (t *Rbtree) Delete(item Item) {
    if item == nil {
        return
    }
    // 查找要删除的节点y
    y := t.search(item)
    if y == t.Nil {
        return
    }
    if y.left != t.Nil && y.right != t.Nil {
        // 如果y有两个子节点
        // 把y的后继节点z复制到y，然后删除z（转化为要删除的节点至多有一个子节点）
        z := y.right
        for z.left != t.Nil {
            z = z.left
        }
        y.Item = z.Item
        y = z
    }
    // x是y的唯一子节点
    x := y.left
    if x == t.Nil {
        x = y.right
    }
    // 删除y
    p := y.parent
    if p == t.Nil {
        t.root = x
    } else if y == p.left {
        p.left = x
    } else {
        p.right = x
    }
    x.parent = p
    t.size--
    // 如果删除的节点是黑色，破环了性质5，需要修复
    if y.color == BLACK {
        t.deleteFix(x)
    }
}

// 删除修复
func (t *Rbtree) deleteFix(x *Node) {
    for x != t.root {
        // 如果x是红色，变成黑色即可
        if x.color == RED {
            x.color = BLACK
            break
        }
        // 如果x是黑色，分多种情况处理
        // p, b 分别是x的父节点，兄弟节点
        p := x.parent
        if x == p.left {
            b := p.right
            if b.color == RED {
                // 如果兄弟节点是红色。通过变色+旋转，使得兄弟节点为黑色
                b.color = BLACK
                p.color = RED
                t.leftRotate(p)
                b = p.right
            }
            if b.left.color == BLACK && b.right.color == BLACK {
                // 如果兄弟节点的两个子节点都是黑色。
                // 把兄弟节点变成红色以恢复当前子树的平衡，迭代向上处理外部的不平衡
                b.color = RED
                x = p
            } else {
                // 兄弟节点至少有一个红色子节点
                if b.left.color == RED {
                    // 如果异侧是红色子节点。变色+旋转到同侧
                    b.left.color = BLACK
                    b.color = RED
                    t.rightRotate(b)
                    b = p.right
                }
                // 红色子节点在同侧。变色+旋转
                b.color = p.color
                p.color = BLACK
                b.right.color = BLACK
                t.leftRotate(p)
                break
            }
        } else {
            // 跟上面情况对称
            b := p.left
            if b.color == RED {
                b.color = BLACK
                p.color = RED
                t.rightRotate(p)
                b = p.left
            }
            if b.left.color == BLACK && b.right.color == BLACK {
                b.color = RED
                x = p
            } else {
                if b.right.color == RED {
                    b.right.color = BLACK
                    b.color = RED
                    t.leftRotate(b)
                    b = p.left
                }
                b.color = p.color
                p.color = BLACK
                b.left.color = BLACK
                t.rightRotate(p)
                break
            }
        }
    }
    // 根节点是黑色
    t.root.color = BLACK
}

// 查找节点
func (t *Rbtree) Get(item Item) Item {
    if item == nil {
        return nil
    }
    x := t.search(item)
    if x == t.Nil {
        return nil
    }
    return x.Item
}

func (t *Rbtree) search(item Item) *Node {
    x := t.root
    for x != t.Nil {
        if item.Less(x.Item) {
            x = x.left
        } else if x.Less(item) {
            x = x.right
        } else {
            break
        }
    }
    return x
}

// 查找最小节点
func (t *Rbtree) Min() Item {
    x := t.root
    if x == t.Nil {
        return nil
    }
    for x.left != t.Nil {
        x = x.left
    }
    return x.Item
}

// 查找最大节点
func (t *Rbtree) Max() Item {
    x := t.root
    if x == t.Nil {
        return nil
    }
    for x.right != t.Nil {
        x = x.right
    }
    return x.Item
}

type Func func(Item)

// 升序遍历
func (t *Rbtree) Ascend(f Func) {
    t.ascend(t.root, f)
}

func (t *Rbtree) ascend(x *Node, f Func) {
    if x == t.Nil {
        return
    }
    t.ascend(x.left, f)
    f(x.Item)
    t.ascend(x.right, f)
}

// 降序遍历
func (t *Rbtree) Descend(f Func) {
    t.descend(t.root, f)
}

func (t *Rbtree) descend(x *Node, f Func) {
    if x == t.Nil {
        return
    }
    t.descend(x.right, f)
    f(x.Item)
    t.descend(x.left, f)
}
