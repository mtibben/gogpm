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
	if !fileExists(depsFile) {
		return fmt.Errorf("%s file does not exist", depsFile)
	}

	_, err := exec.LookPath("go")
	if err != nil {
		return errors.New("Go is currently not installed or in your PATH\n")
	}

	pkgs, err := parseDepFile()
	if err != nil {
		return err
	}

	for pkg, version := range pkgs {
		log.Printf("Getting %s\n", pkg)
		_, _, err = execCmd(fmt.Sprintf(`go get -u -d "%s"`, pkg))
		if err != nil {
			return err
		}

		err = setPackageToVersion(pkg, version)
		if err != nil {
			return err
		}
	}

	log.Println("All Done")

	return nil
}
