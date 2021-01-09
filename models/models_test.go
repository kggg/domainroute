package models

import (
	"fmt"
	"testing"
)

func TestRouteTables(t *testing.T) {
	tables, err := getRouteTables()
	if err != nil {
		t.Errorf("%v", err)
	}
	for _, v := range tables {
		fmt.Println(v)
	}
}
