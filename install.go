package main

import (
	"errors"
	"fmt"
	"log"
	"os/exec"

	"github.com/mtibben/gogpm/vcs"
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

		update := ""
		rr, err := vcs.RepoRootForImportPath(dep)
		if err == nil {
			absoluteVcsPath := installPath(rr.Root)
			cur, _ := rr.Vcs.CurrentTag(absoluteVcsPath)
			if cur != version {
				update = "-u"
			}
		}

		out, err := execCmd(fmt.Sprintf(`go get -d %s "%s/..."`, update, dep))
		if err != nil {
			log.Println(out)
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
