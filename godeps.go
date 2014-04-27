package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func usage() {
	usagestring := `
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

      $ gpm bootstrap     # Downloads all external top-level packages required by
                          # your application and generates a Godeps file with their
                          # latest tags or revisions.

      $ gogpm version     # Outputs version information
      $ gogpm help        # Prints this message

`
	fmt.Fprintf(os.Stderr, usagestring)
	os.Exit(2)
}

var depsFile = "Godeps"

var workingDir string

func main() {
	// parse flags and opts
	flag.Parse()
	command := flag.Arg(0)

	var err error
	workingDir, err = os.Getwd()
	fmt.Println(workingDir)
	if err != nil {
		panic(err)
	}

	depsFile = workingDir + "/" + depsFile

	// Command Line Parsing
	switch command {
	case "version":
		fmt.Println(">> gogpm 0.1 (gpm v1.2.1 equiv)")

	case "install", "":
		if !fileExists(depsFile) {
			panic(fmt.Sprintf(">> %s file does not exist.\n", depsFile))
		}

		_, err := exec.LookPath("go")
		if err != nil {
			panic(">> Go is currently not installed or in your PATH\n")
		}

		install()

	case "bootstrap":
		if fileExists(depsFile) {
			panic(">> A Godeps file exists within this directory.")
		}

		bootstrap()

	default:
		usage()
	}
}

func fileExists(path string) bool {
	fi, err := os.Stat(path)
	if !os.IsNotExist(err) && !fi.IsDir() {
		return true
	}

	return false
}

func dirExists(path string) bool {
	fi, err := os.Stat(path)
	if !os.IsNotExist(err) && fi.IsDir() {
		return true
	}

	return false
}

func parseDepFile() map[string]string {
	b, err := ioutil.ReadFile(depsFile)
	if err != nil {
		panic(err)
	}

	deps := make(map[string]string, len(b))
	for _, line := range strings.Split(string(b), "\n") {
		if line == "" {
			continue
		}
		d := strings.Split(line, " ")

		if len(d) != 2 {
			panic("Couldn't parse " + depsFile)
		}

		deps[d[0]] = d[1]
	}

	return deps
}

func appendLineToDepFile(line string) {
	f, err := os.OpenFile(depsFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString(line + "\n"); err != nil {
		panic(err)
	}
}

func execCmd(cmd string) string {
	command := exec.Command("bash")
	var b bytes.Buffer
	command.Stdin = bytes.NewBufferString(cmd)
	command.Stdout = &b
	// command.Stderr = os.Stderr
	err := command.Run()
	if err != nil {
		panic(err)
	}

	return b.String()
}
