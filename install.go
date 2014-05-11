package main

import (
	"log"

	"github.com/mtibben/gogpm/vcs"
)

// Iterates over Godep file dependencies and sets
// the specified version on each of them.
func install() error {
	deps, err := readDepFile()
	if err != nil {
		return err
	}

	for dep, wantedVersion := range deps {
		curVersion := ""
		rr, err := vcs.RepoRootForImportPath(dep)
		if err == nil {
			absoluteVcsPath := installPath(rr.Root)
			curVersion, _ = rr.Vcs.CurrentTag(absoluteVcsPath)
		}

		if curVersion == wantedVersion {
			log.Printf("Checked %s\n", dep)
		} else {
			log.Printf("Getting %s\n", dep)
			_, err := execCmd("go", "get", "-d", "-u", dep+"/...")
			if err != nil {
				return err
			}

			err = setPackageToVersion(dep, wantedVersion)
			if err != nil {
				return err
			}
		}
	}

	log.Println("All Done")

	return nil
}
