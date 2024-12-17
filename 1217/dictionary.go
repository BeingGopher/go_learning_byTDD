package main

type Dictionary map[string]string

func (d Dictionary) Search(word string) string {
	//接受者、实例类型； 空格后是方法，方法内是传参；最后是返回值类型
	return d[word]
}

func Search(dictionary map[string]string, word string) string {
	//函数定义：参数+类型；括号外是返回值类型
	return dictionary[word]
}
