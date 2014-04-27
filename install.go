package main

import (
	"fmt"
	"os"
)

// Iterates over Godep file dependencies and sets
// the specified version on each of them.
func install() {
	for pkg, version := range parseDepFile() {
		fmt.Printf(">> Getting package %s\n", pkg)
		execCmd(fmt.Sprintf(`go get -u -d "%s"`, pkg))

		fmt.Printf(">> Setting %s to version %s\n", pkg, version)
		setPackageToVersion(pkg, version)
	}

	fmt.Println(">> All Done")
}

func setPackageToVersion(pkg, version string) {
	mustChdir(installPath(pkg))

	if dirExists(".hg") {
		execCmd(fmt.Sprintf(`hg update -q "%s"`, version))
	} else if dirExists(".git") {
		execCmd(fmt.Sprintf(`git checkout -q "%s"`, version))
	} else if dirExists(".bzr") {
		execCmd(fmt.Sprintf(`bzr revert -q -r "%s"`, version))
	} else if dirExists(".svn") {
		execCmd(fmt.Sprintf(`svn update -r "%s"`, version))
	}
}

func installPath(pkg string) string {
	return os.Getenv("GOPATH") + "/src/" + pkg
}

func mustChdir(path string) {
	err := os.Chdir(path)
	if err != nil {
		panic(err)
	}
}
