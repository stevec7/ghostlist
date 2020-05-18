package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/stevec7/ghostlist/cmd/ghostlist/version"
	"github.com/stevec7/ghostlist/pkg/ghostlist"
)

func main() {
	collapseP := flag.Bool("c", false, "Collapse the hostlist into a pdsh style string, e.g. host1,host2,host3")
	expandP := flag.Bool("e", false, "Expand the hostlist into a comma separate list of all hosts, e.g. host[1-3]")
	intersectP := flag.Bool("i", false, "Intersection between 2 pdsh style hostlists, eg: hosts[1-3] hosts[1-2]")
	diffP := flag.Bool("d", false, "Difference between 2 pdsh style hostlists, eg: hosts[1-3] hosts[1-2]")
	showVersion := flag.Bool("V", false, "show the version")
	flag.Parse()

	if *showVersion {
		version.Show()
		os.Exit(0)
	}

	var result string
	var err error
	args := flag.Args()

	switch {
	case *expandP:
		if flag.NArg() != 1 {
			fmt.Printf("Error must supply 1 argument, a pdsh style hostlist")
			os.Exit(1)
		}
		r, err := ghostlist.ExpandHostList(args[0])
		if err != nil {
			fmt.Printf("Error %s\n", err)
			os.Exit(1)
		}
		sort.Strings(r)
		result = strings.Join(r, ",")

	case *collapseP:
		if flag.NArg() != 1 {
			fmt.Printf("Error must supply 1 argument, a pdsh style hostlist")
			os.Exit(1)
		}

		hosts := strings.Split(args[0], ",")
		result, err = ghostlist.CollectHostList(hosts)
		if err != nil {
			fmt.Printf("Error %s\n", err)
			os.Exit(1)
		}
	case *intersectP:
		if flag.NArg() != 2 {
			fmt.Printf("Error must supply 2 arguments, both pdsh style hostlists")
			os.Exit(1)
		}
		result, err = ghostlist.Intersection(args[0], args[1])
		if err != nil {
			fmt.Printf("Error %s\n", err)
			os.Exit(1)
		}

	case *diffP:
		if flag.NArg() != 2 {
			fmt.Printf("Error must supply 2 arguments, both pdsh style hostlists")
			os.Exit(1)
		}
		result, err = ghostlist.Difference(args[0], args[1])
		if err != nil {
			fmt.Printf("Error %s\n", err)
			os.Exit(1)
		}
	default:
		fmt.Printf("Please choose an operation, see --help for usage")
	}
	fmt.Println(result)
}
