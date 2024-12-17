package main

import "errors"

type Dictionary map[string]string

//重构

var ErrNotFound = errors.New("你脑袋怎么尖尖的")

func (d Dictionary) Search(word string) (string, error) {
	//接受者、实例类型； 空格后是方法，方法内是传参；最后是返回值类型
	definition, ok := d[word]
	//map特性，可以返回两个值。第二个值是一个布尔值，表示是否成功找到 key。
	if !ok {
		return " ", ErrNotFound
		//errors.New函数创建一个新的错误实例，错误信息为中文字符串。前面是返回一个空字符串，后面是返回错误信息
	}

	return definition, nil
}

func Search(dictionary map[string]string, word string) string {
	//函数定义：参数+类型；括号外是返回值类型
	return dictionary[word]
}
