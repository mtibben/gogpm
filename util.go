package main

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
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
	var cmdargs []string
	var out bytes.Buffer

	cmdfields := strings.Fields(cmd)
	if len(cmdfields) > 1 {
		cmdargs = cmdfields[1:]
	}

	command := exec.Command(cmdfields[0], cmdargs...)
	command.Stdout = &out
	command.Stderr = &out

	err := command.Run()
	if err != nil {
		return out.String(), err
	}

	return out.String(), nil
}
