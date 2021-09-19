package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

var version string = "v0.3.3"

func main() {
	app := cli.NewApp()
	app.Usage = "convert go test output to junit xml!"
	app.EnableBashCompletion = true
	app.Version = version
	app.Commands = []*cli.Command{
		{
			Name:   "parse",
			Usage:  "parse input from `FILE` (defaults to stdin if not set)",
			Action: actionParse,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "input",
					Aliases: []string{"i"},
					Usage:   "parse input from `FILE` (defaults to stdin if not set)",
				},
				&cli.StringFlag{
					Name:    "output",
					Aliases: []string{"o"},
					Usage:   "write output to `FILE` (defaults to stdout if not set)",
				},
				&cli.BoolFlag{
					Name:  "fail",
					Usage: "return with a non-zero exit status in the case a parsed test failed",
				},
			},
		},
		{
			Name:  "test",
			Usage: "execute go test and output result as junit xml report",
			UsageText: "go2junit test -- [arguments...]\n\n" +
				"   the test subcommand executes 'go test' with the given arguments and directly converts its output to junit xml",
			Action: actionTest,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "output",
					Aliases: []string{"o"},
					Usage:   "write output to `FILE` (defaults to stdout if not set)",
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
