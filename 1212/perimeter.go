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
