package main

import (
	"bufio"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"text/template"
)

type EnumCase struct {
	Name     string
	RawValue string
}

type EnumType struct {
	Name   string
	Values []EnumCase
	Nested []EnumType
}

func parseFolder(filepath string) (EnumType, error) {
	files, err := ioutil.ReadDir(filepath)
	if err != nil {
		return EnumType{}, err
	}
	enumCases := []EnumCase{}
	enumTypes := []EnumType{}
	for _, f := range files {
		if strings.HasSuffix(f.Name(), "imageset") {
			rawValue := strings.Split(f.Name(), ".")[0]
			name := strings.Title(rawValue)
			name = strings.Replace(name, "_", "__", -1)
			name = strings.Replace(name, "-", "_", -1)
			enum := EnumCase{name, rawValue}
			enumCases = append(enumCases, enum)
		} else if !strings.Contains(f.Name(), ".") {
			folderName := path.Join(filepath, f.Name())
			enumType, err := parseFolder(folderName)
			if err != nil {
				return EnumType{}, err
			}
			enumTypes = append(enumTypes, enumType)
		}
	}
	_, file := path.Split(filepath)
	name := strings.Title(file)
	return EnumType{name, enumCases, enumTypes}, nil
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
		enumType, err := parseFolder(c.String("input"))
		if err != nil {
			panic(err)
		}
		enumType.Name = "Image"
		f, err := os.Create(c.String("output"))
		defer f.Close()
		writer := bufio.NewWriter(f)
		templ, err := template.New("swiftTemplate").Parse(swiftTemplate)
		err = templ.Execute(writer, enumType)
		if err != nil {
			panic(err)
		}
		writer.Flush()

		if c.Bool("warn") {
			warn()
		}

	}

	app.Run(os.Args)
}
