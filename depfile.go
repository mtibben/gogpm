package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// dep => version
type dependencyMap map[string]string

func readDepFile() (dependencyMap, error) {
	if !fileExists(depsFile) {
		return nil, fmt.Errorf("%s file does not exist", depsFile)
	}

	b, err := ioutil.ReadFile(depsFile)
	if err != nil {
		return nil, err
	}

	deps := make(dependencyMap, len(b))
	for _, line := range strings.Split(string(b), "\n") {
		if line == "" {
			continue
		}
		d := strings.Split(line, " ")

		if len(d) != 2 {
			return nil, fmt.Errorf("Couldn't parse depfile " + depsFile)
		}

		deps[d[0]] = d[1]
	}

	return deps, nil
}

func writeDepFile(deps dependencyMap) error {
	f, err := os.OpenFile(depsFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	for dep, version := range deps {
		line := fmt.Sprintf("%s %s", dep, version)

		if _, err = f.WriteString(line + "\n"); err != nil {
			return err
		}
	}
	return nil
}
