package main

import (
	"fmt"
	"os/exec"
)

// Iterates over Godep file dependencies and sets
// the specified version on each of them.
func install() {
	if !fileExists(depsFile) {
		panic(fmt.Sprintf(">> %s file does not exist.\n", depsFile))
	}

	_, err := exec.LookPath("go")
	if err != nil {
		panic(">> Go is currently not installed or in your PATH\n")
	}

	for pkg, version := range parseDepFile() {
		fmt.Printf(">> Getting %s\n", pkg)
		execCmd(fmt.Sprintf(`go get -u -d "%s"`, pkg))
		setPackageToVersion(pkg, version)
	}

	fmt.Println(">> All Done")
}
