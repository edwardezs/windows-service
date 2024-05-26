package cli

import "github.com/urfave/cli"

func New() *cli.App {
	return &cli.App{
		Name:     "Example Windows service",
		Flags:    []cli.Flag{CfgFlag},
		Before:   SetupService,
		Commands: []cli.Command{ServiceCmd},
	}
}
