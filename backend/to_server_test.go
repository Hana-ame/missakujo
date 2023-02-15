package backend

import (
	"fmt"
	"testing"
	"time"
)

func TestTIme(t *testing.T) {
	fmt.Println(time.Now().Local().Zone())
}

func TestDelete(t *testing.T) {
	since, err := time.Parse(timeForm, "2023-02-11 11:19:34")
	if err != nil {
		t.Fatal(err)
	}
	until, err := time.Parse(timeForm, "2023-02-28 12:34:56")
	if err != nil {
		t.Fatal(err)
	}

	_, offset := time.Now().Local().Zone()

	sinceInt := since.UnixMilli()
	untilInt := until.UnixMilli()

	sinceInt -= int64(offset) * 1000
	untilInt -= int64(offset) * 1000

	HandleDelete(
		"p1.a9z.dev",
		"9a3qtdtypj",
		"kan1imxGYSfggpHtRZCLMhu35ykPjyi7",
		sinceInt, untilInt,
		999,
		"true", "true",
		"123",
	)
}
