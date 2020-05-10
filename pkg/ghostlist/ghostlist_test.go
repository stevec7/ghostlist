package ghostlist

import (
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

func TestCHL(t *testing.T) {
	// simple
	hostnames := []string{"test1", "test2", "test3"}
	a, _ := CollectHostList(hostnames)
	b := "test[1-3]"
	r := a == b
	if !r {
		t.Errorf("Hostlist was wrong, expected: (%v), got: (%v)", b, a)
	}

	// test intermediate range
	hostnames = []string{"test1-2-3", "test1-2-4", "test1-2-5", "test1-2-7"}
	a, _ = CollectHostList(hostnames)
	b = "test1-2-[3-5,7]"
	r = a == b
	if !r {
		t.Errorf("Hostlist was wrong, expected: (%v), got: (%v)", b, a)
	}

	// test intermediate range with different middles
	hostnames = []string{"test1-2-3", "test1-2-4", "test1-2-5", "test1-2-7", "test1-3-12", "test1-3-14", "test1-3-15", "test1-3-16"}
	a, _ = CollectHostList(hostnames)
	b = "test1-2-[3-5,7],test1-3-[12,14-16]"
	r = a == b
	if !r {
		t.Errorf("Hostlist was wrong, expected: (%v), got: (%v)", b, a)
	}

	// test jumble of random junk
	hostnames = []string{"test0001", "test0002", "test0003", "test04", "test05", "test07-01", "test07-02", "test-07-02-01", "test07-02-02", "test07-02-03", "test07-001", "test07-002", "test07-001-001", "test07-002-002", "test1000", "test2001", "test3002"}
	a, _ = CollectHostList(hostnames)
	// love you long line
	b = "test-07-02-01,test07-001-001,test07-002-002,test07-02-[02-03],test07-[001-002,01-02],test[0001-0003,04-05,1000,2001,3002]"
	r = a == b
	if !r {
		t.Errorf("Hostlist was wrong, expected: (%v), got: (%v)", b, a)
	}
}

func TestERL(t *testing.T) {
	a, _ := expandRangeList("host", "1-2,10")
	b := []string{"host1", "host2", "host10"}

	if !equal(a, b) {
		t.Errorf("Hostlist was wrong, expected: (%v), got: (%v)", b, a)
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
	a, _ := ExpandHostList("n[09-11],d[01-02]")
	b := []string{"d01", "d02", "n09", "n10", "n11"}

	if !equal(a, b) {
		t.Errorf("Hostlist was wrong, expected: (%v), got: (%v)", b, a)
	}
}

func TestEHLComplex(t *testing.T) {
	a, _ := ExpandHostList("x[1-2]y[1-3][001-004]")
	b := []string{"x1y1001", "x1y1002", "x1y1003", "x1y1004", "x1y2001", "x1y2002", "x1y2003", "x1y2004",
		"x1y3001", "x1y3002", "x1y3003", "x1y3004", "x2y1001", "x2y1002", "x2y1003", "x2y1004", "x2y2001",
		"x2y2002", "x2y2003", "x2y2004", "x2y3001", "x2y3002", "x2y3003", "x2y3004"}

	if !equal(a, b) {
		t.Errorf("Hostlist was wrong, expected: (%v), got: (%v)", b, a)
	}
}
