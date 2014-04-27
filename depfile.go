package main

import (
	"io/ioutil"
	"os"
	"strings"
)

func parseDepFile() map[string]string {
	b, err := ioutil.ReadFile(depsFile)
	if err != nil {
		panic(err)
	}

	deps := make(map[string]string, len(b))
	for _, line := range strings.Split(string(b), "\n") {
		if line == "" {
			continue
		}
		d := strings.Split(line, " ")

		if len(d) != 2 {
			panic("Couldn't parse depfile " + depsFile)
		}

		deps[d[0]] = d[1]
	}

	return deps
}

func appendLineToDepFile(line string) {
	f, err := os.OpenFile(depsFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString(line + "\n"); err != nil {
		panic(err)
	}
}
