package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func setPackageToVersion(pkg, version string) {
	dir := installPath(pkg)
	mustChdir(dir)
	vcs, err := vcsForDir(dir)
	if err != nil {
		fmt.Printf(">> Ignored %s, not top-level package.\n", pkg)
	} else {
		fmt.Printf(">> Setting %s to version %s\n", pkg, version)
		vcs.TagSync(version)
	}
}

func installPath(pkg string) string {
	return filepath.Clean(os.Getenv("GOPATH") + "/src/" + pkg)
}

// Returns the latest tag (or, failing that latest revision)
// for an installed package.
func findLastTagOrHEAD(pkg string) (string, error) {
	dir := installPath(pkg)
	mustChdir(dir)
	vcs, err := vcsForDir(dir)
	if err != nil {
		return "", err
	}

	return vcs.CurrentTag(), nil
}
