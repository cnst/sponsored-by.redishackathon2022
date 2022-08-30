package main

import (
	"bytes"
	"fmt"
	"os/exec"

	//"github.com/go-redis/redis"
	"github.com/mediocregopher/radix.v2/redis"
)

var rClient, _ = redis.Dial("tcp", "localhost:6379")

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
			//fmt.Printf("\n old:%s\n%s\n{%s}\n, new:%s\n", comm, spon, path, nc[1])
			handleEntry(comm, spon, path)
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
			//for _, ss := range path {
			//	fmt.Printf("%s\n", ss)
			//}
			//fmt.Println("")
		}
	}
	handleEntry(comm, spon, path)

	return
}

func handleEntry(comm []byte, spon []byte, path [][]byte) {

	fmt.Printf("\n\n\n new entry: %s\n%s\n[%v]{%s}\n\n", comm, spon, len(path), path)
	s := bytes.Split(spon, []byte(":"))
	if len(s) >= 2 {
		spon = s[1]
	}
	spon = bytes.TrimSpace(spon)
	fmt.Printf("s: {%s}\n", spon)

	resp := rClient.Cmd("HSET",
		fmt.Sprintf("commit:%s", comm),
		"sponsored", fmt.Sprintf("%s", spon),
		"path", fmt.Sprintf("%s", path),
	)
	if resp != nil {
		fmt.Println(resp)
	}
}
