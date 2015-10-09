package main

import (
	"bufio"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"log"
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
	l := log.New(os.Stderr, "", 0)
	app.Name = "resourceful"
	app.Usage = "Add strong typing to imageNamed: in Swift apps"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "input, i",
			Value:  "Images.xcassets",
			Usage:  "The xcassets directory that contains your app's images",
			EnvVar: "SCRIPT_INPUT_FILE_0",
		},
		cli.StringFlag{
			Name:   "output, o",
			Value:  "Resourceful.swift",
			Usage:  "The destination file for generated Swift code",
			EnvVar: "SCRIPT_OUTPUT_FILE_0",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:  "warn",
			Usage: "Generate xcode warnings for usage of imageNamed",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "warndirectory, w",
					Usage:  "The directory to search for legacy usage of imageNamed",
					EnvVar: "SRCROOT",
				},
			},
			Action: func(c *cli.Context) {
				err := warn(c.String("warndirectory"))
				if err != nil {
					l.Println("Error detecting warnings.")
					os.Exit(1)
				}
			},
		},
	}
	app.Action = func(c *cli.Context) {
		enumType, err := parseFolder(c.String("input"))
		if err != nil {
			l.Println("Error parsing input file.")
			os.Exit(1)
		}
		enumType.Name = "Image"
		f, err := os.Create(c.String("output"))
		if err != nil {
			l.Println("Error writing output file")
			os.Exit(1)
		}
		defer f.Close()
		writer := bufio.NewWriter(f)
		templ := template.Must(template.New("swiftTemplate").Parse(swiftTemplate))
		err = templ.Execute(writer, enumType)
		if err != nil {
			panic(err)
		}
		writer.Flush()
	}

	app.Run(os.Args)
}
