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

const write = "write"
const sleep = "sleep"

func Countdown(out io.Writer, sleeper Sleeper) {
	//使用fmt.Fprint传入一个io.Writer（例如 *bytes.Buffer）并发送一个 string。

	for i := countdownStart; i > 0; i-- {
		//向下计数的输出1秒的暂停，Go 可以通过 time.Sleep 实现这个功能
		sleeper.Sleep()
		fmt.Fprintln(out, i)

	}

	sleeper.Sleep()
	fmt.Fprint(out, finalWord)
}

// 将函数应用到 main中。这样的话，我们就有了一些可工作的软件来确保我们的工作正在取得进展。
// 在测试的支持下，将功能切分成小的功能点，并使其首尾相连顺利的运行。

type Sleeper interface {
	Sleep()
}

//监视器（spies）是一种 mock，它可以记录依赖关系是怎样被使用的。

type SpySleeper struct {
	Calls int
}

type ConfigurableSleeper struct {
	duration time.Duration
}

func (o *ConfigurableSleeper) Sleep() {
	time.Sleep(o.duration)
}

func (s *SpySleeper) Sleep() {
	s.Calls++
}

type CountdownOperationsSpy struct {
	Calls []string
}

func (s *CountdownOperationsSpy) Sleep() {
	s.Calls = append(s.Calls, sleep)
}

func (s *CountdownOperationsSpy) Write(p []byte) (n int, err error) {
	s.Calls = append(s.Calls, write)
	return
}

func main() {
	sleeper := &ConfigurableSleeper{1 * time.Second}
	Countdown(os.Stdout, sleeper)
}
