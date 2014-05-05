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

func execCmd(cmd string) (string, error) {
	command := exec.Command("bash")
	var out bytes.Buffer
	command.Stdin = bytes.NewBufferString(cmd)
	command.Stdout = &out
	command.Stderr = &out
	err := command.Run()
	if err != nil {
		return out.String(), err
	}

	return out.String(), nil
}
