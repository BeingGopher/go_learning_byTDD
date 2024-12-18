package main

import (
	"bytes"
	"fmt"
)

func Countdown(out *bytes.Buffer) {
	//使用fmt.Fprint传入一个io.Writer（例如 *bytes.Buffer）并发送一个 string。
	fmt.Fprint(out, "3")
}
