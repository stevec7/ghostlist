package ghostlist

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
    "github.com/golang-collections/collections/set"
)

const MAX_SIZE int = 100000

/*
Collect a hostlist string from a string slice of hosts.

We start grouping from the rightmost numerical part.
Duplicates are removed.
*/
func CollectHostList(hostlist []string) (string, error) {
	var leftRight []leftRightRec

	hostlist = removeDups(hostlist)

	for _, host := range hostlist {
		s := strings.TrimSpace(host)
		if host == "" {
			continue
		}

		re := regexp.MustCompile(`[][,]`)
		if re.Match([]byte(host)) {
			return "", errors.New("Forbidden characters in host list, [][,]")
		}

        rec := leftRightRec{l: s, r: ""}
		leftRight = append(leftRight, rec)
	}
	looping := true
	for {
		leftRight, looping = CollectHostListOne(leftRight)
		if !looping {
			break
		}
	}
	var results []string
	for _, i := range leftRight {
		s := fmt.Sprintf("%s%s", i.l, i.r)
		results = append(results, s)
	}
	return strings.Join(results, ","), nil
}



/*
Collect a hostlist string from a list of hosts (left+right).

The input is a list of tuples (left, right). The left part
is analyzed, while the right part is just passed along
(it can contain already collected range expressions).
*/
func CollectHostListOne(leftRight []leftRightRec) ([]leftRightRec, bool) {
	var sL []sortListT
	//var remaining []string //ill handle the set stuff after by just removing the dupes
	remaining := set.New()

	for _, lr := range leftRight {
		left := lr.l
		right := lr.r
		host := fmt.Sprintf("%s%s", string(left), string(right))
        remaining.Insert(host)

		re := regexp.MustCompile(`^(.*?)([0-9]+)?([^0-9]*)$`)
		groups := re.FindStringSubmatch(string(left))
		pre := groups[1]
		numStr := groups[2]
		suf := groups[3]

		suf = fmt.Sprintf("%s%s", suf, string(right))

		if numStr == "" {
            tmp := sortListT{
                // i am using "None" here because the stupid python version uses None, and
                //  go doesn't allow nil string assignment
                preSuf: prefixSuffix{prefix: pre, suffix: "None"},
                numInt: 0,
                numWidth: 0,
                host: host,
            }
			sL = append(sL, tmp)
		} else {
			numI, _ := strconv.Atoi(numStr)
			numW := len(numStr)
            sL = append(sL, sortListT{
                preSuf: prefixSuffix{prefix: pre, suffix: suf},
                numInt: numI,
                numWidth: numW,
                host: host,
            })

		}
	}

    // TODO: need to figure out a nice way to sort in place, using the prefix, and then the suffix

    needsAnotherLoop := false

    var results []leftRightRec
    for _, g := range groupBy(sL) {
        if g.preSuf.suffix == "None" {
            results = append(results, leftRightRec{l: "", r: g.preSuf.prefix})
            remaining.Remove(g.preSuf.prefix)
        } else {
            var rL []rangeList
            for _, m := range g.members {
                if ok := remaining.Has(m.host); !ok {
                    continue
                }

				numInt := m.numInt
                low := m.numInt
                for {
                    newhost := fmt.Sprintf("%s%0*d%s", m.preSuf.prefix, m.numWidth,
						m.numInt, m.preSuf.suffix)
                    if ok := remaining.Has(newhost); ok {
                        remaining.Remove(newhost)
                        numInt += 1
                    } else {
                        break
                    }
                }
                high := numInt - 1
                rL = append(rL, rangeList{low, high, m.numWidth})

            }
            needsAnotherLoop = true
            if len(rL) == 1 && rL[0].low == rL[0].high {
                results = append(results, leftRightRec{l: g.preSuf.prefix,
					r: fmt.Sprintf("%0*d%s", rL[0].numWidth, rL[0].low, g.preSuf.suffix)})
            } else {
                var tmp []string
                for _, i := range rL {
                    tmp = append(tmp, formatRange(i.low, i.high, i.numWidth))
                }
                results = append(results, leftRightRec{l: g.preSuf.prefix,
					r: fmt.Sprintf("[%s]%s", strings.Join(tmp, ","), g.preSuf.suffix) })
            }
        }
    }

	needsAnotherLoop = false

    return results, needsAnotherLoop
}

/*
Expand a hostlist expression string to a slice

Example: expand_hostlist("n[9-11],d[01-02]") ==>
         ['n9', 'n10', 'n11', 'd01', 'd02']

Duplicates will be removed, and the results will be sorted
*/
func ExpandHostList(hostlist string) ([]string, error) {
	var results []string
	bracketLevel := 0
	part := ""

	for _, c := range fmt.Sprintf("%s,", hostlist) {
		if string(c) == "," && bracketLevel == 0 {
			if len(part) > 0 {
				r, err := ExpandPart(part)
				if err != nil {
					return []string{}, err
				}
				results = append(results, r...)
			}
			part = ""

		} else {
			part += string(c)
		}

		if string(c) == "[" {
			bracketLevel += 1
		} else if string(c) == "]" {
			bracketLevel -= 1
		}

		if bracketLevel > 1 {
			return []string{}, errors.New("Error, nested brackets.")
		} else if bracketLevel < 0 {
			return []string{}, errors.New("Error, unbalanced brackets.")
		}
	}

	if bracketLevel > 0 {
		return []string{}, errors.New("Error, unbalanced brackets")
	}

	// remove dups
	results = removeDups(results)

	// sort
	err := sortHostlist(&results)

	if err != nil {
		return results, err
	}

	return results, nil
}

// Expand a part (e.g. "x[1-2]y[1-3][1-3]") (no outer level commas).
func ExpandPart(s string) ([]string, error) {
	if s == "" {
		return []string{""}, nil
	}

	re := regexp.MustCompile(`([^,\[]*)(\[[^\]]*\])?(.*)`)
	groups := re.FindStringSubmatch(s)

	prefix := groups[1]
	rangeList := groups[2]
	rest := groups[3]

	restExpanded, err := ExpandPart(rest)
	if err != nil {
		return []string{}, err
	}

	var usExpanded []string
	if rangeList == "" {
		usExpanded = []string{prefix}
	} else {
		usExpanded, err = ExpandRangeList(prefix, rangeList[1:len(rangeList)-1])
		if err != nil {
			return []string{}, err
		}
	}

	if (len(usExpanded) * len(restExpanded)) > MAX_SIZE {
		return []string{}, errors.New("results too large")
	}

	var results []string
	for _, u := range usExpanded {
		for _, r := range restExpanded {
			results = append(results, fmt.Sprintf("%s%s", u, r))
		}
	}
	return results, nil
}

// Expand a rangelist (e.g. "1-10,14"), putting a prefix before.
func ExpandRangeList(prefix, rnglist string) ([]string, error) {
	var results []string
	for _, r := range strings.Split(rnglist, ",") {
		result, err := ExpandRange(prefix, r)
		if err != nil {
			return []string{}, err
		}
		results = append(results, result...)
	}
	return results, nil
}

// Expand a range (e.g. 1-10 or 14), putting a prefix before.
func ExpandRange(prefix, rng string) ([]string, error) {
	matcher := regexp.MustCompile(`^[0-9]+$`)

	// single number
	if matcher.Match([]byte(rng)) {
		return []string{fmt.Sprintf("%s%s", prefix, rng)}, nil
	}

	matcher = regexp.MustCompile(`^([0-9]+)-([0-9]+)$`)
	if !matcher.Match([]byte(string(rng))) {
		return []string{}, errors.New("Malformed host string, bad range.")
	}
	groups := matcher.FindStringSubmatch(rng)

	stringLow := groups[1]
	stringHigh := groups[2]
	low, _ := strconv.Atoi(stringLow)
	high, _ := strconv.Atoi(stringHigh)
	width := len(stringLow)

	if high < low {
		return []string{}, errors.New("start > stop")
	} else if (high - low) > MAX_SIZE {
		return []string{}, errors.New("range too large")
	}

	var results []string
	for i := low; i < high+1; i++ {
		s := fmt.Sprintf("%s%0*d", prefix, width, i)
		results = append(results, s)
	}

	return results, nil
}
