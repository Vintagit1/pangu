package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/urfave/cli"

	"github.com/vinta/pangu"
)

const (
	NAME    = "pangu-axe"
	USAGE   = "Paranoid text spacing"
	VERSION = "2.5.6"
	AUTHOR  = "Vinta Chen"
	EMAIL   = "vinta.chen@gmail.com"
)

// PREFIX is prefix of outpu filename
// TODO
var PREFIX = "readable."

func prefixFilename(path, specified string) string {
	if len(specified) > 0 {
		return specified
	}

	filename := filepath.Base(path)
	newFilename := PREFIX + filename

	return newFilename
}

func processFile(errc chan error, filename, o string) {
	var fw *os.File
	var err error

	_, err = os.Stat(filename)
	if err != nil {
		errc <- err
		return
	}

	switch o {
	case "stdout", "STDOUT":
		fw = os.Stdout
	case "stderr", "STDERR":
		fw = os.Stderr
	default:
		newFilename := prefixFilename(filename, o)
		fw, err = os.Create(newFilename)
		if err != nil {
			errc <- err
			return
		}
		defer fw.Close()
	}

	err = pangu.SpacingFile(filename, fw)
	errc <- err
}

func main() {
	app := cli.NewApp()
	app.Name = NAME
	app.Usage = USAGE
	app.Version = VERSION
	app.Author = AUTHOR
	app.Email = EMAIL
	app.Commands = []cli.Command{
		{
			Name:    "text",
			Usage:   "Performs paranoid text spacing on text",
			Aliases: []string{"t"},
			Action: func(c *cli.Context) {
				if len(c.Args()) == 0 {
					cli.ShowSubcommandHelp(c)
					return
				}

				text := c.Args().First()
				fmt.Println(pangu.SpacingText(text))
			},
		},
		{
			Name:    "file",
			Usage:   "Performs paranoid text spacing on files",
			Aliases: []string{"f"},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "output, o",
					Value: "",
					Usage: fmt.Sprintf(`Specifies the output file name. If not specified, the output file name will be "%sfilename.ext"`, PREFIX),
				},
			},
			Action: func(c *cli.Context) {
				if len(c.Args()) == 0 {
					cli.ShowSubcommandHelp(c)
					return
				}

				o := c.String("output")

				if len(c.Args()) > 1 && len(o) > 0 {
					color.Red(`can't use the "-output" flag with multiple files`)
					os.Exit(1)
				}

				errc := make(chan error)

				for _, filename := range c.Args() {
					go processFile(errc, filename, o)
				}

				for _ = range c.Args() {
					err := <-errc
					if err != nil {
						color.Red("%s", err)
					}
				}
			},
		},
	}

	app.Run(os.Args)
}
