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
		_, got := dictionary.Search("unknown")

		assertError(t, got, ErrNotFound)
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

// 创建新的辅助函数来简化测试，并调用ErrNotFound变量
func assertError(t *testing.T, got, want error) {
	t.Helper()

	if got != want {
		t.Errorf("got error '%s' want '%s'", got, want)
	}
}

// 编写添加新单词功能测试
func TestAdd(t *testing.T) {
	dictionary := Dictionary{}
	dictionary.Add("test", "再多看一眼就会爆炸")

	want := "再多看一眼就会爆炸"
	got, err := dictionary.Search("test")
	if err != nil {
		t.Fatal("should find added word:", err)
	}

	if want != got {
		t.Errorf("got '%s' want '%s'", got, want)
	}

}
