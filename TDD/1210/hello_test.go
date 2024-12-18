package main

import "testing"

func TestHello2(t *testing.T) {
	assertCorrectMessage := func(t *testing.T, got, want string) {
		t.Helper()
		if got != want {
			t.Errorf("got '%q' want '%q'", got, want)
		}
	}
	t.Run("saying 鸡你太美 to people", func(t *testing.T) {
		got := Hello("蔡徐坤", "")
		want := "Hello, 蔡徐坤"

		assertCorrectMessage(t, got, want)
	})

	t.Run("empty string defaults to 'world'", func(t *testing.T) {
		got := Hello("", "")
		want := "Hello, World"

		assertCorrectMessage(t, got, want)
	})

	t.Run("in Spanish", func(t *testing.T) {
		got := Hello("Elodie", "Spanish")
		want := "Hola, Elodie"
		assertCorrectMessage(t, got, want)
	})

	t.Run("in French", func(t *testing.T) {
		got := Hello("蔡徐坤", "French")
		want := "Bonjour, 蔡徐坤"

		assertCorrectMessage(t, got, want)
	})

	t.Run("in Chinese", func(t *testing.T) {
		got := Hello("蔡徐坤", chinese)
		want := "你好, 蔡徐坤"
		assertCorrectMessage(t, got, want)
	})
}
