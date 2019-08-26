package ghostlist_test

import (
	"github.com/stevec7/ghostlist/pkg/ghostlist"
	"testing"
)

func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func TestERL(t *testing.T) {
	a, _:= ghostlist.ExpandRangeList("host", "1-2,10")
	b := []string{"host1", "host2", "host10"}

	if !equal(a, b) {
		t.Errorf("Hostlist was wrong, expected: (%v), got: (%v_", b, a)
	}
}
/*
func TestEHLNoPadding(t *testing.T) {
	a, _ := ghostlist.ExpandHostList("n[9-11],d[01-02]")
	b := []string{"d01" ,"d02", "n9", "n10", "n11"}

	if !equal(a, b) {
		t.Errorf("Hostlist was wrong, expected: (%v), got: (%v_", b, a)
	}
}
*/
func TestEHLPadding(t *testing.T) {
	a, _ := ghostlist.ExpandHostList("n[09-11],d[01-02]")
	b := []string{"d01" ,"d02", "n09", "n10", "n11"}

	if !equal(a, b) {
		t.Errorf("Hostlist was wrong, expected: (%v), got: (%v_", b, a)
	}
}

func TestEHLComplex(t *testing.T) {
	a, _ := ghostlist.ExpandHostList("x[1-2]y[1-3][001-004]")
	b := []string{"x1y1001", "x1y1002", "x1y1003", "x1y1004", "x1y2001", "x1y2002", "x1y2003", "x1y2004",
		"x1y3001", "x1y3002", "x1y3003", "x1y3004", "x2y1001", "x2y1002", "x2y1003", "x2y1004", "x2y2001",
		"x2y2002", "x2y2003", "x2y2004", "x2y3001", "x2y3002", "x2y3003", "x2y3004"}

	if !equal(a, b) {
		t.Errorf("Hostlist was wrong, expected: (%v), got: (%v_", b, a)
	}
}
