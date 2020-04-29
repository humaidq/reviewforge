package main

import (
	"log"
	"os"

	"git.sr.ht/~humaid/reviewforge/cmd"

	"github.com/urfave/cli/v2"
)

// VERSION specifies the version of reviewforge
var VERSION = "0.1.0"

func main() {
	app := cli.NewApp()
	app.Name = "reviewforge"
	app.Usage = "The unified code review platform"
	app.Version = VERSION
	app.Commands = []*cli.Command{
		cmd.CmdStart,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
