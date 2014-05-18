package main

import (
	"log"

	"github.com/mtibben/gogpm/vcs"
)

func goget(dep string) (string, error) {
	return execCmd("go", "get", "-d", "-u", dep)
}

// Iterates over Godep file dependencies and sets
// the repo revision to the version required version
func install() error {
	deps, err := readDepFile()
	if err != nil {
		return err
	}

	for dep, wantedVersion := range deps {
		pkg, err := vcs.PackageForImportPath(dep)
		if err != nil {
			return err
		}

		curVersion, _ := pkg.CurrentTagOrRevision()

		if curVersion == wantedVersion {
			log.Printf("Checked %s\n", dep)
		} else {

			log.Printf("Getting %s\n", dep)

			_, err := goget(dep)
			if err != nil {
				// Sometimes we get the error message
				//     package xxx: unrecognized import path "xxx"
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
	}

	log.Println("All Done")

	return nil
}
