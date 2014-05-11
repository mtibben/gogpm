# Gogpm

A go implementation of [gpm](https://github.com/pote/gpm)


## Goals
 * Keep things small and simple
 * Use `go get` wherever possible
 * No external dependencies
 * Use `go get`'s own logic for package import path resolution (the vcs package)


## Differences from gpm
 * Gogpm finds dependencies in tests (gpm does not)
 * The bootstrap command takes optional arguments for packages so that you can create a Godeps file with dependencies from multiple packages
 * When running "bootstrap", gogpm does not update the package if it already exists
 * No implicit command for gogpm. You'll need to `gogpm install`


## Example usage

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
