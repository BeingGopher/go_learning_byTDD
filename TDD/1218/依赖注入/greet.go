package main

import (
	"fmt"
	"io"
	"net/http"
)

func Greet(writer io.Writer, name string) {
	//fmt.Fprintf 和 fmt.Printf 一样，只不过 fmt.Fprintf 会接收一个 Writer 参数，用于把字符串传递过去，而 fmt.Printf 默认是标准输出。
	fmt.Fprintf(writer, "Hello, %s", name)
}

func MyGreeterHandler(w http.ResponseWriter, r *http.Request) {
	Greet(w, "world") //还可以将数据写入网络
}

// 当你编写一个 HTTP 处理器（handler）时，你需要给出 http.ResponseWriter 和用于创建请求的 http.Request。在你实现服务器时，你使用 writer 写入了请求。
func main() {
	http.ListenAndServe(":5000", http.HandlerFunc(MyGreeterHandler))
}
