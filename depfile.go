package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func parseDepFile() (map[string]string, error) {
	b, err := ioutil.ReadFile(depsFile)
	if err != nil {
		return nil, err
	}

	deps := make(map[string]string, len(b))
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

func appendLineToDepFile(line string) error {
	f, err := os.OpenFile(depsFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}

	defer f.Close()

	if _, err = f.WriteString(line + "\n"); err != nil {
		return err
	}

	return nil
}
