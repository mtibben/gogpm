# Gogpm
A go dependency manager, originally a go implementation of [gpm](https://github.com/pote/gpm)

Features
 * picks up dependencies in tests
 * use args to specify which packages to collect deps for (e.g. ./...)

TODO:
 * get rid of bash requirements
 * clean up and simplify vcs
 * improve documentation
 * make depfile parser more flexible, allow comments, etc
 * read dep files from other dep managers (like godep)
