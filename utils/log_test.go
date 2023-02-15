package utils

import (
	"os"
	"testing"
)

func TestLog(t *testing.T) {
	fw, err := NewFileWriter("123.txt")
	if err != nil {
		t.Error(err)
	}
	fw.WriteString("123\n")
	fw.Log("some message loged")
	fw.Close()
}

func TestAppend(t *testing.T) {
	AppendStringToFile("123.txt", "some message appened")
}

// func TestOverwrite(t *testing.T) {
// 	err := OverwriteStringToFile("456.txt", "some message\n")
// 	t.Error(err)
// }

func TestDelete(t *testing.T) {
	if err := os.Remove("123.txt"); err != nil {
		t.Error(err)
	}
}
