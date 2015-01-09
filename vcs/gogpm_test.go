package vcs

import (
	"go/build"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"
)

func TestPresentPackageRepoSimpleGopath(t *testing.T) {
	withDummyBuildContextSingleGopath(t, func(gopath string) {
		repo, _ := PackageForImportPath("github.com/fake/library")
		context := build.Default
		context.GOPATH = gopath
		repo.ctx = &context

		expected := path.Join(gopath, "src", "/github.com/fake/library")
		actual := repo.Dir()
		if actual != expected {
			t.Errorf("expected Dir to be %v but it was %v", expected, actual)
		}
	})
}

func TestNotYetInstalledPackageRepoSimpleGopath(t *testing.T) {
	withDummyBuildContextSingleGopath(t, func(gopath string) {
		repo, _ := PackageForImportPath("github.com/fake/uninstalledlibrary")
		context := build.Default
		context.GOPATH = gopath
		repo.ctx = &context

		expected := path.Join(gopath, "src", "/github.com/fake/uninstalledlibrary")
		actual := repo.Dir()
		if actual != expected {
			t.Errorf("expected Dir to be %v but it was %v", expected, actual)
		}
	})
}

func withDummyBuildContextSingleGopath(t *testing.T, testFunc func(string)) {
	// setup
	fakeGoPath, err := ioutil.TempDir("", "gopath")
	if err != nil {
		t.Fatal(err)
	}
	err = os.MkdirAll(path.Join(fakeGoPath, "src/github.com/fake/library"), 0777)
	if err != nil {
		t.Fatal(err)
	}

	// run test
	testFunc(fakeGoPath)

	// cleanup
	os.RemoveAll(fakeGoPath)
}

func TestMultipleGopathSingleInstall(t *testing.T) {
	withDummyBuildContextMultipleGopath(t, func(gopathOne string, gopathTwo string) {
		repo, _ := PackageForImportPath("github.com/fake/library")
		context := build.Default
		context.GOPATH = strings.Join([]string{gopathOne, gopathTwo}, ":")
		repo.ctx = &context

		expected := path.Join(gopathTwo, "src", "/github.com/fake/library")
		actual := repo.Dir()
		if actual != expected {
			t.Errorf("expected Dir to be %v but it was %v", expected, actual)
		}
	})
}

func TestMultipleGopathNoInstall(t *testing.T) {
	withDummyBuildContextMultipleGopath(t, func(gopathOne string, gopathTwo string) {
		repo, _ := PackageForImportPath("github.com/fake/uninstalledlibrary")
		context := build.Default
		context.GOPATH = strings.Join([]string{gopathOne, gopathTwo}, ":")
		repo.ctx = &context

		expected := path.Join(gopathOne, "src", "/github.com/fake/uninstalledlibrary")
		actual := repo.Dir()
		if actual != expected {
			t.Errorf("expected Dir to be %v but it was %v", expected, actual)
		}
	})
}

func withDummyBuildContextMultipleGopath(t *testing.T, testFunc func(string, string)) {
	// setup
	fakeGoPath, err := ioutil.TempDir("", "gopath")
	if err != nil {
		t.Fatal(err)
	}
	fakeGoPathTwo, err := ioutil.TempDir("", "gopathtwo")
	if err != nil {
		t.Fatal(err)
	}
	err = os.MkdirAll(path.Join(fakeGoPathTwo, "src/github.com/fake/library"), 0777)
	if err != nil {
		t.Fatal(err)
	}

	// run test
	testFunc(fakeGoPath, fakeGoPathTwo)

	// cleanup
	os.RemoveAll(fakeGoPath)
}
