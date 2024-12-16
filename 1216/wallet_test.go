package main

import (
	"fmt"
	"testing"
)

func TestWallet(t *testing.T) {

	wallet := Wallet{}

	wallet.Deposit(10)

	got := wallet.Balance()
	//通过打印语句调试
	fmt.Println("address of balance in test is", &wallet.balance)

	want := 10

	if got != want {
		t.Errorf("got %d want %d", got, want)
	}
	//目前的状态，会报错。当调用一个函数或方法时，参数会被复制。

}
