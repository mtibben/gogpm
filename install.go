package main

import (
	"log"

	"github.com/mtibben/gogpm/vcs"
)

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
			_, err := execCmd("go", "get", "-d", "-u", dep+"/...")
			if err != nil {
				return err
			}

			err = pkg.SetRevision(wantedVersion)
			if err != nil {
				return err
			}
		}
	}

	log.Println("All Done")

	return nil
}
