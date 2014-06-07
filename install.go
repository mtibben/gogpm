package main

import "github.com/mtibben/gogpm/vcs"

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

		if pkg.IsCurrentTagOrRevision(wantedVersion) {
			logVerbose.Printf("Checked %s\n", dep)
		} else {

			logVerbose.Printf("Getting %s\n", dep)

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

			logVerbose.Printf("Setting %s to %s\n", dep, wantedVersion)

			err = pkg.SetRevision(wantedVersion)
			if err != nil {
				return err
			}
		}
	}

	logVerbose.Println("All Done")

	return nil
}
