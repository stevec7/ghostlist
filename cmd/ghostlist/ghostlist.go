package main

import (
	"flag"
	"fmt"
	"github.com/stevec7/ghostlist"
	"log"
)

func main() {
	expand := flag.String("expand", "host00[1-3]", "Expand the hostlist" )
	flag.Parse()

	result, err := ghostlist.ExpandHostList(*expand)

	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	fmt.Printf("%v\n", result)
}
