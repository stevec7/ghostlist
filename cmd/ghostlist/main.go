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
	expandP := flag.String("e", "host00[1-3]", "Expand the hostlist")
	showVersion := flag.Bool("V", false, "show the version")
	flag.Parse()

	if *showVersion {
		version.Show()
		os.Exit(0)
	}

	result, err := ghostlist.ExpandHostList(*expandP)

	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	fmt.Printf("%v\n", strings.Join(result, ","))
}
