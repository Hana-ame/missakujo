package utils

import (
	"fmt"
	"testing"
)

func tester[T any](v any) (vv T, err error) {
	defer func() {
		e := recover()
		if e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()

	vv = v.(T)
	return
}

func TestXxx(t *testing.T) {

	a, err := tester[string](1)
	fmt.Println(a, err)

	b, err := tester[int]("")
	fmt.Println(b, err)

	c, err := tester[int](213)
	fmt.Println(c, err)
}
