package main

import (
	"bytes"
	"fmt"
	"os/exec"
)

func convert() {
	cmd := exec.Command("/usr/bin/git",
		"log", "-n3", "--name-only", "--grep=Sponsored")
	cmd.Dir = "../../freebsd/freebsd-src/"
	//cmd.Dir = "/home/cnst/github.com/freebsd/freebsd-src/"

	output, err := cmd.Output()

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%s\n", output)

	lines := bytes.Split(output, []byte("\n"))
	fmt.Printf("\noutput byte length: %v; number of lines: %v\n\n", len(output), len(lines))

	for i, s := range lines {
		//fmt.Printf("%v:%s\n", i, s)
		if bytes.HasPrefix(s, []byte("commit")) {
			fmt.Printf("new commit at line %v: %s\n", i, s)
		}
	}

	return
}
