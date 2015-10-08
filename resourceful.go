package main

import (
	"bufio"
	"fmt"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

type SwiftEnum struct {
	Name     string
	RawValue string
}

type TemplateArgs struct {
	EnumValues []SwiftEnum
}

func main() {
	app := cli.NewApp()
	app.Name = "resourceful"
	app.Usage = "Add strong typing to imageNamed: in Swift apps"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "input, i",
			Value: "Images.xcassets",
			Usage: "The xcassets directory that contains your app's images",
		},
		cli.StringFlag{
			Name:  "output, o",
			Value: "Resourceful.swift",
			Usage: "The destination file for generated Swift code",
		},
		cli.BoolFlag{
			Name:  "warn, w",
			Usage: "Generate xcode warnings for usage of imageNamed",
		},
	}
	app.Action = func(c *cli.Context) {
		files, err := ioutil.ReadDir(c.String("input"))
		if err != nil {
			panic(err)
		}
		enumValues := []SwiftEnum{}
		for _, f := range files {
			if strings.HasSuffix(f.Name(), "imageset") {
				rawValue := strings.Split(f.Name(), ".")[0]
				name := strings.Title(rawValue)
				name = strings.Replace(name, "_", "__", -1)
				name = strings.Replace(name, "-", "_", -1)
				enum := SwiftEnum{name, rawValue}
				enumValues = append(enumValues, enum)
			}
		}
		f, err := os.Create(c.String("output"))
		defer f.Close()
		writer := bufio.NewWriter(f)
		templateArgs := TemplateArgs{enumValues}
		templ, err := template.New("swiftTemplate").Parse(swiftTemplate)
		err = templ.Execute(writer, templateArgs)
		if err != nil {
			panic(err)
		}
		writer.Flush()

		if c.Bool("warn") {
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

			match.Stdin, err = find.StdoutPipe()
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

	}

	app.Run(os.Args)
}
