package main

import (
	"testing"
)

func TestSum(t *testing.T) {
	want := 5
	res := sum(2,3)
	if want != res {
		t.Errorf("Sum func error")
	}
}
