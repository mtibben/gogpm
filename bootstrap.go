package main

import (
	"fmt"
	"strings"
)

func bootstrap() {
	if fileExists(depsFile) {
		panic(">> A Godeps file exists within this directory.")
	}

	fmt.Println(">> Installing dependencies.")
	execCmd("go get -d")

	depString := execCmd(`go list -f '{{join .Deps "\n"}}' | xargs go list -f '{{if not .Standard}}{{.ImportPath}}{{end}}'`)
	dependencies := strings.Split(depString, "\n")

	for _, pkg := range dependencies {
		version, err := findLastTagOrHEAD(pkg)

		// If no repo file is found it means we are inside a repo's
		// subdirectory tree, we can just ignore this package.
		if err != nil {
			fmt.Printf(">> Ignored %s, not top-level package.\n", pkg)
			continue
		}

		fmt.Printf(`>> Adding package "%s" version "%s" to Godeps.`, pkg, version)
		fmt.Println()
		appendLineToDepFile(fmt.Sprintf("%s %s", pkg, version))
		// Sets a given package to a given revision using
		// the appropriate VCS.
		setPackageToVersion(pkg, version)

	}
	fmt.Println(">> All Done.")
}
