package main

import (
	"fmt"
	"strings"
)

func bootstrap() {
	fmt.Println(">> Installing dependencies.")
	execCmd("go get -d")

	depString := execCmd(`go list -f '{{join .Deps "\n"}}' | xargs go list -f '{{if not .Standard}}{{.ImportPath}}{{end}}'`)
	dependencies := strings.Split(depString, "\n")

	for _, pkg := range dependencies {
		version := findLastTagOrHEAD(pkg)

		if version != "" {
			fmt.Printf(`>> Adding package "%s" version "%s" to Godeps.`, pkg, version)
			fmt.Println()
			appendLineToDepFile(fmt.Sprintf("%s %s", pkg, version))
			// Sets a given package to a given revision using
			// the appropriate VCS.
			setPackageToVersion(pkg, version)
		}

	}
	fmt.Println(">> All Done.")
}

// Returns the latest tag (or, failing that latest revision)
// for an installed package.
func findLastTagOrHEAD(pkg string) string {
	mustChdir(installPath(pkg))

	if dirExists(".git") {

		// FIXME: there should be a better way,  but git tags returns in alphabetical order.
		version := strings.TrimSpace(execCmd(`git tag | xargs -I@ git log --format=format:"%ai @%n" -1 @ | sort | awk '{print $4}' | tail -1`))

		if version != "" {
			return version
		}

		version = strings.TrimSpace(execCmd(`git log -n 1 --pretty=oneline | cut -d " " -f 1`))
		fmt.Printf(`>> No tags on package "%s", setting version to latest revision.`, pkg)
		fmt.Println()
		return version

	} else if dirExists(".bzr") {

		version := strings.TrimSpace(execCmd(`bzr tags | tail -1 | cut -d " " -f 1`))

		if version != "" {
			return version
		}

		version = strings.TrimSpace(execCmd(`bzr log -r-1 --log-format=line | cut -d ":" -f 1`))
		fmt.Printf(`>> No tags on package "%s", setting version to latest revision.`, pkg)
		fmt.Println()
		return version

	} else if dirExists(".hg") {

		version := strings.TrimSpace(execCmd(`hg parents --template "{latesttag}"`))

		if version != "" {
			return version
		}

		version = strings.TrimSpace(execCmd(`hg log --template "{node}" -l 1`))
		fmt.Printf(`>> No tags on package "%s", setting version to latest_revision.`, pkg)
		fmt.Println()
		return version
	} else {
		// If no repo file is found it means we are inside a repo's
		// subdirectory tree, we can just ignore this package.

		fmt.Printf(">> Ignored %s, not top-level package.", pkg)
		fmt.Println()
	}

	return ""
}
