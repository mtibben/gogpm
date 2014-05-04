package main

import (
	"errors"
	"fmt"
	"log"
	"strings"
)

func bootstrap() error {
	if fileExists(depsFile) {
		return errors.New("A Godeps file already exists within this directory")
	}

	log.Println("Installing dependencies")
	_, _, err := execCmd("go get -d")
	if err != nil {
		return err
	}

	depString, _, err := execCmd(`go list -f '{{join .Deps "\n"}}' | xargs go list -f '{{if not .Standard}}{{.ImportPath}}{{end}}'`)
	if err != nil {
		return err
	}

	dependencies := strings.Split(depString, "\n")

	for _, pkg := range dependencies {
		version, err := findLastTagOrHEAD(pkg)

		// If no repo file is found it means we are inside a repo's
		// subdirectory tree, we can just ignore this package.
		if err != nil {
			log.Printf("Ignored %s, not top-level package\n", pkg)
			continue
		}

		log.Printf(`Adding package "%s" version "%s" to Godeps`, pkg, version)
		appendLineToDepFile(fmt.Sprintf("%s %s", pkg, version))
		// Sets a given package to a given revision using
		// the appropriate VCS.
		err = setPackageToVersion(pkg, version)
		if err != nil {
			return err
		}

	}
	log.Println("All Done")

	return nil
}
