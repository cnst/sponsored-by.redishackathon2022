package main

import (
	"bytes"
	"fmt"
	"os/exec"
)

func convert() {
	cmd := exec.Command("/usr/bin/git",
		"log", "-n5", "--name-only", "--grep=Sponsored")
	cmd.Dir = "../../freebsd/freebsd-src/"
	//cmd.Dir = "/home/cnst/github.com/freebsd/freebsd-src/"

	output, err := cmd.Output()

	if err != nil {
		fmt.Println(err)
		return
	}
	//fmt.Printf("%s\n", output)

	lines := bytes.Split(output, []byte("\n"))
	fmt.Printf("\noutput byte length: %v; number of lines: %v\n\n", len(output), len(lines))

	var comm []byte
	var spon []byte
	var path [][]byte
	
	for i, s := range lines {
		//fmt.Printf("%v:%s\n", i, s)
		if bytes.HasPrefix(s, []byte("commit ")) {
			fmt.Printf("new commit at line %v: %s\n", i, s)
			nc := bytes.Split(s, []byte(" "))
			fmt.Printf("\n old:%s\n%s\n{%s}\n, new:%s\n", comm, spon, path, nc[1])
			comm = nc[1]
			spon = nil
			path = nil
		} else if bytes.HasPrefix(s, []byte("Author")) {
			//auth = s
		} else if bytes.HasPrefix(s, []byte("Date")) {
			//date = s
		} else if bytes.HasPrefix(s, []byte("    Sponsored")) {
			spon = s
		} else if bytes.HasPrefix(s, []byte("    ")) {
			continue
		} else if len(s) > 0 {
			path = append(path, s)
			for _, ss := range path {
				fmt.Printf("%s\n", ss)
			}
			fmt.Println("")
		}
	}

	return
}
