package main

import (
	"bytes"
	"testing"
)

func TestGreet(t *testing.T) {
	buffer := bytes.Buffer{} //bytes 包中的 buffer 类型实现了 Writer 接口。
	Greet(&buffer, "蔡徐坤")

	got := buffer.String()
	want := "Hello, 蔡徐坤"

	if got != want {
		t.Errorf("got '%s' want '%s'", got, want)
	}
}
