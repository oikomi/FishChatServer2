package lib

import (
	"testing"
)

func TestDiff(t *testing.T) {
	a := []string{"1", "2", "3", "4"}
	b := []string{"3", "4", "5", "6"}

	c := diff(a, b)
	d := diff(b, a)

	if len(c) != 2 || c[0] != "1" || c[1] != "2" {
		t.Errorf("diff(a-b) error, get=%v", c)
	}

	if len(d) != 2 || d[0] != "5" || d[1] != "6" {
		t.Errorf("diff(b-a) error, get=%v", d)
	}
}
