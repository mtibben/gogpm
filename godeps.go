package main

import (
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
      $ gogpm version     # Outputs version information
      $ gogpm help        # Prints this message

`
	fmt.Fprintf(os.Stderr, usagestring)
	os.Exit(2)
}

const depsFile = "Godeps"

func main() {
	// parse flags and opts
	command := flag.Arg(0)

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

		setDependencies(depsFile)

	default:
		usage()
	}
}

// Iterates over Godep file dependencies and sets
// the specified version on each of them.
func setDependencies(depFile string) {
	for pkg, version := range parseDepFile(depFile) {
		installPath := os.Getenv("GOPATH") + "/src/" + pkg
		fmt.Printf(">> Getting package %s\n", pkg)

		execCmd(fmt.Sprintf(`go get -u -d "%s"`, pkg))
		fmt.Printf(">> Setting %s to version %s\n", pkg, version)

		err := os.Chdir(installPath)
		if err != nil {
			panic(err)
		}

		if dirExists(".hg") {
			execCmd(fmt.Sprintf(`hg update -q "%s"`, version))
		} else if dirExists(".git") {
			execCmd(fmt.Sprintf(`git checkout -q "%s"`, version))
		} else if dirExists(".bzr") {
			execCmd(fmt.Sprintf(`bzr revert -q -r "%s"`, version))
		} else if dirExists(".svn") {
			execCmd(fmt.Sprintf(`svn update -r "%s"`, version))
		}
	}

	fmt.Println(">> All Done")
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

func parseDepFile(depFile string) map[string]string {
	b, err := ioutil.ReadFile(depFile)
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

func execCmd(cmd string) {
	command := exec.Command("bash", "-c", cmd)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	err := command.Run()
	if err != nil {
		panic(err)
	}
}
