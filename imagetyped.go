package main

import (
	"bufio"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"os"
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
	app.Name = "imagetyped"
	app.Usage = "Add strong typing to imageNamed: in Swift apps"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "input, i",
			Value: "Images.xcassets",
			Usage: "The xcassets directory that contains your app's images",
		},
		cli.StringFlag{
			Name:  "output, o",
			Value: "ImageTyped.swift",
			Usage: "The destination file for generated Swift code",
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
		templ, err := template.New("template.swift").ParseFiles("template.swift")
		err = templ.Execute(writer, templateArgs)
		if err != nil {
			panic(err)
		}
		writer.Flush()
	}

	app.Run(os.Args)
}
