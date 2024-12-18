package integers

import (
	"fmt"
	"testing"
)

func TestAdder(t *testing.T) {
	AssertCorrectMessage := func(t *testing.T, expected, sum int) {
		t.Helper()
		if sum != expected {
			t.Errorf("expected '%d' but got '%d'", expected, sum)
		}

	}

	sum := Add(2, 2)
	expected := 4

	AssertCorrectMessage(t, expected, sum)
}

func ExampleAdd() {
	sum := Add(1, 5)
	fmt.Println(sum)
	// Output: 6
}
