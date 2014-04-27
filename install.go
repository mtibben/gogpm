package main

import "fmt"

// Iterates over Godep file dependencies and sets
// the specified version on each of them.
func install() {
	for pkg, version := range parseDepFile() {
		fmt.Printf(">> Getting %s\n", pkg)
		execCmd(fmt.Sprintf(`go get -u -d "%s"`, pkg))
		setPackageToVersion(pkg, version)
	}

	fmt.Println(">> All Done")
}
