package vcs

import (
	"go/build"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestPresentPackageRepoSimpleGopath(t *testing.T) {
	withDummyBuildContext(t, func(gopath string) {
		repo, _ := PackageForImportPath("github.com/fake/library")
		expected := path.Join(gopath, "src", "/github.com/fake/library")
		actual := repo.Dir()
		if actual != expected {
			t.Errorf("expected Dir to be %v but it was %v", expected, actual)
		}
	})
}

func TestNotYetInstalledPackageRepoSimpleGopath(t *testing.T) {
	withDummyBuildContext(t, func(gopath string) {
		repo, _ := PackageForImportPath("github.com/fake/uninstalledlibrary")
		expected := path.Join(gopath, "src", "/github.com/fake/uninstalledlibrary")
		actual := repo.Dir()
		if actual != expected {
			t.Errorf("expected Dir to be %v but it was %v", expected, actual)
		}
	})
}

func withDummyBuildContext(t *testing.T, testFunc func(string)) {
	fakeGoPath, err := ioutil.TempDir("", "gopath")
	if err != nil {
		t.Fatal(err)
	}
	err = os.MkdirAll(path.Join(fakeGoPath, "src/github.com/fake/library"), 0777)
	if err != nil {
		t.Fatal(err)
	}

	context := build.Default
	context.GOPATH = fakeGoPath
	importContext = &context

	testFunc(fakeGoPath)

	importContext = nil
}
