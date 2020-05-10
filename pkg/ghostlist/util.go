package ghostlist

import (
	"fmt"
	"sort"
)

func formatRange(low, high, width int) string {
	if low == high {
		return fmt.Sprintf("%0*d", width, low)
	}
	return fmt.Sprintf("%0*d-%0*d", width, low, width, high)
}

func groupBy(s []sortListT) []sortListG {
	keyloc := make(map[prefixSuffix]*sortListG)
	var r []sortListG

	for _, item := range s {
		if _, ok := keyloc[item.preSuf]; !ok {
			g := sortListG{preSuf: item.preSuf, members: []sortListT{}}
			r = append(r, g)
			g.members = append(g.members, item)
			keyloc[item.preSuf] = &g
		} else {
			l := keyloc[item.preSuf]
			l.members = append(l.members, item)
		}
	}

	for i, v := range r {
		r[i] = *keyloc[v.preSuf]
	}
	return r
}

func removeDups(hostlist []string) []string {
	// this is tedious, and/or i am dumb
	m := make(map[string]int, len(hostlist))
	for i, v := range hostlist {
		m[v] = i
	}
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
