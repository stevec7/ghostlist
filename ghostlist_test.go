package main

import (
	"github.com/stevec7/ghostlist"
	"testing"
)

func TestERL(t *testing.T) {
	a, _:= ExpandRangeList("host", "1-2,10")
	expected := []string{"host1", "host2", "host10"}
	if a != expected {
		t.Errorf("Hostlist was wrong, expected: (%v), got: (%v_", expected, a)
	}

}
