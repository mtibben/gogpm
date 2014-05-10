package main

import (
	"errors"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/mtibben/gogpm/vcs"
)

func uniq(set []string) (newset []string) {
	sort.Strings(set)

	lastitem := ""
	for _, item := range set {
		if item != lastitem {
			newset = append(newset, item)
			lastitem = item
		}
	}

	return newset
}

func bootstrap(pgs []string) error {

	packages := strings.Join(pgs, " ")

	// check if Godeps already exists
	if fileExists(depsFile) {
		return errors.New("A Godeps file already exists within this directory")
	}

	log.Println("Installing dependencies")

	// go get dependencies if they're not already present (without updating)
	_, err := execCmd("go get -d " + packages)
	if err != nil {
		return err
	}

	// get a list of dependencies for the packages, including test dependencies
	depListStr, err := execCmd(fmt.Sprintf(`go list -f '{{join .Deps "\n"}}%s{{ join .TestImports "\n" }}' %s`, "\n", packages))
	if err != nil {
		log.Println(depListStr)
		return err
	}
	depList := uniq(strings.Split(strings.TrimSpace(depListStr), "\n"))

	// filter out standard library
	depListStr, err = execCmd(`go list -f '{{if not .Standard}}{{.ImportPath}}{{end}}' ` + strings.Join(depList, " "))
	if err != nil {
		log.Println(depListStr)
		return err
	}

	dependencies := append(uniq(strings.Split(strings.TrimSpace(depListStr), "\n")), pgs...)

	deps := dependencyMap{}

	for _, importPath := range dependencies {
		// ignore relative import paths
		if importPath[0] == '.' {
			continue
		}

		// get the repo for the import path
		rr, err := vcs.RepoRootForImportPath(importPath)
		if err != nil {
			return err
		}

		rootVcsPath := rr.Root

		// if dep has already been found
		if _, exists := deps[rootVcsPath]; exists {
			continue
		}

		// find the directory location of the repo
		absoluteVcsPath := installPath(rootVcsPath)

		// if the dep is the directory we're working from
		if absoluteVcsPath == workingDir {
			continue
		}

		// current version of repo
		version, err := rr.Vcs.CurrentTag(absoluteVcsPath)
		if err != nil {
			return err
		}

		log.Printf(`Adding package "%s" version "%s"`, rootVcsPath, version)
		deps[rootVcsPath] = version
	}

	log.Printf("Writing Godeps file")
	err = writeDepFile(deps)
	if err != nil {
		return err
	}

	log.Println("All Done")

	return nil
}
