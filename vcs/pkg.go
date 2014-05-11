package vcs

import (
	"os"
	"path/filepath"
)

type Package struct {
	rr *RepoRoot
}

func PackageFromImportPath(importPath string) (*Package, error) {
	reporoot, err := RepoRootForImportPath(importPath)
	if err != nil {
		return nil, err
	}

	return &Package{
		rr: reporoot,
	}, nil
}

func (p *Package) RootImportPath() string {
	return p.rr.Root
}

func (p *Package) Path() string {
	return filepath.Join(os.Getenv("GOPATH"), "src", p.rr.Root)
}

func (p *Package) CurrentRevision() (string, error) {
	return p.rr.Vcs.CurrentTag(p.Path())
}

func (p *Package) SetRevision(version string) error {
	return p.rr.Vcs.Checkout(p.Path(), version)
}
