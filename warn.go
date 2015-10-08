package main

import (
	"fmt"
	"os"
	"os/exec"
)

func warn() {
	srcroot := os.Getenv("SRCROOT")
	if srcroot == "" {
		panic(fmt.Errorf("SRCROOT should not be nil"))
	}
	find := exec.Command(
		"/usr/bin/find",
		srcroot,
		"-name",
		"*.swift",
		"-print0",
	)
	match := exec.Command(
		"xargs",
		"-0",
		"egrep",
		"--with-filename",
		"--line-number",
		"--only-matching",
		"UIImage\\(named",
	)
	sed := exec.Command(
		"sed",
		"s/UIImage.*/ warning: legacy use of imageNamed; consider using Resourceful/",
	)

	stdin, err := find.StdoutPipe()
	match.Stdin = stdin
	if err != nil {
		panic(err)
	}
	sed.Stdin, err = match.StdoutPipe()
	if err != nil {
		panic(err)
	}
	sed.Stdout = os.Stderr

	err = match.Start()
	if err != nil {
		panic(err)
	}
	err = sed.Start()
	if err != nil {
		panic(err)
	}
	err = find.Run()
	if err != nil {
		panic(err)
	}
	err = match.Wait()
	if err != nil {
		panic(err)
	}
	err = sed.Wait()
	if err != nil {
		panic(err)
	}
}
