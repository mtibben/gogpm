package vcs

import (
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func (v *vcsCmd) currentRevision(dir string) (tag string, err error) {
	out, err := v.runOutput(dir, v.currentRevisionCmd.cmd)
	if err != nil {
		return
	}

	if v.currentRevisionCmd.pattern == "" {
		tag = strings.TrimSpace(string(out))
	} else {
		re := regexp.MustCompile(`(?m-s)` + v.currentRevisionCmd.pattern)
		m := re.FindStringSubmatch(tag)
		if len(m) > 1 {
			tag = m[1]
		} else {
			err = errors.New("Regex didn't match")
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

func (p *PackageRepo) CurrentRevision() (string, error) {
	return p.rr.vcs.currentRevision(p.Dir())
}

func (p *PackageRepo) SetRevision(version string) error {
	return p.rr.vcs.checkout(p.Dir(), version)
}
