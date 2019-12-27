package main

import (
	"flag"
	"fmt"
	"github.com/stevec7/ghostlist/cmd/ghostlist/version"
	"github.com/stevec7/ghostlist/pkg/ghostlist"
	"log"
	"os"
	"strings"
)

func main() {
	collapseP := flag.String("c", "", "Collapse the hostlist, e.g. host1,host2,host3")
	expandP := flag.String("e", "", "Expand the hostlist, e.g. host[1-3]")
	showVersion := flag.Bool("V", false, "show the version")
	flag.Parse()

	if *showVersion {
		version.Show()
		os.Exit(0)
	}

	if *expandP != "" {
		result, err := ghostlist.ExpandHostList(*expandP)
		if err != nil {
			log.Fatalf("Error: %s\n", err)
		}

		fmt.Printf("%v\n", strings.Join(result, ","))
	} else if *collapseP != "" {
		result, err := ghostlist.CollectHostList(*collapseP)
		if err != nil {
			log.Fatalf("Error: %s\n", err)
		}
		fmt.Printf("%s\n", result)
	} else {
		log.Fatalf("You must enter a value for the [-c|-e] args")
	}


}
