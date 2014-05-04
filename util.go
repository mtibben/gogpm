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

func execCmd(cmd string) (string, string, error) {
	command := exec.Command("bash")
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	command.Stdin = bytes.NewBufferString(cmd)
	command.Stdout = &stdout
	command.Stderr = &stderr
	err := command.Run()
	if err != nil {
		return "", "", err
	}

	return stdout.String(), stderr.String(), nil
}
