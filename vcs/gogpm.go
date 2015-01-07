package vcs

import (
	"errors"
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
	rr *repoRoot
}

func PackageForImportPath(importPath string) (*PackageRepo, error) {
	reporoot, err := repoRootForImportPath(importPath)
	if err != nil {
		return nil, err
	}

	return &PackageRepo{
		rr: reporoot,
	}, nil
}

func (p *PackageRepo) RootImportPath() string {
	return p.rr.root
}

func (p *PackageRepo) Dir() string {
	// split gopath
	goPath := os.Getenv("GOPATH")
	paths := strings.Split(goPath, ":")

	// construct path options for repo
	fullPaths := []string{}
	for _, path := range paths {
		fullPaths = append(fullPaths, filepath.Join(path, "src", p.rr.root))
	}

	// return first instance where lib exists
	for _, path := range fullPaths {
		_, err := os.Stat(path)
		if err == nil {
			return path
		}
	}

	// if not installed, put it in FIRST gopath
	return fullPaths[0]
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
