package main

import "fmt"

// 定义常量，避免每次使用函数时创建字符串实例
const helloPrefix = "Hello, "
const spanish = "Spanish"
const spanishHelloPrefix = "Hola, "
const french = "French"
const frenchHelloPrefix = "Bonjour, "
const chinese = "Chinese"
const chineseHelloPrefix = "你好, "

func Hello(name string, language string) string {
	if name == "" {
		name = "World"
	}

	return greetingPrefix(language) + name
}

// 私有函数，小写字母开头，只在内部调用？
func greetingPrefix(language string) (prefix string) {
	switch language {
	case french:
		prefix = frenchHelloPrefix
	case spanish:
		prefix = spanishHelloPrefix
	case chinese:
		prefix = chineseHelloPrefix
	default:
		prefix = helloPrefix
	}
	//已经设置了命名返回值，所以直接return即可，不需要return prefix
	return
}

func main() {
	fmt.Println(Hello("蔡徐坤!", ""))
}
