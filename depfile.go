package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

// dep => version
type dependencyMap map[string]string

func depfileExists() bool {
	fi, err := os.Stat(depsFile)
	return !os.IsNotExist(err) && !fi.IsDir()
}

func readDepFile() (dependencyMap, error) {
	if !depfileExists() {
		return nil, fmt.Errorf("%s file does not exist", depsFile)
	}

	b, err := ioutil.ReadFile(depsFile)
	if err != nil {
		return nil, err
	}

	deps := make(dependencyMap, len(b))
	lines := strings.Split(string(b), "\n")
	for i, line := range lines {

		// strip comments from line
		lineAndComment := strings.SplitN(line, "#", 2)
		line = lineAndComment[0]
		line = strings.TrimSpace(line)

		// ignore empty lines
		if line == "" {
			continue
		}

		// split on whitespace
		d := strings.Fields(line)

		if len(d) != 2 {
			return nil, fmt.Errorf("Couldn't parse line %d of depfile %s", i, depsFile)
		}

		deps[d[0]] = d[1]
	}

	return deps, nil
}

func writeDepFile(deps dependencyMap) error {

	// sort the dep paths
	klen := 0
	var keys []string
	for k := range deps {
		keys = append(keys, k)
		if len(k) > klen {
			klen = len(k)
		}
	}
	sort.Strings(keys)

	f, err := os.OpenFile(depsFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, dep := range keys {
		version := deps[dep]

		line := fmt.Sprintf("%-[3]*[1]s %[2]s", dep, version, klen+2)

		if _, err = f.WriteString(line + "\n"); err != nil {
			return err
		}
	}
	return nil
}
