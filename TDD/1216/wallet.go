package main

import (
	"errors"
	"fmt"
)

type Bitcoin int

type Wallet struct {
	//在 Go 中，如果一个符号（例如变量、类型、函数等）是以小写符号开头，那么它在 定义它的包之外 就是私有的。
	balance Bitcoin
}

func (w *Wallet) Deposit(amount Bitcoin) {
	//使用「receiver」变量访问结构体内部的 balance 字段。
	//fmt.Println("address of balance in deposit is", &w.balance)
	w.balance += amount
}

// 用指针来解决这个问题。指针让我们指向某个值，然后修改它。(添加*)

// 指向wallet的指针

func (w *Wallet) Balance() Bitcoin {
	return w.balance
}

type Stringer interface {
	String() string
}

func (b Bitcoin) String() string {
	return fmt.Sprintf("%d BTC", b)
}

//不能透支钱包里的钱，需要引出一个错误

var InsufficientFundsError = errors.New("cannot withdraw, insufficient funds")

func (w *Wallet) Withdraw(amount Bitcoin) error {

	if amount > w.balance {
		return InsufficientFundsError
	}
	w.balance -= amount
	return nil
}
