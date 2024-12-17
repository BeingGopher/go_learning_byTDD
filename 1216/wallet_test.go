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

//当你传值给函数或方法时，Go 会复制这些值。因此，如果你写的函数需要更改状态，你就需要用指针指向你想要更改的值

//Go 取值的副本在大多数时候是有效的，但是有时候你不希望你的系统只使用副本，在这种情况下你需要传递一个引用。例如，非常庞大的数据或者你只想有一个实例（比如数据库连接池）

//指针可以是 nil

//当函数返回一个的指针，你需要确保检查过它是否为 nil，否则你可能会抛出一个执行异常，编译器在这里不能帮到你

//nil 非常适合描述一个可能丢失的值

//在错误中检查字符串会导致测试不稳定。因此，我们用一个有意义的值重构了，这样就更容易测试代码，同时对于我们 API 的用户来说也更简单。
