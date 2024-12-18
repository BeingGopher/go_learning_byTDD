package iteration

//repeat括号后面的string是代表返回的结果的类型

func Repeat(character string, n int) string {
	var repeated string
	for i := 0; i < n; i++ {
		repeated += character
	}
	return repeated
}
