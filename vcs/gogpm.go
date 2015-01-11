package vcs

import (
	"errors"
	"go/build"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func (v *vcsCmd) currentRevision(dir string) (revision string, err error) {
	out, err := v.runOutput(dir, v.currentRevisionCmd.cmd)
	if err != nil {
		return
	}

	if v.currentRevisionCmd.pattern == "" {
		revision = strings.TrimSpace(string(out))
	} else {
		re := regexp.MustCompile(`(?m-s)` + v.currentRevisionCmd.pattern)
		m := re.FindStringSubmatch(revision)
		if len(m) > 1 {
			revision = m[1]
		} else {
			err = errors.New("Regex didn't match")
		}
	}

	return
}

func (v *vcsCmd) currentTag(dir, currentRevision string) (tag string, err error) {
	out, err := v.runOutput(dir, v.currentTagCmd.cmd, "curRev", currentRevision)
	if err != nil {
		return
	}

	if v.currentTagCmd.pattern == "" {
		tag = strings.TrimSpace(string(out))
	} else {
		re := regexp.MustCompile(`(?m-s)` + v.currentTagCmd.pattern)
		m := re.FindStringSubmatch(tag)
		if len(m) > 1 {
			tag = m[1]
		}
	}

	return
}

func (v *vcsCmd) checkout(dir, tag string) error {
	_, err := v.runOutput(dir, v.tagSyncCmd, "tag", tag)
	return err
}

type PackageRepo struct {
	rr  *repoRoot
	ctx build.Context
}

func PackageForImportPath(importPath string) (*PackageRepo, error) {
	reporoot, err := repoRootForImportPath(importPath)
	if err != nil {
		return nil, err
	}

	return &PackageRepo{
		rr:  reporoot,
		ctx: build.Default,
	}, nil
}

func (p *PackageRepo) RootImportPath() string {
	return p.rr.root
}

func (p *PackageRepo) Dir() string {
	// get go/build to search GOPATH for where it is installed
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	pkg, err := p.ctx.Import(p.rr.root, pwd, build.FindOnly)
	if err == nil {
		return pkg.Dir
	}

	// TODO when ready to drop support for go 1.3 & under
	// we should add a check for if the err is a build.MultiplePackageError

	// if uninstalled go get/gogpm will put in first GOPATH entry
	goPath := p.ctx.GOPATH
	paths := filepath.SplitList(goPath)

	if len(paths) == 0 {
		panic("GOPATH not defined")
	}

	return filepath.Join(paths[0], "src", p.rr.root)
}

// IsCurrentTagOrRevision checks if the given version matches
// the current tag or revision
func (p *PackageRepo) IsCurrentTagOrRevision(version string) bool {
	dir := p.Dir()

	currentRev, _ := p.rr.vcs.currentRevision(dir)
	if currentRev == version {
		return true
	}

	currentTag, _ := p.rr.vcs.currentTag(dir, currentRev)
	return currentTag == version
}

func (p *PackageRepo) CurrentTagOrRevision() (string, error) {
	dir := p.Dir()

	currentRev, err := p.rr.vcs.currentRevision(dir)
	if err != nil {
		return "", err
	}
	currentTag, err := p.rr.vcs.currentTag(dir, currentRev)
	if err != nil {
		return "", err
	}

	if currentTag != "" {
		return currentTag, nil
	} else {
		return currentRev, nil
	}
}

func (p *PackageRepo) SetRevision(version string) error {
	return p.rr.vcs.checkout(p.Dir(), version)
}
