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
	return filepath.Join(os.Getenv("GOPATH"), "src", p.rr.root)
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
