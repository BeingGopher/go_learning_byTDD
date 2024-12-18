package main

import (
	"bytes"
	"reflect"
	"testing"
)

func TestCountdown(t *testing.T) {
	//让 Countdown 函数将数据写到某处，io.writer 就是作为 Go 的一个接口来抓取数据的一种方式。
	/*
		在 main 中，我们将信息发送到 os.Stdout，所以用户可以看到 Countdown 的结果打印到终端
		在测试中，我们将发送到 bytes.Buffer，所以我们的测试能够抓取到正在生成的数据*/

	t.Run("sleep after every print", func(t *testing.T) {
		spySleepPrinter := &CountdownOperationsSpy{}
		Countdown(spySleepPrinter, spySleepPrinter)

		want := []string{
			sleep,
			write,
			sleep,
			write,
			sleep,
			write,
			sleep,
			write,
		}

		if !reflect.DeepEqual(want, spySleepPrinter.Calls) {
			t.Errorf("wanted calls %v, got %v", want, spySleepPrinter.Calls)
		}
	})

	t.Run("prints 3 to Go!", func(t *testing.T) {
		buffer := &bytes.Buffer{}
		Countdown(buffer, &CountdownOperationsSpy{})

		got := buffer.String()
		want := `3
2
1
Go!`

		if got != want {
			t.Errorf("got '%s' want '%s'", got, want)
		}
	})
}
