package main

import (
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"sort"
	"strings"

	"github.com/mtibben/gogpm/vcs"
)

// uniq filters a string slice for unique vals
func uniq(set []string) (uniqSet []string) {
	sort.Strings(set)

	lastitem := ""
	for _, item := range set {
		if item != lastitem {
			uniqSet = append(uniqSet, item)
			lastitem = item
		}
	}

	return uniqSet
}

func bootstrap(packages []string) error {

	// check if Godeps already exists
	if depfileExists() {
		return errors.New("A Godeps file already exists within this directory")
	}

	log.Println("Installing dependencies")

	// go get dependencies if they're not already present (without updating)
	_, err := execCmd("go", append([]string{"get", "-d"}, packages...)...)
	if err != nil {
		return err
	}

	// get a list of dependencies for the packages, including test dependencies
	tmpl := fmt.Sprintf(`{{ join .Deps "\n" }}%s{{ join .TestImports "\n" }}`, "\n")
	depListStr, err := execCmd("go", append([]string{"list", "-f", tmpl}, packages...)...)
	if err != nil {
		return err
	}
	depList := uniq(strings.Split(strings.TrimSpace(depListStr), "\n"))

	// filter out standard library
	depListStr, err = execCmd("go", append([]string{"list", "-f", "{{if not .Standard}}{{.ImportPath}}{{end}}"}, depList...)...)
	if err != nil {
		return err
	}

	dependencies := uniq(strings.Split(strings.TrimSpace(depListStr), "\n"))
	dependencies = append(dependencies, packages...)

	deps := dependencyMap{}

	for _, importPath := range dependencies {
		// ignore relative import paths
		if importPath[0] == '.' {
			continue
		}

		// get the repo for the import path
		pkg, err := vcs.PackageFromImportPath(importPath)
		if err != nil {
			return err
		}

		// if dep has already been found
		if _, exists := deps[pkg.RootImportPath()]; exists {
			continue
		}

		// if the dep is the directory we're working from
		pp, err := filepath.Abs(pkg.Path())
		if err != nil {
			return err
		}
		wd, _ := filepath.Abs(workingDir)
		if err != nil {
			return err
		}
		if pp == wd {
			continue
		}

		// current revision of repo
		version, err := pkg.CurrentRevision()
		if err != nil {
			return err
		}

		log.Printf(`Adding package "%s" version "%s"`, pkg.RootImportPath(), version)
		deps[pkg.RootImportPath()] = version
	}

	log.Printf("Writing Godeps file")
	err = writeDepFile(deps)
	if err != nil {
		return err
	}

	log.Println("All Done")

	return nil
}
