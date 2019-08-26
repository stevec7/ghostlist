package ghostlist

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

const MAX_SIZE int = 100000

/*
Collect a hostlist string from a string slice of hosts.

We start grouping from the rightmost numerical part.
Duplicates are removed.
 */
func CollectHostList(hostlist []string) (string, error) {
	var leftRight []string

	for _, host := range hostlist {
		s := strings.TrimSpace(host)
		if host == "" {
			continue
		}

		re := regexp.MustCompile(`[][,]`)
		if re.Match([]byte(host)) {
			return "", errors.New("Forbidden characters in host list, [][,]")
		}

		lr := []string{s, ""}
		leftRight = append(leftRight, lr...)
	}
	looping := true
	for {
		leftRight, looping = CollectHostListOne(leftRight)
		if !looping {
			break
		}
	}
	var results []string
	for l, r := range leftRight {
		s := fmt.Sprintf("%s%s", l, r)
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
func CollectHostListOne(leftRight []string) ([]string, bool){
	//var sortList []string
	var sortList []interface{}
	var remaining  []string	//ill handle the set stuff after by just removing the dupes

	for _, lr := range leftRight {
		left := lr[0]
		right := lr[1]
		host := fmt.Sprintf("%s%s", left, right)
		remaining = append(remaining, host)

		re := regexp.MustCompile(`^(.*?)([0-9]+)?([^0-9]*)$`)
		groups := re.FindStringSubmatch(string(left))
		prefix := groups[1]
		numStr := groups[2]
		suffix := groups[3]

		suffix = fmt.Sprintf("%s%s", suffix, right)

		if numStr == "" {
			fmt.Println("What the heck...")
			s1 := []interface{}{prefix, nil}
			s2 := []interface{}{nil, nil, host}
			s3 := []interface{}{s1, s2}
			sortList = append(sortList, s3)
		} else {
			numInt, _ := strconv.Atoi(numStr)
			numWidth := len(numStr)
			s1 := []string{prefix, suffix}
			s2 := []interface{}{numInt, numWidth, host}
			s3 := []interface{}{s1, s2}
			sortList = append(sortList, s3)
		}
	}

	var results []string
	needsAnotherLoop := false


	return []string{"a", "b", "c"}, true
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
func ExpandPart(s string) ([]string, error){
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
	for i := low; i < high + 1; i++ {
		s := fmt.Sprintf("%s%0*d", prefix, width, i)
		results = append(results, s)
	}

	return results, nil
}

func removeDups(hostlist []string) ([]string) {
	// this is tedious, and/or i am dumb
	m := make(map[string]int, len(hostlist))
	for i, v := range hostlist {
		m[v] = i
	}
	var keys []string
	for k, _ := range m {
		keys = append(keys, k)
	}
	return keys
}

func sortHostlist(hostlist *[]string) error {
	sort.Strings(*hostlist)

	return nil
}

/*
func main() {
	a, err := ExpandRangeList("host", "1-7,10")
	if err != nil {
		fmt.Println("Error")
	}
	fmt.Println("a: ", a)
	b, err := ExpandPart("x[1-2]y[1-3][1-3]")
	if err != nil {
		fmt.Println("Error")
	}
	fmt.Println("b: ", b)
	fmt.Println("len(b): ", len(b))

	//c, err := ExpandHostList("n[9-11],d[01-02]")
	c, err := ExpandHostList("host0[100-102],x[500-503],host0100,host0115,host0116,d[01-05,07],host0116")
	if err != nil {
		fmt.Println("Error: ", err)
	}
	fmt.Println("c: ", c)

	d := removeDups(c)
	//sort.Strings(d)
	err = sortHostlist(&d)

	if err != nil {
		fmt.Println("He didnt fly so good")
	}
	fmt.Println("removed dupes: ", d)
	fmt.Println("sorted: ", d)
	//c, _ := ExpandRangeList("y[1-3][1-3]")
	//fmt.Println("c: ", c)
}
*/
