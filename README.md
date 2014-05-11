# Gogpm

Gogpm is a go implementation of [gpm](https://github.com/pote/gpm)

### gogpm

gogpm is a tool for managing package dependency versions

gogpm leverages the power of the `go get` command and the underlying version
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


### Usage:

    $ gogpm bootstrap [packages]    # Downloads all top-level packages required by the listed
                                    # import paths and generates a Godeps file with their
                                    # latest tags or revisions.
                                    # For more about specifying packages, see 'go help packages'.

    $ gogpm install                 # Parses the Godeps file, installs dependencies and sets
                                    # them to the appropriate version.

    $ gogpm version                 # Outputs version information

    $ gogpm help                    # Prints this message


### Example usage

```bash
$ ls .
Godeps  foo.go  foo_test.go

$ cat Godeps
github.com/nu7hatch/gotrail               v0.0.2
github.com/replicon/fast-archiver         v1.02
launchpad.net/gocheck                     r2013.03.03   # Bazaar repositories are supported
code.google.com/p/go.example/hello/...    ae081cd1d6cc  # And so are Mercurial ones

$ gpm install
>> Getting github.com/nu7hatch/gotrail
>> Setting github.com/nu7hatch/gotrail to version v0.0.2
>> Getting code.google.com/p/go.example/hello/...
>> Setting code.google.com/p/go.example/hello/... to version ae081cd1d6cc
>> Getting launchpad.net/gocheck
>> Setting launchpad.net/gocheck to version r2013.03.03
>> Getting github.com/replicon/fast-archiver
>> Setting github.com/replicon/fast-archiver to version v1.02
>> All Done
```


## Goals
 * Keep things small and simple
 * Use `go get` wherever possible
 * No external dependencies
 * Use `go get`'s own logic for package import path resolution (the vcs package)


## Differences from gpm
 * Gogpm finds dependencies in tests
 * The bootstrap command takes optional arguments for packages so that you can create a Godeps file with dependencies from multiple packages
 * When running "bootstrap", gogpm does not update the package if it already exists
 * Gogpm does not implicitly run the install command when no arguments are given. You'll need to `gogpm install`
