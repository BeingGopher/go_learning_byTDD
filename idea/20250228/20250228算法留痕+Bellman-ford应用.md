# 20250228算法留痕+Bellman-ford应用

## SPFA算法

[94.城市间货物运输 I](https://kamacoder.com/problempage.php?pid=1152)

优化思路：**只需要对 上一次松弛的时候更新过的节点作为出发节点所连接的边 进行松弛就够了**。

> 实践过程中自己忽略了visited数组的更新逻辑，只要加入队列，就标记以访问，只要取出，就标记未访问。
>
> 这样做是为了保证队列中的不重复加入，队列外的元素可以重复加入。

```go
package main

import (
	"bufio"
	"container/list"
	"fmt"
	"math"
	"os"
)

type Edge struct {
	Target int
	Weight int
}

func main() {
	var n, m int
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	fmt.Sscanf(scanner.Text(), "%d %d", &n, &m)

	edges := make([][]Edge, n+1)
	for i := 0; i < m; i++ {
		var s, t, v int
		scanner.Scan()
		fmt.Sscanf(scanner.Text(), "%d %d %d", &s, &t, &v)
		edges[s] = append(edges[s], Edge{Target: t, Weight: v})
	}

	start := 1
	end := n
	minDist := make([]int, n+1)
	for i := range minDist {
		minDist[i] = math.MaxInt64
	}
	minDist[start] = 0
	visited := make([]bool, n+1)
	q := list.New()
	q.PushBack(start)

	for q.Len() > 0 {
		e := q.Remove(q.Front()).(int)//已经在队列里的，不重复添加，其他的需要重复添加
		visited[e] = false // 标记为未访问，以便后续可能再次入队
		for _, v := range edges[e] {
			if minDist[e]+v.Weight < minDist[v.Target] {
				minDist[v.Target] = minDist[e] + v.Weight
				if !visited[v.Target] {
					visited[v.Target] = true
					q.PushBack(v.Target)
				}
			}
		}
	}

	// 输出结果
	if minDist[end] == math.MaxInt64 {
		fmt.Println("unconnected")
	} else {
		fmt.Println(minDist[end])
	}
}
```

> [!WARNING]
>
> 本题无负权回路，所以不会陷入死循环

##  bellman_ford之判断负权回路

[95.城市间货物运输 II](https://kamacoder.com/problempage.php?pid=1152)

本体有负权回路，不做处理会陷入死循环，这里使用SPFA算法，同时加入访问节点的计数器，如果访问节点次数超过节点数-1，说明存在死循环。

在极端情况下，即：所有节点都与其他节点相连，每个节点的入度为 n-1 （n为节点数量），所以每个节点最多加入 n-1 次队列。

那么如果节点加入队列的次数 超过了 n-1次 ，那么该图就一定有负权回路。

```go
package main

import (
	"bufio"
	"container/list"
	"fmt"
	"math"
	"os"
)

type Edge struct {
	Target int
	Weight int
}

func main() {
	var n, m int
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	fmt.Sscanf(scanner.Text(), "%d %d", &n, &m)

	edges := make([][]Edge, n+1)
	for i := 0; i < m; i++ {
		var s, t, v int
		scanner.Scan()
		fmt.Sscanf(scanner.Text(), "%d %d %d", &s, &t, &v)
		edges[s] = append(edges[s], Edge{Target: t, Weight: v})
	}

	start := 1
	end := n
	minDist := make([]int, n+1)
	for i := range minDist {
		minDist[i] = math.MaxInt64
	}
	minDist[start] = 0
	visited := make([]bool, n+1)
	q := list.New()
	q.PushBack(start)
    countVisisted := make([]int, n+1)
    countVisisted[start]++
    flag := false
	for q.Len() > 0 {
		e := q.Remove(q.Front()).(int)//已经在队列里的，不重复添加，其他的需要重复添加
		visited[e] = false // 标记为未访问，以便后续可能再次入队
		for _, v := range edges[e] {
			if minDist[e]+v.Weight < minDist[v.Target] {
				minDist[v.Target] = minDist[e] + v.Weight
				if !visited[v.Target] {
					visited[v.Target] = true
					q.PushBack(v.Target)
					countVisisted[v.Target]++
					if countVisisted[v.Target] > n-1 {
					    q.Remove(q.Front())//这里通过影响q的长度，间接影响外层循环的退出，但执行依赖队列长度，逻辑不清晰
					    flag = true
					    break
					}
				}
			}
		}
	}

	// 输出结果
	if flag {
	    fmt.Println("circle")
	} else if minDist[end] == math.MaxInt64 {
		fmt.Println("unconnected")
	} else {
		fmt.Println(minDist[end])
	}
	
}
```

> [!IMPORTANT]
>
> q.Remove(q.Front())//这里通过影响q的长度，间接影响外层循环的退出，但执行依赖队列长度，逻辑不清晰

以下是显式地退出外层循环，代码健壮性更好，逻辑更清晰：

```go
package main

import (
	"bufio"
	"container/list"
	"fmt"
	"math"
	"os"
)

type Edge struct {
	Target int
	Weight int
}

func main() {
	var n, m int
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	fmt.Sscanf(scanner.Text(), "%d %d", &n, &m)

	edges := make([][]Edge, n+1)
	for i := 0; i < m; i++ {
		var s, t, v int
		scanner.Scan()
		fmt.Sscanf(scanner.Text(), "%d %d %d", &s, &t, &v)
		edges[s] = append(edges[s], Edge{Target: t, Weight: v})
	}

	start := 1
	end := n
	minDist := make([]int, n+1)
	for i := range minDist {
		minDist[i] = math.MaxInt64
	}
	minDist[start] = 0
	visited := make([]bool, n+1)
	q := list.New()
	q.PushBack(start)
    countVisisted := make([]int, n+1)
    countVisisted[start]++
    flag := false
	for q.Len() > 0 {
		e := q.Remove(q.Front()).(int)//已经在队列里的，不重复添加，其他的需要重复添加
		visited[e] = false // 标记为未访问，以便后续可能再次入队
		for _, v := range edges[e] {
			if minDist[e]+v.Weight < minDist[v.Target] {
				minDist[v.Target] = minDist[e] + v.Weight
				if !visited[v.Target] {
					visited[v.Target] = true
					q.PushBack(v.Target)
					countVisisted[v.Target]++
					if countVisisted[v.Target] > n-1 {
					    flag = true
					    break
					}
				}
				
			}
		}
		if flag {//退出外层循环，代码更清晰
		    break
		}
	}

	// 输出结果
	if flag {
	    fmt.Println("circle")
	} else if minDist[end] == math.MaxInt64 {
		fmt.Println("unconnected")
	} else {
		fmt.Println(minDist[end])
	}
	
}
```

## bellman_ford之单源有限最短路

[96.城市间货物运输 III](https://kamacoder.com/problempage.php?pid=1154)

限制k个节点，说明最多做k+1次松弛。

> ### **`visited` 的作用：**
>
> 1. **避免重复加入队列**：
>    - 在一次松弛操作中，如果某个节点 `to` 已经被加入队列，则不需要再次加入。
>    - 这是因为，在一次松弛操作中，每个节点只需要被处理一次。
>
> ------
>
> ### **`visited` 的更新逻辑：**
>
> 1. **检查是否需要更新最短距离**：
>    - 如果 `minDist[to] > minDistCopy[from] + price`，说明通过 `from` 到 `to` 的路径比当前记录的最短路径更短，需要更新 `minDist[to]`。
> 2. **检查 `visited[to]` 的状态**：
>    - 如果 `visited[to]` 为 `true`，说明节点 `to` 已经在本轮松弛操作中被加入队列，不需要再次加入，直接跳过。
>    - 如果 `visited[to]` 为 `false`，说明节点 `to` 尚未被加入队列，可以继续处理。
> 3. **更新 `visited[to]` 并加入队列**：
>    - 将 `visited[to]` 设置为 `true`，标记节点 `to` 已经被加入队列。
>    - 将节点 `to` 加入队列 `que`，以便在后续的松弛操作中处理。
>
> ------
>
> ### **为什么 `visited` 需要在本轮松弛操作中重置？**
>
> - `visited` 数组的作用范围是 **当前一轮松弛操作**。
> - 在每一轮松弛操作开始时，`visited` 会被重置为 `false`，因为每个节点在每一轮中都需要重新被处理。
> - 如果 `visited` 没有被重置，可能会导致某些节点在后续轮次中无法被正确处理。

> ### **为什么需要 `queSize`？**(BFS逻辑的关键)
>
> 1. **区分不同轮次的节点**：
>    - 在每一轮松弛操作中，我们只处理 **当前队列中的节点**，这些节点是通过上一轮松弛操作加入队列的。
>    - 如果没有 `queSize`，新加入队列的节点会立即被处理，导致无法区分当前轮次和下一轮次的节点。
> 2. **确保每一轮只处理当前队列中的节点**：
>    - `queSize` 记录了当前队列的长度，即当前轮次需要处理的节点数量。
>    - 在循环中，`queSize--` 确保只处理当前轮次的节点，而不会处理后续加入队列的节点。
> 3. **维护算法的正确性**：
>    - 在 Bellman-Ford 算法或其变种中，每一轮松弛操作对应的是“经过 `k` 条边”的最短路径。
>    - 如果没有 `queSize`，可能会导致某些节点在错误的轮次中被处理，从而影响算法的正确性。

> ### **为什么需要复制 `minDist`？**
>
> 1. **确保每一轮松弛操作基于上一轮的结果**：
>    - 在 Bellman-Ford 算法或其变种中，每一轮松弛操作的目标是基于 **上一轮的结果** 更新最短路径。
>    - 如果不复制 `minDist`，当前轮次中某个节点的更新可能会影响后续节点的更新，导致错误的结果。
> 2. **避免当前轮次中更新的干扰**：
>    - 如果在当前轮次中直接使用 `minDist`，某个节点的更新可能会被后续节点使用，从而导致 **同一轮次中多次更新**，违反了 Bellman-Ford 算法的基本规则。
> 3. **维护算法的正确性**：
>    - Bellman-Ford 算法的核心是每一轮松弛操作只更新 **经过 `k` 条边** 的最短路径。
>    - 如果不复制 `minDist`，可能会导致某些节点的路径长度被错误地更新为 **经过 `k+1` 条边** 的结果。

```go
package main

import (
	"bufio"
	"container/list"
	"fmt"
	"math"
	"os"
)

type Edge struct {
	to  int
	val int
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)


	scanner.Scan()
	var n, m int
	fmt.Sscanf(scanner.Text(), "%d %d", &n, &m)


	grid := make([]*list.List, n+1)
	for i := 0; i <= n; i++ {
		grid[i] = list.New()
	}


	for i := 0; i < m; i++ {
		scanner.Scan()
		var p1, p2, val int
		fmt.Sscanf(scanner.Text(), "%d %d %d", &p1, &p2, &val)
		grid[p1].PushBack(Edge{to: p2, val: val})
	}


	scanner.Scan()
	var start, end, k int
	fmt.Sscanf(scanner.Text(), "%d %d %d", &start, &end, &k)
	k++//k+1次松弛


	minDist := make([]int, n+1)
	minDistCopy := make([]int, n+1)
	for i := range minDist {
		minDist[i] = math.MaxInt32
	}
	minDist[start] = 0


	que := list.New()
	que.PushBack(start)


	for k > 0 && que.Len() > 0 {
		visited := make([]bool, n+1)
		copy(minDistCopy, minDist)//要根据上一次的minDist的值来更新
		queSize := que.Len()

		for queSize > 0 {
			node := que.Front().Value.(int)
			que.Remove(que.Front())
			queSize--

			for e := grid[node].Front(); e != nil; e = e.Next() {
				edge := e.Value.(Edge)
				from := node
				to := edge.to
				price := edge.val

				if minDist[to] > minDistCopy[from]+price {
					minDist[to] = minDistCopy[from] + price
					if visited[to] {
						continue
					}
					visited[to] = true
					que.PushBack(to)
				}
			}
		}
		k--
	}

	// 输出结果
	if minDist[end] == math.MaxInt32 {
		fmt.Println("unreachable")
	} else {
		fmt.Println(minDist[end])
	}
}
```

