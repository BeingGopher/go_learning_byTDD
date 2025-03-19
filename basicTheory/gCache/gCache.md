# gCache

## LRU实现

```go
package main

import (
	"fmt"
)

// Node 定义双向链表节点
type Node struct {
	Key  int
	Val  int
	Prev *Node
	Next *Node
}

// LRUCache 定义 LRU 缓存结构
type LRUCache struct {
	capacity int           // 缓存容量
	cache    map[int]*Node  // 哈希表用于O(1)查找
	head     *Node          // 虚拟头节点
	tail     *Node          // 虚拟尾节点
}

// Constructor 初始化 LRU 缓存
func Constructor(capacity int) LRUCache {
	return LRUCache{
		capacity: capacity,
		cache:    make(map[int]*Node),
		head:     &Node{}, // 虚拟头节点
		tail:     &Node{}, // 虚拟尾节点
		head.Next = tail  // 初始化头尾连接
		tail.Prev = head,
	}
}

// Get 获取键对应的值，并将该键移动到最前面
func (this *LRUCache) Get(key int) int {
	if node, exists := this.cache[key]; exists {
		this.moveToFront(node)
		return node.Val
	}
	return -1
}

// Put 插入或更新键值对，若超出容量则删除最久未使用的元素
func (this *LRUCache) Put(key int, value int) {
	if this.capacity == 0 {
		return
	}

	if node, exists := this.cache[key]; exists {
		node.Val = value
		this.moveToFront(node)
		return
	}

	// 创建新节点
	newNode := &Node{Key: key, Val: value}
	this.cache[key] = newNode
	this.addToFront(newNode)

	// 检查是否超出容量
	if len(this.cache) > this.capacity {
		// 删除最久未使用的节点（虚拟尾节点的前驱）
		last := this.tail.Prev
		this.removeNode(last)
		delete(this.cache, last.Key)
	}
}

// addToFront 将节点添加到链表头部
func (this *LRUCache) addToFront(node *Node) {
	this.removeNode(node)
	node.Next = this.head.Next
	node.Prev = this.head
	this.head.Next.Prev = node
	this.head.Next = node
}

// moveToFront 将节点移动到链表头部
func (this *LRUCache) moveToFront(node *Node) {
	this.removeNode(node)
	this.addToFront(node)
}

// removeNode 删除指定节点
func (this *LRUCache) removeNode(node *Node) {
	prev := node.Prev
	next := node.Next

	prev.Next = next
	next.Prev = prev
}

// 测试用例
func main() {
	cache := Constructor(2)

	cache.Put(1, 1)
	cache.Put(2, 2)
	fmt.Println(cache.Get(1)) // 返回 1，此时缓存顺序为 [1,2]
	cache.Put(3, 3)         // 删除 2，缓存变为 [1,3]
	fmt.Println(cache.Get(2)) // 返回 -1
	fmt.Println(cache.Get(3)) // 返回 3
	cache.Put(4, 4)         // 删除 1，缓存变为 [3,4]
	fmt.Println(cache.Get(1)) // 返回 -1
	fmt.Println(cache.Get(3)) // 返回 3
	fmt.Println(cache.Get(4)) // 返回 4

	// 清空缓存测试
	emptyCache := Constructor(0)
	emptyCache.Put(1, 1) // 无操作
	fmt.Println(emptyCache.Get(1)) // 返回 -1
}
```

