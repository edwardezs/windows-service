package cli

import "github.com/urfave/cli"

func New(svcName string) *cli.App {
	return &cli.App{
		Name:     svcName,
		Flags:    []cli.Flag{CfgFlag},
		Before:   SetupService,
		Commands: ServiceCmd,
	}
}
