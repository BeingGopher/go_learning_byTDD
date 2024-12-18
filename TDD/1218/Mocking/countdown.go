package main

import (
	"fmt"
	"io"
	"os"
	"time"
)

//*bytes.Buffer 可以运行，但最好使用通用接口代替。

// 将变量重构为命名变量
const finalWord = "Go!"
const countdownStart = 3

func Countdown(out io.Writer) {
	//使用fmt.Fprint传入一个io.Writer（例如 *bytes.Buffer）并发送一个 string。

	for i := countdownStart; i > 0; i-- {
		//向下计数的输出1秒的暂停，Go 可以通过 time.Sleep 实现这个功能
		time.Sleep(1 * time.Second)
		fmt.Fprintln(out, i)
	}
	time.Sleep(1 * time.Second)
	fmt.Fprint(out, finalWord)
}

// 将函数应用到 main中。这样的话，我们就有了一些可工作的软件来确保我们的工作正在取得进展。
// 在测试的支持下，将功能切分成小的功能点，并使其首尾相连顺利的运行。
func main() {
	Countdown(os.Stdout)
}
