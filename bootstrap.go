package main

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/mtibben/gogpm/vcs"
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

	depListStr, _, err := execCmd(`go list -f '{{join .Deps "\n"}}' | xargs go list -f '{{if not .Standard}}{{.ImportPath}}{{end}}'`)
	if err != nil {
		return err
	}

	depListStr = strings.TrimSpace(depListStr)
	dependencies := strings.Split(depListStr, "\n")

	deps := map[string]string{}

	for _, importPath := range dependencies {
		rr, err := vcs.RepoRootForImportPath(importPath)
		rootVcsPath := rr.Root

		if _, exists := deps[rootVcsPath]; exists {
			continue
		}

		absoluteVcsPath := installPath(rootVcsPath)

		if absoluteVcsPath == workingDir {
			continue
		}

		version, err := rr.Vcs.CurrentTag(absoluteVcsPath)
		if err != nil {
			return err
		}

		log.Printf(`Adding package "%s" version "%s"`, rootVcsPath, version)
		deps[rootVcsPath] = version
	}

	log.Printf("Writing Godeps file")
	for dep, version := range deps {
		appendLineToDepFile(fmt.Sprintf("%s %s", dep, version))
	}

	log.Println("All Done")

	return nil
}
