# 20250227算法留痕+最短路径（dijkstra堆优化+**Bellman_ford**）

## **dijkstra（堆优化版）**

[47.参加科学大会（第六期模拟笔试）](https://kamacoder.com/problempage.php?pid=1047)

对于节点多但边少的稀疏图，通过维护边的方式来求出最短路径（邻接表）

> 由于小顶堆的特性，每次弹出的都是节点中距离最小的，如果出现距离大于当前节点的minDist中的值，说明这是一个旧值，直接跳过即可。

```go
package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"math"
	"os"
)

// 定义一个最小堆结构
type MinHeap struct {
	nodes [][2]int // [距离, 节点]
}

// 实现 heap.Interface 的方法
func (h MinHeap) Len() int {
	return len(h.nodes)
}

func (h MinHeap) Less(i, j int) bool {
	return h.nodes[i][0] < h.nodes[j][0]
}

func (h MinHeap) Swap(i, j int) {
	h.nodes[i], h.nodes[j] = h.nodes[j], h.nodes[i]
}

func (h *MinHeap) Push(x interface{}) {
	h.nodes = append(h.nodes, x.([2]int))
}

func (h *MinHeap) Pop() interface{} {
	n := len(h.nodes)
	x := h.nodes[n-1]
	h.nodes = h.nodes[:n-1]
	return x
}

func main() {
	var N, M int
	scanner := bufio.NewScanner(os.Stdin)

	// 读取第一行的 N 和 M
	scanner.Scan()
	fmt.Sscanf(scanner.Text(), "%d %d", &N, &M)

	// 初始化邻接表
	adj := make([][][]int, N+1)
	for i := range adj {
		adj[i] = make([][]int, 0)
	}

	// 读取边的信息
	for i := 0; i < M; i++ {
		var p1, p2, val int
		scanner.Scan()
		fmt.Sscanf(scanner.Text(), "%d %d %d", &p1, &p2, &val)
		adj[p1] = append(adj[p1], []int{p2, val}) // 无向图可以同时添加 p2 -> p1
	}

	start := 1
	end := N

	// 存储从源点到每个节点的最短距离
	minDist := make([]int, N+1)
	for i := range minDist {
		minDist[i] = math.MaxInt64
	}
	minDist[start] = 0 // 起始点到自身的距离为 0

	// 初始化最小堆
	h := &MinHeap{}
	heap.Init(h)
	heap.Push(h, [2]int{0, start}) // [距离, 节点]

	// Dijkstra 算法主循环
	for h.Len() > 0 {
		// 弹出堆顶元素
		current := heap.Pop(h).([2]int)
		dist, node := current[0], current[1]

		// 如果当前节点的距离大于已记录的最短距离，跳过
		if dist > minDist[node] {
			continue
		}

		// 遍历当前节点的所有邻接点
		for _, edge := range adj[node] {
			nextNode, weight := edge[0], edge[1]
			newDist := dist + weight

			// 如果找到更短的路径，更新最短距离并将其加入堆中
			if newDist < minDist[nextNode] {
				minDist[nextNode] = newDist
				heap.Push(h, [2]int{newDist, nextNode})
			}
		}
	}

	// 输出结果
	if minDist[end] == math.MaxInt64 {
		fmt.Println(-1) // 不能到达终点
	} else {
		fmt.Println(minDist[end]) // 到达终点的最短路径
	}
}
```

定义visited数组的版本(更容易理解)

```go
package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"math"
	"os"
)

// 定义一个最小堆结构
type MinHeap struct {
	nodes [][2]int // [距离, 节点]
}

// 实现 heap.Interface 的方法
func (h MinHeap) Len() int {
	return len(h.nodes)
}

func (h MinHeap) Less(i, j int) bool {
	return h.nodes[i][0] < h.nodes[j][0]
}

func (h MinHeap) Swap(i, j int) {
	h.nodes[i], h.nodes[j] = h.nodes[j], h.nodes[i]
}

func (h *MinHeap) Push(x interface{}) {
	h.nodes = append(h.nodes, x.([2]int))
}

func (h *MinHeap) Pop() interface{} {
	n := len(h.nodes)
	x := h.nodes[n-1]
	h.nodes = h.nodes[:n-1]
	return x
}

func main() {
	var N, M int
	scanner := bufio.NewScanner(os.Stdin)

	// 读取第一行的 N 和 M
	scanner.Scan()
	fmt.Sscanf(scanner.Text(), "%d %d", &N, &M)

	// 初始化邻接表
	adj := make([][][]int, N+1)
	for i := range adj {
		adj[i] = make([][]int, 0)
	}

	// 读取边的信息
	for i := 0; i < M; i++ {
		var p1, p2, val int
		scanner.Scan()
		fmt.Sscanf(scanner.Text(), "%d %d %d", &p1, &p2, &val)
		adj[p1] = append(adj[p1], []int{p2, val}) // 无向图可以同时添加 p2 -> p1
	}

	start := 1
	end := N

	// 存储从源点到每个节点的最短距离
	minDist := make([]int, N+1)
	for i := range minDist {
		minDist[i] = math.MaxInt64
	}
	minDist[start] = 0 // 起始点到自身的距离为 0

	// 记录节点是否被访问过
	visited := make([]bool, N+1)

	// 初始化最小堆
	h := &MinHeap{}
	heap.Init(h)
	heap.Push(h, [2]int{0, start}) // [距离, 节点]

	// Dijkstra 算法主循环
	for h.Len() > 0 {
		// 弹出堆顶元素
		current := heap.Pop(h).([2]int)
		dist, node := current[0], current[1]

		// 如果当前节点已经被访问过，跳过
		if visited[node] {
			continue
		}

		// 标记当前节点为已访问
		visited[node] = true

		// 遍历当前节点的所有邻接点
		for _, edge := range adj[node] {
			nextNode, weight := edge[0], edge[1]
			newDist := dist + weight

			// 如果找到更短的路径，更新最短距离并将其加入堆中
			if newDist < minDist[nextNode] {
				minDist[nextNode] = newDist
				heap.Push(h, [2]int{newDist, nextNode})
			}
		}
	}

	// 输出结果
	if minDist[end] == math.MaxInt64 {
		fmt.Println(-1) // 不能到达终点
	} else {
		fmt.Println(minDist[end]) // 到达终点的最短路径
	}
}
```

## Bellman_ford

[94.城市间货物运输 I](https://kamacoder.com/problempage.php?pid=1152)

松弛：即求当前节点的minDist和相邻节点minDist加上节点间距离的最小值

```go
minDist[B] = min(minDist[A] + value, minDist[B])
```

**对所有边松弛一次，相当于计算 起点到达 与起点一条边相连的节点 的最短距离**。那对所有边松弛两次 可以得到与起点 两条边相连的节点的最短距离......

对所有边松弛 n-1 次 就一定能得到 起点到达 终点的最短距离。

```go
package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func main() {
	var n, m int
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	fmt.Sscanf(scanner.Text(), "%d %d", &n, &m)

	// 初始化图的边
	edges := make([][]int, m)
	for i := 0; i < m; i++ {
		var s, t, v int
		scanner.Scan()
		fmt.Sscanf(scanner.Text(), "%d %d %d", &s, &t, &v)
		edges[i] = []int{s, t, v}
	}

	// Bellman-Ford 算法
	start := 1
	end := n
	minDist := make([]int, n+1)
	for i := range minDist {
		minDist[i] = math.MaxInt32
	}
	minDist[start] = 0

	// 放松所有边 n-1 次
	for i := 1; i <= n-1; i++ {
		for _, edge := range edges {
			s, t, v := edge[0], edge[1], edge[2]
			if minDist[s] != math.MaxInt32 && minDist[s]+v < minDist[t] {
				minDist[t] = minDist[s] + v
			}
		}
	}


	// 输出结果
	if minDist[end] == math.MaxInt32 {
		fmt.Println("unconnected")
	} else {
		fmt.Println(minDist[end])
	}
}
```

> 自己在处理节点距离更新的时候有些问题，题目中给的是节点s到节点t的权值为v，应该使用minDist[s]+v < minDist[t]！！！



> [!IMPORTANT]
>
> 个人总结：
>
> 如果是维护节点，通常需要找到特殊的节点（比如入度为0或者入度为2，一个表示是起点，一个表示有环）
>
> 如果是维护边，通常顺序遍历即可

