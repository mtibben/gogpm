package main

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
)

// Iterates over Godep file dependencies and sets
// the specified version on each of them.
func install() error {
	deps, err := readDepFile()
	if err != nil {
		return err
	}

	_, err = exec.LookPath("go")
	if err != nil {
		return errors.New("Go is currently not installed or in your PATH\n")
	}

	for dep, version := range deps {
		log.Printf("Getting %s\n", dep)
		_, _, err := execCmd(fmt.Sprintf(`go get -u -d "%s/..."`, dep))
		if err != nil {
			return err
		}

		err = setPackageToVersion(dep, version)
		if err != nil {
			return err
		}
	}

	log.Println("All Done")

	return nil
}
