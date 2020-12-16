package ghostlist

import (
	"sort"
	"strconv"
	"unicode"
)

func nextHostSpan(s string) (remainder string, span string, n bool) {
	if s == "" {
		return
	}
	for i, r := range s {
		if i == 0 {
			n = unicode.IsNumber(r)
		} else if n != unicode.IsNumber(r) {
			return s[i:], s[:i], n
		}
	}
	return "", s, n
}

func sortHost(s1, s2 string) bool {
	for {
		var sp1, sp2 string
		var i, j bool
		s1, sp1, i = nextHostSpan(s1)
		s2, sp2, j = nextHostSpan(s2)
		if sp1 == "" && sp2 == "" {
			return false
		}
		if i && j {
			if n1, err := strconv.ParseUint(sp1, 10, 64); err == nil {
				if n2, err := strconv.ParseUint(sp2, 10, 64); err == nil {
					if n1 != n2 {
						return n1 < n2
					}
					continue
				}
			}
		}
		if sp1 != sp2 {
			return sp1 < sp2
		}
	}
}

func sortHostlist(hostlist *[]string) {
	sort.Slice(*hostlist, func(i, j int) bool {
		return sortHost((*hostlist)[i], (*hostlist)[j])
	})
}
