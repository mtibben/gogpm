package vcs

import (
	"fmt"
	"go/build"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPresentPackageRepoSimpleGopath(t *testing.T) {
	withDummyBuildContextSingleGopath(t, func(gopath string) {
		repo, _ := PackageForImportPath("github.com/fake/library")
		repo.ctx = build.Context{Compiler: "gc", GOPATH: gopath}

		expected := filepath.Join(gopath, "src", "github.com", "fake", "library")
		actual := repo.Dir()
		if actual != expected {
			t.Errorf("expected Dir to be %v but it was %v", expected, actual)
		}
	})
}

func TestNotYetInstalledPackageRepoSimpleGopath(t *testing.T) {
	withDummyBuildContextSingleGopath(t, func(gopath string) {
		repo, _ := PackageForImportPath("github.com/fake/uninstalledlibrary")
		repo.ctx = build.Context{Compiler: "gc", GOPATH: gopath}

		expected := filepath.Join(gopath, "src", "github.com", "fake", "uninstalledlibrary")
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
	if err = os.MkdirAll(filepath.Join(fakeGoPath, "src", "github.com", "fake", "library"), 0777); err != nil {
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
		gopath := strings.Join([]string{gopathOne, gopathTwo}, fmt.Sprintf("%c", os.PathListSeparator))
		repo.ctx = build.Context{Compiler: "gc", GOPATH: gopath}

		expected := filepath.Join(gopathTwo, "src", "github.com", "fake", "library")
		actual := repo.Dir()
		if actual != expected {
			t.Errorf("expected Dir to be %v but it was %v", expected, actual)
		}
	})
}

func TestMultipleGopathNoInstall(t *testing.T) {
	withDummyBuildContextMultipleGopath(t, func(gopathOne string, gopathTwo string) {
		repo, _ := PackageForImportPath("github.com/fake/uninstalledlibrary")
		gopath := strings.Join([]string{gopathOne, gopathTwo}, fmt.Sprintf("%c", os.PathListSeparator))
		repo.ctx = build.Context{Compiler: "gc", GOPATH: gopath}

		expected := filepath.Join(gopathOne, "src", "github.com", "fake", "uninstalledlibrary")
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
	if err = os.MkdirAll(filepath.Join(fakeGoPathTwo, "src", "github.com", "fake", "library"), 0777); err != nil {
		t.Fatal(err)
	}

	// run test
	testFunc(fakeGoPath, fakeGoPathTwo)

	// cleanup
	os.RemoveAll(fakeGoPath)
}
