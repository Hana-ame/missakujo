package utils

import (
	"fmt"
	"testing"
)

func Test1(t *testing.T) {
	r, e := ResolveUser("meromero@p1.a9z.dev")
	if e != nil {
		t.Error(e)
	}
	fmt.Println(r)
}
