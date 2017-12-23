package main

import (
	"os"

	"github.com/urfave/cli"
)

func main() {
	cmd := cli.NewApp()
	cmd.Usage = "A simple proxy for interacting with the Horizon server"
	cmd.Version = "0.1.0"
	cmd.Commands = []cli.Command{
		ServeCommand,
	}

	cmd.Run(os.Args)
}
