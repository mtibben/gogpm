package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type vcsCmd struct {
	name string
	cmd  string

	latestTagCmd []string
	tagSyncCmd   string
}

func (v *vcsCmd) LatestTag() (tag string) {
	for _, c := range v.latestTagCmd {
		tag = strings.TrimSpace(execCmd(c))
		if tag != "" {
			break
		}
	}

	return
}

func (v *vcsCmd) TagSync(tag string) {
	execCmd(fmt.Sprintf(v.tagSyncCmd, tag))
}

var vcsList = []*vcsCmd{
	vcsHg,
	vcsGit,
	vcsBzr,
	vcsSvn,
}

var vcsHg = &vcsCmd{
	name: "Mercurial",
	cmd:  "hg",
	latestTagCmd: []string{
		`hg parents --template "{latesttag}"`,
		`hg log --template "{node}" -l 1`,
	},
	tagSyncCmd: `hg update -q "%s"`,
}

var vcsGit = &vcsCmd{
	name: "Git",
	cmd:  "git",
	latestTagCmd: []string{
		`git tag | xargs -I@ git log --format=format:"%ai @%n" -1 @ | sort | awk '{print $4}' | tail -1`,
		`git log -n 1 --pretty=oneline | cut -d " " -f 1`,
	},
	tagSyncCmd: `git checkout -q "%s"`,
}

var vcsBzr = &vcsCmd{
	name: "Bazaar",
	cmd:  "bzr",
	latestTagCmd: []string{
		`bzr tags | tail -1 | cut -d " " -f 1`,
		`bzr log -r-1 --log-format=line | cut -d ":" -f 1`,
	},
	tagSyncCmd: `bzr revert -q -r "%s"`,
}

var vcsSvn = &vcsCmd{
	name:       "Subversion",
	cmd:        "svn",
	tagSyncCmd: `svn update -r "%s"`,
}

func vcsForDir(origDir string) (vcs *vcsCmd, err error) {
	dir := filepath.Clean(origDir)

	for _, vcs := range vcsList {
		if fi, err := os.Stat(filepath.Join(dir, "."+vcs.cmd)); err == nil && fi.IsDir() {
			return vcs, nil
		}
	}

	return nil, fmt.Errorf("directory %q is not using a known version control system", origDir)
}
