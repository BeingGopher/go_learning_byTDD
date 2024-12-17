package main

import "math"

func Perimeter(rectangle Rectangle) float64 {
	return 2 * (rectangle.n + rectangle.m)
}

func Area(rectangle Rectangle) float64 {
	return rectangle.n * rectangle.m
}

//定义结构体，比如长方形

type Rectangle struct {
	n float64
	m float64
}

func (r Rectangle) Area() float64 {
	return r.n * r.m
}

type Circle struct {
	r float64
}

//把类型的第一个字母作为接收者变量

func (c Circle) Area() float64 {
	return math.Pi * c.r * c.r
}

//通过声明一个接口，辅助函数能从具体类型解耦而只关心方法本身需要做的工作。

type Shape interface {
	Area() float64
}

type Triangle struct {
	a float64
	b float64
}

func (c Triangle) Area() float64 {
	return c.a * c.b / 2
}
