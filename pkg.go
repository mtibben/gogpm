package main

import (
	"log"
	"os"
	"path/filepath"
)

func setPackageToVersion(pkg, version string) error {
	dir := installPath(pkg)

	err := os.Chdir(dir)
	if err != nil {
		return err
	}

	vcs, err := vcsForDir(dir)
	if err != nil {
		log.Printf("Ignored %s, not top-level package\n", pkg)
		return nil
	} else {
		log.Printf("Setting %s to version %s\n", pkg, version)
		return vcs.TagSync(version)
	}
}

func installPath(pkg string) string {
	return filepath.Clean(os.Getenv("GOPATH") + "/src/" + pkg)
}

// Returns the latest tag (or, failing that latest revision)
// for an installed package.
func findLastTagOrHEAD(pkg string) (string, error) {
	dir := installPath(pkg)

	err := os.Chdir(dir)
	if err != nil {
		return "", err
	}

	vcs, err := vcsForDir(dir)
	if err != nil {
		return "", err
	}

	return vcs.LatestTag()
}
