package main

import (
	"flag"
	"fmt"
	"io/ioutil"
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
    github.com/replicon/fast-archiver   v1.02    # Tag
    github.com/nu7hatch/gotrail         2eb79d1f # Revisions


Usage: gogpm [-v] <command>

    $ gogpm [-v] bootstrap [packages]  # Downloads and installs the packages
                                       # named by the import paths along with
                                       # their dependencies (executes
                                       # go get -d [packages]).
                                       #
                                       # Generates a Godeps file with the
                                       # package's current tags or revisions.
                                       # For more about specifying packages,
                                       # see 'go help packages'.

    $ gogpm [-v] install               # Parses the Godeps file, installs
                                       # dependencies and sets them to the
                                       # appropriate version.

    $ gogpm version                    # Outputs version information

    $ gogpm help                       # Prints this message

The -v flag makes the output more verbose.
`
)

var depsFile, workingDir string

var logErr = log.New(os.Stderr, "", 0)
var logVerbose = log.New(ioutil.Discard, "", 0)

func initLogging(verbose bool) {
	log.SetFlags(0)
	log.SetPrefix("")
	if verbose {
		logVerbose = log.New(os.Stdout, "", 0)
	}
}

func init() {
	var err error

	// parse flags and opts
	flag.Usage = func() {
		logErr.Println(usage)
		os.Exit(1)
	}
	// flag.BoolVar(verbose, "verbose", false, "Verbose")
	verbose := flag.Bool("v", false, "Verbose")
	flag.Parse()

	initLogging(*verbose)

	// get the working directory
	workingDir, err = os.Getwd()
	if err != nil {
		logErr.Fatalln(err.Error())
	}

	// lockfile name
	depsFile = filepath.Join(workingDir, lockfileName)
}

func main() {
	var err error

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
		flag.Usage()
	}

	if err != nil {
		logErr.Fatalln(err.Error())
	}
}
