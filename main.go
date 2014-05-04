package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const (
	lockfile = "Godeps"
	usage    = `
SYNOPSIS

    gogpm leverages the power of the go get command and the underlying version
    control systems used by it to set your Go dependencies to desired versions,
    thus allowing easily reproducible builds in your Go projects.

    A Godeps file in the root of your Go application is expected containing
    the import paths of your packages and a specific tag or commit hash
    from its version control system, an example Godeps file looks like this:

    $ cat Godeps
    # This is a comment
    github.com/nu7hatch/gotrail         v0.0.2
    github.com/replicon/fast-archiver   v1.02   #This is another comment!
    github.com/nu7hatch/gotrail         2eb79d1f03ab24bacbc32b15b75769880629a865

USAGE
      $ gogpm             # Same as 'install'.

      $ gogpm install     # Parses the Godeps file, installs dependencies and sets
                          # them to the appropriate version.

      $ gogpm bootstrap   # Downloads all external top-level packages required by
                          # your application and generates a Godeps file with their
                          # latest tags or revisions.

      $ gogpm version     # Outputs version information
      $ gogpm help        # Prints this message

`
)

var depsFile, workingDir string

func init() {
	log.SetPrefix(">> ")
	log.SetFlags(0)
}

func main() {
	// parse flags and opts
	flag.Parse()
	command := flag.Arg(0)

	// get the working directory
	var err error
	workingDir, err = os.Getwd()
	if err != nil {
		panic(err)
	}

	// lock file name
	depsFile = filepath.Clean(filepath.Join(workingDir, lockfile))

	// Command Line Parsing
	switch command {
	case "version":
		fmt.Println("gogpm 0.1-alpha (gpm v1.2.1 equiv)")

	case "bootstrap":
		err = bootstrap()

	case "install", "":
		err = install()

	default:
		err = errors.New(usage)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error()+"\n")
		os.Exit(1)
	}
}
