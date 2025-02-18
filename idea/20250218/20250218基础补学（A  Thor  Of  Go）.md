# 20250218基础补学（A  Thor  Of  Go）

## 裸返回值

Go 的返回值可被命名，它们会被视作定义在函数顶部的变量。

返回值的命名应当能反应其含义，它可以作为文档使用。

没有参数的 `return` 语句会直接返回已命名的返回值，也就是「裸」返回值。

裸返回语句应当仅用在下面这样的短函数中。在长的函数中它们会影响代码的可读性。

```go
package main

import "fmt"

func split(sum int) (x, y int) {
	x = sum * 4 / 9
	y = sum - x
	return
}

func main() {
	fmt.Println(split(17))//7 10
}
```

## 短变量声明

在函数中，短赋值语句 `:=` 可在隐式确定类型的 `var` 声明中使用。

函数外的每个语句都 **必须** 以关键字开始（`var`、`func` 等），因此 `:=` 结构不能在函数外使用。

```go
package main

import "fmt"

func main() {
	var i, j int = 1, 2
	k := 3
	c, python, java := true, false, "no!"

	fmt.Println(i, j, k, c, python, java)//1 2 3 true false no!
}
```

## for循环

初始化语句和后置语句是可选的（分号也可以直接去掉）。

```go
package main

import "fmt"

func main() {
	sum := 1
	for ; sum < 1000; {
		sum += sum
	}
	fmt.Println(sum)//1024
}

package main
//-------------------------------------------
import "fmt"

func main() {
	sum := 1
	for sum < 1000 {
		sum += sum
	}
	fmt.Println(sum)
}
```

## if 和简短语句

和 `for` 一样，`if` 语句可以在条件表达式前执行一个简短语句。

该语句声明的变量作用域仅在 `if` 之内。

```go
package main

import (
	"fmt"
	"math"
)

func pow(x, n, lim float64) float64 {
	if v := math.Pow(x, n); v < lim {
		return v
	}
	return lim
}

func main() {
	fmt.Println(
		pow(3, 2, 10),
		pow(3, 3, 20),
	)//9 20
}
```

## 结构体指针

p指向结构体v（即p是指针，类型为*Vertex）

```go
package main

import "fmt"

type Vertex struct {
	X int
	Y int
}

func main() {
	v := Vertex{1, 2}
	p := &v
	p.X = 1e9
	fmt.Println(v)
}
```

## 切片类似数组的引用

切片就像数组的引用 切片并不存储任何数据，它只是描述了底层数组中的一段。

更改切片的元素会修改其底层数组中对应的元素。

和它共享底层数组的切片都会观测到这些修改。

## nil 切片

切片的零值是 `nil`。

nil 切片的长度和容量为 0 且没有底层数组。

##  向切片追加元素

为切片追加新的元素是种常见的操作，为此 Go 提供了内置的 `append` 函数。内置函数的[文档](https://tour.go-zh.org/pkg/builtin/#append)对该函数有详细的介绍。

```go
func append(s []T, vs ...T) []T
```

`append` 的第一个参数 `s` 是一个元素类型为 `T` 的切片，其余类型为 `T` 的值将会追加到该切片的末尾。

`append` 的结果是一个包含原切片所有元素加上新添加元素的切片。

当 `s` 的底层数组太小，不足以容纳所有给定的值时，它就会分配一个更大的数组。 返回的切片会指向这个新分配的数组。

##  range 遍历

`for` 循环的 `range` 形式可遍历切片或映射。

当使用 `for` 循环遍历切片时，每次迭代都会返回两个值。 第一个值为当前元素的下标，第二个值为该下标所对应元素的一份副本。

> [!WARNING]
>
> 所以如果需要修改切片内容的话，不能通过第二个值，因为第二个值是副本，因该通过下标来修改。

## map 映射

`map` 映射将键映射到值。

映射的零值为 `nil` 。`nil` 映射既没有键，也不能添加键。

`make` 函数会返回给定类型的映射，并将其初始化备用。

在 Go 语言中，闭包是一个非常重要的概念，它可以帮助我们实现一些灵活的功能。让我们逐步解析你提到的这段话，特别是结合一个具体的例子来理解闭包的概念。

## **闭包**

闭包是一个函数值，它引用了其函数体之外的变量。换句话说，闭包不仅包含函数的代码，还包含了函数所引用的外部变量的环境。

### 1. **“绑定”变量的概念**

闭包可以访问并修改它引用的外部变量。即使这个外部变量在函数体之外定义，闭包仍然可以操作它。这种行为就像是闭包被“绑定”到了这些变量上，闭包的生命周期也依赖于这些变量。

```go
package main

import "fmt"

// adder 函数返回一个闭包
func adder() func(int) int {
    sum := 0 // 外部变量
    return func(delta int) int { // 返回的闭包
        sum += delta // 闭包引用并修改外部变量
        return sum
    }
}

func main() {
    // 调用 adder 函数，返回一个闭包
    add := adder()

    fmt.Println(add(1)) // 输出 1，sum = 1
    fmt.Println(add(2)) // 输出 3，sum = 3
    fmt.Println(add(3)) // 输出 6，sum = 6

    // 再次调用 adder 函数，返回一个新的闭包
    add2 := adder()
    fmt.Println(add2(1)) // 输出 1，sum = 1
    fmt.Println(add2(2)) // 输出 3，sum = 3
}
```

### 2. **解释代码**

- **`adder` 函数**：
  - `adder` 函数内部定义了一个变量 `sum`，并返回了一个匿名函数（闭包）。
  - 这个匿名函数可以访问并修改 `sum`，即使 `sum` 是在匿名函数的外部定义的。
- **闭包的“绑定”**：
  - 每次调用 `adder()` 时，都会创建一个新的 `sum` 变量，并返回一个新的闭包。
  - 每个闭包都绑定到它自己的 `sum` 变量上。这意味着每个闭包都有自己的独立状态。
- **输出结果**：
  - 第一次调用 `add(1)`，`sum` 从 0 变为 1，返回 1。
  - 第二次调用 `add(2)`，`sum` 从 1 变为 3，返回 3。
  - 第三次调用 `add(3)`，`sum` 从 3 变为 6，返回 6。
  - 而 `add2` 是另一个闭包，它有自己的 `sum`，从 0 开始计算。

### 3. **总结**

闭包的核心特性是它能够引用并操作外部变量，这些变量的生命周期与闭包绑定在一起。在 `adder` 的例子中，每个闭包都绑定到它自己的 `sum` 变量上，因此每个闭包都有独立的状态。这种特性使得闭包非常适合用于实现一些需要保存状态的场景，比如计数器、累加器等。

## 斐波纳契闭包

```go
package main

import "fmt"

// fibonacci 是返回一个「返回一个 int 的函数」的函数
func fibonacci() func() int {
	// 使用闭包保存状态
	start := 0
	end := 1
	return func() int {
		// 保存当前值
		result := start
		// 更新状态
		start, end = end, start+end
		return result
	}
}

func main() {
	f := fibonacci()
	for i := 0; i < 10; i++ {
		fmt.Println(f())
	}
}
```

## 声明方法

你可以为非结构体类型声明方法。

在此例中，我们看到了一个带 `Abs` 方法的数值类型 `MyFloat`。

你只能为在同一个包中定义的接收者类型声明方法，而不能为其它别的包中定义的类型 （包括 `int` 之类的内置类型）声明方法。

（译注：就是接收者的类型定义和方法声明必须在同一包内。）

##  指针类型的接收者

你可以为指针类型的接收者声明方法。

这意味着对于某类型 `T`，接收者的类型可以用 `*T` 的文法。 （此外，`T` 本身不能是指针，比如不能是 `*int`。）

例如，这里为 `*Vertex` 定义了 `Scale` 方法。

指针接收者的方法可以修改接收者指向的值（如这里的 `Scale` 所示）。 由于方法经常需要修改它的接收者，指针接收者比值接收者更常用。

试着移除第 16 行 `Scale` 函数声明中的 `*`，观察此程序的行为如何变化。

若使用值接收者，那么 `Scale` 方法会对原始 `Vertex` 值的副本进行操作。（对于函数的其它参数也是如此。）`Scale` 方法必须用指针接收者来更改 `main` 函数中声明的 `Vertex` 的值。

> 使用指针接收者的原因有二：
>
> 首先，方法能够修改其接收者指向的值。
>
> 其次，这样可以避免在每次调用方法时复制该值。若值的类型为大型结构体时，这样会更加高效。
>
> 在本例中，`Scale` 和 `Abs` 接收者的类型为 `*Vertex`，即便 `Abs` 并不需要修改其接收者。
>
> 通常来说，所有给定类型的方法都应该有值或指针接收者，但并不应该二者混用。 （我们会在接下来几页中明白为什么。）

## Reader

> 由于每次调用 `Read` 方法都会返回字符 `'A'`，并且不会返回 `io.EOF` 错误，因此可以不断调用 `Read` 方法来获取无限个字符 `'A'`，从而实现了 ASCII 字符 `'A'` 的无限流。

```go
package main

import (
    "golang.org/x/tour/reader"
)

// MyReader 定义一个自定义的读取器类型
type MyReader struct{}

// Read 为 MyReader 类型实现 Read 方法，用于实现 io.Reader 接口
func (r MyReader) Read(b []byte) (int, error) {
    // 检查传入的字节切片长度是否为 0
    if len(b) == 0 {
        return 0, nil
    }
    // 向字节切片的第一个元素写入字符 'A'
    b[0] = 'A'
    // 返回读取的字节数（这里为 1）和 nil 错误
    return 1, nil
}

func main() {
    // 验证 MyReader 是否实现了 io.Reader 接口
    reader.Validate(MyReader{})
}
```

