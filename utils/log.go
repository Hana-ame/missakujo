package utils

import (
	"fmt"
	"os"
	"time"
)

type FileWriter struct {
	*os.File
}

func NewFileWriter(fn string) (*FileWriter, error) {
	file, err := os.OpenFile(fn, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	return &FileWriter{File: file}, nil
}

// func (f *FileWriter) Write(b []byte) (int, error) {
// 	return f.Write(b)
// }

// func (f *FileWriter) WriteString(s string) (int, error) {
// 	return f.WriteString(s)
// }

// func (f *FileWriter) Close() error {
// 	return f.Close()
// }

func (f *FileWriter) Log(s string) {
	f.WriteString("[")
	f.WriteString(time.Now().Format("15:04:05"))
	f.WriteString("]")
	f.WriteString("\t")
	f.WriteString(s)
	f.WriteString("\n")
}

func (f *FileWriter) Println(a ...any) error {
	s := fmt.Sprintln(a...)
	if _, err := f.WriteString(s); err != nil {
		return err
	}
	return nil
}

func AppendStringToFile(fn string, s string) error {
	file, err := NewFileWriter(fn)
	if err != nil {
		return err
	}

	if _, err := file.WriteString(s); err != nil {
		return err
	}

	file.WriteString("\n\n\n")

	return file.Close()
}

// func OverwriteStringToFile(fn string, s string) error {
// 	file, err := NewFileWriter(fn)
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()

// 	if err := file.Truncate(0); err != nil { //  truncate 456.txt: Access is denied.
// 		return err
// 	}

// 	if _, err := file.Seek(0, 0); err != nil {
// 		return err
// 	}

// 	if _, err := file.WriteString(s); err != nil {
// 		return err
// 	}

// 	return nil
// }
