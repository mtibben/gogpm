package main

import (
	"bytes"
	"log"
	"os/exec"
)

func execCmd(cmdname string, cmdargs ...string) (string, error) {
	var out bytes.Buffer

	command := exec.Command(cmdname, cmdargs...)
	command.Stdout = &out
	command.Stderr = &out

	err := command.Run()
	if err != nil {
		log.Printf("Error while executing: %s %v\n", cmdname, cmdargs)
		log.Println(out.String())
		return out.String(), err
	}

	return out.String(), nil
}
