package main

import "testing"

func TestSearch(t *testing.T) {
	//map 关键字开头，需要两种类型。第一个是键的类型，写在 [] 中。第二个是值的类型，跟在 [] 之后。
	dictionary := Dictionary{"test": "我不是人，我只是个怪物"}
	//值的类型可以是任意类型

	t.Run("known word", func(t *testing.T) {
		//参数一个是自测试的名称，另一个是包含测试逻辑的匿名函数
		got, _ := dictionary.Search("test")
		want := "我不是人，我只是个怪物"

		assertStrings(t, got, want)
	})

	t.Run("unknown word", func(t *testing.T) {
		_, err := dictionary.Search("unknown")
		want := "你脑袋怎么尖尖的"

		if err == nil {
			t.Fatal("expected to get an error")
		}

		assertStrings(t, err.Error(), want)
	})

}

// 重构
func assertStrings(t *testing.T, got, want string) {
	//t *testing.T同理，t是参数，空格后是类型
	t.Helper()

	if got != want {
		t.Errorf("got '%s' want '%s'", got, want)
	}
}
