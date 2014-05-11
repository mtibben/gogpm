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
	version      = "gogpm 1.0-pre (gpm v1.2.1 equiv)"
	lockfileName = "Godeps"
	usage        = `gogpm is a tool for managing package dependency versions

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


Usage:

    $ gogpm bootstrap [packages]    # Downloads all top-level packages required by the listed
                                    # import paths and generates a Godeps file with their
                                    # latest tags or revisions.
                                    # For more about specifying packages, see 'go help packages'.

    $ gogpm install                 # Parses the Godeps file, installs dependencies and sets
                                    # them to the appropriate version.

    $ gogpm version                 # Outputs version information

    $ gogpm help                    # Prints this message

`
)

var depsFile, workingDir string

func init() {
	var err error

	log.SetPrefix(">> ")
	log.SetFlags(0)

	// get the working directory
	workingDir, err = os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error()+"\n")
		os.Exit(1)
	}

	// lockfile name
	depsFile = filepath.Join(workingDir, lockfileName)
}

func main() {
	var err error

	// parse flags and opts
	flag.Parse()
	command := flag.Arg(0)

	// Command Line Parsing
	switch command {

	case "version":
		fmt.Println(version)

	case "bootstrap":
		args := flag.Args()
		err = bootstrap(args[1:])

	case "install":
		err = install()

	default:
		err = errors.New(usage)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error()+"\n")
		os.Exit(1)
	}
}
