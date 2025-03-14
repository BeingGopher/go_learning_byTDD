# gCache


## http服务端

分布式缓存需要实现节点间的通信，基于http简单且容易实现。

```go
package geecache

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

const defaultBasePath = "/_geecache/"

// HTTPPool implements PeerPicker for a pool of HTTP peers.
type HTTPPool struct {
	// this peer's base URL, e.g. "https://example.net:8000"
	self     string
	basePath string
}

// NewHTTPPool initializes an HTTP pool of peers.
func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: defaultBasePath,
	}
}
```

- `HTTPPool` 有 2 个参数，一个是 self，用来记录自己的地址，包括主机名/IP 和端口。
- 另一个是 basePath，作为节点间通讯地址的前缀，默认是 `/_geecache/`，那么 http://example.com/_geecache/ 开头的请求，就用于节点间的访问。因为一个主机上还可能承载其他的服务，加一段 Path 是一个好习惯。比如，大部分网站的 API 接口，一般以 `/api` 作为前缀。

接下来，实现最为核心的 `ServeHTTP` 方法。

```go
// Log info with server name
func (p *HTTPPool) Log(format string, v ...interface{}) {//...表示这个参数可以接受零个或多个值，这些值会被存储在一个切片中。
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}

// ServeHTTP handle all http requests
func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		panic("HTTPPool serving unexpected path: " + r.URL.Path)
	}
	p.Log("%s %s", r.Method, r.URL.Path)
	// /<basepath>/<groupname>/<key> required
	parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	groupName := parts[0]
	key := parts[1]

	group := GetGroup(groupName)
	if group == nil {
		http.Error(w, "no such group: "+groupName, http.StatusNotFound)
		return
	}

	view, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(view.ByteSlice())
}
```

- ServeHTTP 的实现逻辑是比较简单的，首先判断访问路径的前缀是否是 `basePath`，不是返回错误。
- 我们约定访问路径格式为 `/<basepath>/<groupname>/<key>`，通过 groupname 得到 group 实例，再使用 `group.Get(key)` 获取缓存数据。
- 最终使用 `w.Write()` 将缓存值作为 httpResponse 的 body 返回。

到这里，HTTP 服务端已经完整地实现了。接下来，我们将在单机上启动 HTTP 服务，使用 curl 进行测试。

> 若缓存数据为空，回调函数会从数据库中获取数据

## 一致性哈希

这是缓存从单节点走向分布式节点的关键

1.分布式缓存中，若节点无缓存值需确定从哪获取数据，如 10 个节点随机选节点获取（含自己），存在缓存效率低、存储空间浪费问题（因可能多次从数据源取相同数据）。解决办法是用 hash 算法使给定 key 每次选同一节点，如自定义将 key 各字符 ASCII 码相加后除以 10 取余的 hash 算法 。

2.简单求取 Hash 值解决了缓存性能的问题，但是没有考虑节点数量变化的场景。假设，移除了其中一台节点，只剩下 9 个，那么之前 `hash(key) % 10` 变成了 `hash(key) % 9`，也就意味着几乎缓存值对应的节点都发生了改变。即几乎所有的缓存值都失效了。节点在接收到对应的请求时，均需要重新去数据源获取数据，容易引起 `缓存雪崩`。

> 缓存雪崩：缓存在同一时刻全部失效，造成瞬时DB请求量大、压力骤增，引起雪崩。常因为缓存服务器宕机，或缓存设置了相同的过期时间引起。

### 一致性哈希算法原理

> ### 基本概念
> • **哈希环**：一致性哈希算法将整个哈希空间组织成一个虚拟的圆环，通常使用 0 到 2^{32}- 1 这个范围来表示（对应一个 32 位的无符号整数空间），这个圆环被称为哈希环。
>
> ### 算法步骤
> 1. **节点映射到哈希环**
>     ◦ 对分布式系统中的每个节点（如服务器）通过特定的哈希函数计算出一个哈希值，然后将这个哈希值映射到哈希环上。例如，对于服务器节点 A，经过哈希函数计算得到哈希值 H(A)$，把 H(A)对应的位置标记为节点 A 在哈希环上的位置。
> 2. **数据映射到哈希环**
>     ◦ 对于需要存储的数据（通常用数据的 key 来标识），同样使用相同的哈希函数计算其哈希值，再将这个哈希值映射到哈希环上。比如数据对象 D 的 key 经过哈希函数计算得到 $H(D)$ ，把 $H(D)$ 对应的位置标记为数据 D 在哈希环上的位置。
> 3. **确定数据存储节点**
>     ◦ 在哈希环上，从数据对象所在位置开始，顺时针方向查找第一个遇到的节点，该节点即为存储此数据的节点。例如，数据 D 在哈希环上的位置确定后，沿着顺时针方向找到第一个碰到的节点，假设是节点 B ，那么数据 D 就会被存储到节点 B 上。
>
> ### 节点变动时的处理
>  • **节点增加**：当有新的节点加入分布式系统时，只需重新分配新节点与其顺时针方向相邻节点之间的那部分数据。例如新加入节点 C ，找到其在哈希环上的位置后，把原来由节点 B 负责存储但位于节点 C 和节点 B 顺时针方向之间的数据迁移到节点 C 上，其他节点的数据不受影响。
>  • **节点删除**：当某个节点从系统中移除时，该节点上存储的数据会被迁移到其顺时针方向的下一个节点。比如节点 B 下线，那么原本存储在节点 B 上的数据会被迁移到顺时针方向的下一个节点（如节点 C ）上。
>
> ### 虚拟节点
> 如果服务器的节点过少，容易引起 key 的倾斜。例如上面例子中的 peer2，peer4，peer6 分布在环的上半部分，下半部分是空的。那么映射到环下半部分的 key 都会被分配给 peer2，key 过度向 peer2 倾斜，缓存节点间负载不均。
>
> 在实际应用中，为了使节点在哈希环上分布更加均匀，避免因节点数量较少导致某些区域数据集中，还会引入虚拟节点的概念。每个物理节点可以对应多个虚拟节点，每个虚拟节点都有自己独立的哈希值并映射到哈希环上 。这样能让数据分布更均衡，进一步提升系统的性能和可靠性 。
>
> 而且实现代价非常小，只需要增加一个字典(map)维护真实节点与虚拟节点的映射关系即可。

### 实现

定义map

```go
package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)
// Hash maps bytes to uint32
type Hash func(data []byte) uint32

// Map constains all hashed keys
type Map struct {
	hash     Hash
	replicas int//虚拟节点倍数
	keys     []int // Sorted,哈希环
	hashMap  map[int]string//键是虚拟节点的哈希值，值是真实节点的名称。
}

// New creates a Map instance
func New(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}
```

> - 定义了函数类型 `Hash`，采取依赖注入的方式，允许用于替换成自定义的 Hash 函数，也方便测试时替换，默认为 `crc32.ChecksumIEEE` 算法。
> - `Map` 是一致性哈希算法的主数据结构，包含 4 个成员变量：Hash 函数 `hash`；虚拟节点倍数 `replicas`；哈希环 `keys`；虚拟节点与真实节点的映射表 `hashMap`，键是虚拟节点的哈希值，值是真实节点的名称。
> - 构造函数 `New()` 允许自定义虚拟节点倍数和 Hash 函数。

```go
// Add adds some keys to the hash.
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key
		}
	}
	sort.Ints(m.keys)
}
```

> - `Add` 函数允许传入 0 或 多个真实节点的名称。
> - 对每一个真实节点 `key`，对应创建 `m.replicas` 个虚拟节点，虚拟节点的名称是：`strconv.Itoa(i) + key`，即通过添加编号的方式区分不同虚拟节点。
> - 使用 `m.hash()` 计算虚拟节点的哈希值，使用 `append(m.keys, hash)` 添加到环上。
> - 在 `hashMap` 中增加虚拟节点和真实节点的映射关系。
> - 最后一步，环上的哈希值排序。
> - 私以为，add方法本质类似于扩容

Get方法

```go
// Get gets the closest item in the hash to the provided key.
func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}

	hash := int(m.hash([]byte(key)))
	// Binary search for appropriate replica.
	idx := sort.Search(len(m.keys), func(i int) bool {//标准库二分查找函数，用于在有序切片中查找满足条件的第一个位置。
		return m.keys[i] >= hash
	})

	return m.hashMap[m.keys[idx%len(m.keys)]]
}
```

> - 第一步，计算 key 的哈希值。
> - 第二步，顺时针找到第一个匹配的虚拟节点的下标 `idx`，从 m.keys 中获取到对应的哈希值。如果 `idx == len(m.keys)`，说明应选择 `m.keys[0]`，因为 `m.keys` 是一个环状结构，所以用取余数的方式来处理这种情况。
> - 第三步，通过 `hashMap` 映射得到真实的节点。

## 分布式节点

### 抽象 PeerPicker

```go
package geecache

// PeerPicker is the interface that must be implemented to locate
// the peer that owns a specific key.
type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

// PeerGetter is the interface that must be implemented by a peer.
type PeerGetter interface {
	Get(group string, key string) ([]byte, error)
}
```

> - 抽象出2个接口，PeerPicker 的 `PickPeer()` 方法用于根据传入的 key 选择相应节点 PeerGetter。
> - 接口 PeerGetter 的 `Get()` 方法用于从对应 group 查找缓存值。PeerGetter 就对应于上述流程中的 HTTP 客户端。

### 节点选择与 HTTP 客户端

```go
type httpGetter struct {
	baseURL string
}

func (h *httpGetter) Get(group string, key string) ([]byte, error) {
	u := fmt.Sprintf(
		"%v%v/%v",
		h.baseURL,
		url.QueryEscape(group),
		url.QueryEscape(key),
	)
	res, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned: %v", res.Status)
	}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %v", err)
	}

	return bytes, nil
}

var _ PeerGetter = (*httpGetter)(nil)//接口断言，确保 httpGetter 实现了 PeerGetter 接口。如果没有实现，编译时会报错。
```

> - baseURL 表示将要访问的远程节点的地址，例如 `http://example.com/_geecache/`。
> - 使用 `http.Get()` 方式获取返回值，并转换为 `[]bytes` 类型。

```go
const (
	defaultBasePath = "/_geecache/"
	defaultReplicas = 50
)
// HTTPPool implements PeerPicker for a pool of HTTP peers.
type HTTPPool struct {
	// this peer's base URL, e.g. "https://example.net:8000"
	self        string
	basePath    string
	mu          sync.Mutex // guards peers and httpGetters
	peers       *consistenthash.Map
	httpGetters map[string]*httpGetter // keyed by e.g. "http://10.0.0.2:8008"
}
```

> - 新增成员变量 `peers`，类型是一致性哈希算法的 `Map`，用来根据具体的 key 选择节点。
> - 新增成员变量 `httpGetters`，映射远程节点与对应的 httpGetter。每一个远程节点对应一个 httpGetter，因为 httpGetter 与远程节点的地址 `baseURL` 有关。

```go
func (p *HTTPPool) Set(peers ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.peers = consistenthash.New(defaultReplicas, nil)
	p.peers.Add(peers...)
	p.httpGetters = make(map[string]*httpGetter, len(peers))
	for _, peer := range peers {
		p.httpGetters[peer] = &httpGetter{baseURL: peer + p.basePath}
	}
}

// PickPeer picks a peer according to key
func (p *HTTPPool) PickPeer(key string) (PeerGetter, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if peer := p.peers.Get(key); peer != "" && peer != p.self {
		p.Log("Pick peer %s", peer)
		return p.httpGetters[peer], true
	}
	return nil, false
}

var _ PeerPicker = (*HTTPPool)(nil)
```

> - `Set()` 方法实例化了一致性哈希算法，并且添加了传入的节点。
> - 并为每一个节点创建了一个 HTTP 客户端 `httpGetter`。
> - `PickerPeer()` 包装了一致性哈希算法的 `Get()` 方法，根据具体的 key，选择节点，返回节点对应的 HTTP 客户端。

至此，HTTPPool 既具备了提供 HTTP 服务的能力，也具备了根据具体的 key，创建 HTTP 客户端从远程节点获取缓存值的能力。

## 防止缓存击穿

> **缓存雪崩**：缓存在同一时刻全部失效，造成瞬时DB请求量大、压力骤增，引起雪崩。缓存雪崩通常因为缓存服务器宕机、缓存的 key 设置了相同的过期时间等引起。
>
> **缓存击穿**：一个存在的key，在缓存过期的一刻，同时有大量的请求，这些请求都会击穿到 DB ，造成瞬时DB请求量大、压力骤增。
>
> **缓存穿透**：查询一个不存在的数据，因为不存在则不会写到缓存中，所以每次都会去请求 DB，如果瞬间流量过大，穿透到 DB，导致宕机。

```go
package singleflight

import "sync"

type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

type Group struct {
	mu sync.Mutex       // protects m
	m  map[string]*call
}
```

> - `call` 代表正在进行中，或已经结束的请求。使用 `sync.WaitGroup` 锁避免重入。
> - `Group` 是 singleflight 的主数据结构，管理不同 key 的请求(call)。

```go
func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}
	c := new(call)
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()

	c.val, c.err = fn()
	c.wg.Done()

	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()

	return c.val, c.err
}
```

> - Do 方法，接收 2 个参数，第一个参数是 `key`，第二个参数是一个函数 `fn`。Do 的作用就是，针对相同的 key，无论 Do 被调用多少次，函数 `fn` 都只会被调用一次，等待 fn 调用结束了，返回返回值或错误。
> - 因为会事先判断是否会存在，存在就直接返回，所以只会返回一个值（无论调用多少次do）

## protobuf通信

- 按照 protobuf 的语法，在 `.proto` 文件中定义数据结构，并使用 `protoc` 生成 Go 代码（`.proto` 文件是跨平台的，还可以生成 C、Java 等其他源码文件）。
- 在项目代码中引用生成的 Go 代码。

```properties
syntax = "proto3";

package geecachepb;

message Request {
  string group = 1;
  string key = 2;
}

message Response {
  bytes value = 1;
}

service GroupCache {
  rpc Get(Request) returns (Response);
}
```

> - `Request` 包含 2 个字段， group 和 cache，这与我们之前定义的接口 `/_geecache/<group>/<name>` 所需的参数吻合。
> - `Response` 包含 1 个字段，bytes，类型为 byte 数组，与之前吻合。

```ba
protoc --go_out=. *.proto
s
geecachepb.pb.go  geecachepb.proto
```

> 生成 `geecache.pb.go`

```go
import pb "geecache/geecachepb"

type PeerGetter interface {
	Get(in *pb.Request, out *pb.Response) error
}
```

> 修改 `peers.go` 中的 `PeerGetter` 接口，参数使用 `geecachepb.pb.go` 中的数据类型。

```go
import (
    // ...
	pb "geecache/geecachepb"
	"github.com/golang/protobuf/proto"
)

func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    // ...
	// Write the value to the response body as a proto message.
	body, err := proto.Marshal(&pb.Response{Value: view.ByteSlice()})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(body)
}

func (h *httpGetter) Get(in *pb.Request, out *pb.Response) error {
	u := fmt.Sprintf(
		"%v%v/%v",
		h.baseURL,
		url.QueryEscape(in.GetGroup()),
		url.QueryEscape(in.GetKey()),
	)
    res, err := http.Get(u)
	// ...
	if err = proto.Unmarshal(bytes, out); err != nil {
		return fmt.Errorf("decoding response body: %v", err)
	}

	return nil
}
```

> 修改 `peers.go` 中的 `PeerGetter` 接口，参数使用 `geecachepb.pb.go` 中的数据类型。
>
> - `ServeHTTP()` 中使用 `proto.Marshal()` 编码 HTTP 响应。
> - `Get()` 中使用 `proto.Unmarshal()` 解码 HTTP 响应。