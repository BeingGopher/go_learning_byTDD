package main

import (
	"fmt"
	"io"
	"os"
)

//*bytes.Buffer 可以运行，但最好使用通用接口代替。

func Countdown(out io.Writer) {
	//使用fmt.Fprint传入一个io.Writer（例如 *bytes.Buffer）并发送一个 string。

	for i := 3; i > 0; i-- {
		fmt.Fprintln(out, i)
	}

	fmt.Fprint(out, "Go!")
}

// 将函数应用到 main中。这样的话，我们就有了一些可工作的软件来确保我们的工作正在取得进展。
// 在测试的支持下，将功能切分成小的功能点，并使其首尾相连顺利的运行。
func main() {
	Countdown(os.Stdout)
}
