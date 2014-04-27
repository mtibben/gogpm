package main

import (
	"bytes"
	"os"
	"os/exec"
)

func fileExists(path string) bool {
	fi, err := os.Stat(path)
	return !os.IsNotExist(err) && !fi.IsDir()
}

func dirExists(path string) bool {
	fi, err := os.Stat(path)
	return !os.IsNotExist(err) && fi.IsDir()
}

func execCmd(cmd string) string {
	command := exec.Command("bash")
	var b bytes.Buffer
	command.Stdin = bytes.NewBufferString(cmd)
	command.Stdout = &b
	// command.Stderr = os.Stderr
	err := command.Run()
	if err != nil {
		panic(err)
	}

	return b.String()
}

func mustChdir(path string) {
	err := os.Chdir(path)
	if err != nil {
		panic(err)
	}
}
