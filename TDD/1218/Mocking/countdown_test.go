package main

import (
	"bytes"
	"testing"
)

func TestCountdown(t *testing.T) {
	//让 Countdown 函数将数据写到某处，io.writer 就是作为 Go 的一个接口来抓取数据的一种方式。
	/*
		在 main 中，我们将信息发送到 os.Stdout，所以用户可以看到 Countdown 的结果打印到终端
		在测试中，我们将发送到 bytes.Buffer，所以我们的测试能够抓取到正在生成的数据*/
	buffer := &bytes.Buffer{}
	spySleeper := &SpySleeper{}

	Countdown(buffer, spySleeper)

	got := buffer.String()
	//反引号语法是创建 string 的另一种方式，但是允许你放置东西例如放到新的一行
	want := `3
2
1
Go!`

	if got != want {
		t.Errorf("got '%s' want '%s'", got, want)
	}

	if spySleeper.Calls != 4 {
		t.Errorf("not enough calls to sleeper, want 4 got %d", spySleeper.Calls)
	}
}
