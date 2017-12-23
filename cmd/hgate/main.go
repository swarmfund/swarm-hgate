package main

import (
	"os"

	"github.com/urfave/cli"
)

func main() {
	cmd := cli.NewApp()

	cmd.Commands = []cli.Command{
		ServeCommand,
	}

	cmd.Run(os.Args)
}
