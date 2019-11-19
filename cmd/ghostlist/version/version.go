package version

import "fmt"

//this is added at build time, dont touch it
var version = "undefined"

func Show() {
	fmt.Println(version)
}
