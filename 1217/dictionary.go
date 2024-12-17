package main

type Dictionary map[string]string

//初始化空map方法如下
//var m map[string]string
//或者使用make关键字，如下：
/*
  dictionary = map[string]string{}
或
  dictionary = make(map[string]string)
*/

//重构

// 重构，将错误声明为常量，这需要我们创建自己的 DictionaryErr 类型来实现 error 接口
const (
	ErrNotFound         = DictionaryErr("你脑袋怎么尖尖的")
	ErrWordExists       = DictionaryErr("住手，你们不要再打了啦")
	ErrWordDoesNotExist = DictionaryErr("流脓")
)

type DictionaryErr string

func (e DictionaryErr) Error() string {
	return string(e)
}

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

func (d Dictionary) Add(word, definition string) error {
	_, err := d.Search(word)

	switch err {
	case ErrNotFound:
		d[word] = definition
	case nil:
		return ErrWordExists
	default:
		return err
	}
	return nil
}

func (d Dictionary) Update(word, definition string) error {
	_, err := d.Search(word)

	switch err {
	case ErrNotFound:

		return ErrWordDoesNotExist
	case nil:
		d[word] = definition
	default:
		return err
	}
	return nil
}
