package main

import (
	"flag"
	"fmt"
	"github.com/stevec7/ghostlist/pkg/ghostlist"
	"log"
	"strings"
)

func main() {
	expandP := flag.String("e", "host00[1-3]", "Expand the hostlist")
	flag.Parse()

	result, err := ghostlist.ExpandHostList(*expandP)

	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	fmt.Printf("%v\n", strings.Join(result, ","))
}
