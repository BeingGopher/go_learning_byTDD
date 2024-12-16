package main

import (
	"testing"
)

func TestWallet(t *testing.T) {
	//重构
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

	t.Run("Withdraw with funds", func(t *testing.T) {
		wallet := Wallet{Bitcoin(20)}

		err := wallet.Withdraw(Bitcoin(10))

		assertBalance(t, wallet, Bitcoin(10))

		assertNoError(t, err)
	})

	t.Run("Withdraw insufficient funds", func(t *testing.T) {

		wallet := Wallet{Bitcoin(20)}
		err := wallet.Withdraw(Bitcoin(100))

		assertBalance(t, wallet, Bitcoin(20))

		assertError(t, err, InsufficientFundsError)
	})

}

func assertBalance(t *testing.T, wallet Wallet, want Bitcoin) {
	got := wallet.Balance()

	if got != want {
		t.Errorf("got '%s' want '%s'", got, want)
	}
}

//没有检查 Withdraw 是否成功

func assertNoError(t *testing.T, got error) {
	if got != nil {
		t.Fatal("got an error but didn't want one")
	}
}

func assertError(t *testing.T, got error, want error) {
	if got == nil {
		t.Fatal("didn't get an error but wanted one")
	}

	if got != want {
		t.Errorf("got '%s' want '%s'", got, want)
	}
}
