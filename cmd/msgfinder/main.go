package main

import (
	"os"

	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

var options Options

func main() {
	// Defaults
	options.LogLevel = "INFO"
	options.LogEncoding = "console"

	cliFlags := []cli.Flag{}

	app := &cli.App{
		Name:    "storenode-messages",
		Version: "0.0.1",
		Before:  altsrc.InitInputSourceWithContext(cliFlags, altsrc.NewTomlSourceFromFlagFunc("config-file")),
		Flags:   cliFlags,
		Action: func(c *cli.Context) error {
			Execute(c.Context, options)
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
