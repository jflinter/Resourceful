package main

import (
	"fmt"
	"os"
	"os/exec"
)

func warn(directory string) error {
	if directory == "" {
		return fmt.Errorf("warn directory should not be nil")
	}
	find := exec.Command(
		"/usr/bin/find",
		directory,
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
		"(UIImage|NSImage)\\(named",
	)
	sed := exec.Command(
		"sed",
		"-E",
		"s/(UIImage|NSImage).*/ warning: legacy use of imageNamed; consider using Resourceful/",
	)

	stdin, err := find.StdoutPipe()
	match.Stdin = stdin
	if err != nil {
		return err
	}
	sed.Stdin, err = match.StdoutPipe()
	if err != nil {
		return err
	}
	sed.Stdout = os.Stderr

	err = match.Start()
	if err != nil {
		return err
	}
	err = sed.Start()
	if err != nil {
		return err
	}
	err = find.Run()
	if err != nil {
		return err
	}
	err = match.Wait()
	if err != nil {
		return err
	}
	err = sed.Wait()
	if err != nil {
		return err
	}
	return nil
}
