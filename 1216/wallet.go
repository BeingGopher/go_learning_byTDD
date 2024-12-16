package main

type Wallet struct {
	//在 Go 中，如果一个符号（例如变量、类型、函数等）是以小写符号开头，那么它在 定义它的包之外 就是私有的。
	balance int
}

func (w Wallet) Deposit(amount int) {
	//使用「receiver」变量访问结构体内部的 balance 字段。
	w.balance += amount
}

func (w Wallet) Balance() int {
	return w.balance
}
