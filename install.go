package main

import (
	"errors"
	"log"
	"strings"

	"github.com/mtibben/gogpm/vcs"
)

func goget(dep string) (string, error) {
	return execCmd("go", "get", "-d", "-u", dep)
}

func installDep(dep, wantedVersion string) error {
	pkg, err := vcs.PackageForImportPath(dep)
	if err != nil {
		return err
	}

	if pkg.IsCurrentTagOrRevision(wantedVersion) {
		log.Printf("Checked %s\n", dep)
	} else {

		log.Printf("Getting %s\n", dep)

		// we need the /... else sometimes we get
		// 	   imports github.com/bradfitz/gomemcache: no buildable Go source files in /go/src/github.com/bradfitz/gomemcache
		_, err := goget(dep + "/...")
		if err != nil {
			// Sometimes we get the error message
			//     package gopkg.in/check.v1/...: unrecognized import path "gopkg.in/check.v1/..."
			// It seems the package is downloaded however, so running the command
			// again returns without an error
			_, err := goget(dep)
			if err != nil {
				return err
			}
		}

		log.Printf("Setting %s to %s\n", dep, wantedVersion)

		err = pkg.SetRevision(wantedVersion)
		if err != nil {
			return err
		}
	}
	return nil
}

func combineErrs(errs []error) error {
	errStrings := []string{}
	for _, e := range errs {
		if e != nil {
			errStrings = append(errStrings, e.Error())
		}
	}

	if len(errStrings) > 0 {
		return errors.New(strings.Join(errStrings, "\n"))
	}

	return nil
}

// Iterates over Godep file dependencies and sets
// the repo revision to the version required version
func install() error {
	deps, err := readDepFile()
	if err != nil {
		return err
	}
	errs := []error{}

	if concurrencyDisabled {
		for dep, wantedVersion := range deps {
			err := installDep(dep, wantedVersion)
			errs = append(errs, err)
		}
	} else {
		errChan := make(chan error)
		active := 0

		for dep, wantedVersion := range deps {
			active++
			go func(dep, wantedVersion string) {
				errChan <- installDep(dep, wantedVersion)
			}(dep, wantedVersion)
		}

		for {
			err := <-errChan
			errs = append(errs, err)
			active--
			if active == 0 {
				break
			}
		}
	}

	err = combineErrs(errs)
	if err == nil {
		log.Println("All Done")
	}

	return err

}
