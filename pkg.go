package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/mtibben/gogpm/vcs"
)

func setPackageToVersion(importpath, version string) error {
	rr, err := vcs.RepoRootForImportPath(importpath)

	if err != nil {
		log.Printf("Ignored %s, not top-level package\n", importpath)
		return nil
	} else {
		log.Printf("Setting %s to version %s\n", importpath, version)
		return rr.Vcs.Checkout(installPath(rr.Root), version)
	}
}

func installPath(importpath string) string {
	return filepath.Clean(filepath.Join(os.Getenv("GOPATH"), "src", importpath))
}
