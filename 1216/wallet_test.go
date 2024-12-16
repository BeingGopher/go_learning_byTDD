package main

import (
	"testing"
)

func TestWallet(t *testing.T) {
	//重构
	assertBalance := func(t *testing.T, wallet Wallet, want Bitcoin) {
		got := wallet.Balance()

		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
	}

	wallet := Wallet{}

	wallet.Deposit(Bitcoin(10))

	got := wallet.Balance()
	//通过打印语句调试
	//fmt.Println("address of balance in test is", &wallet.balance)

	want := Bitcoin(10)

	if got != want {
		t.Errorf("got %s want %s", got, want)
	}
	//目前的状态，会报错。当调用一个函数或方法时，参数会被复制。

	t.Run("Deposit", func(t *testing.T) {
		wallet := Wallet{}

		wallet.Deposit(Bitcoin(10))

		assertBalance(t, wallet, Bitcoin(10))
	})

	t.Run("Withdraw", func(t *testing.T) {
		wallet := Wallet{balance: Bitcoin(20)}

		wallet.Withdraw(Bitcoin(10))

		assertBalance(t, wallet, Bitcoin(10))
	})

}
