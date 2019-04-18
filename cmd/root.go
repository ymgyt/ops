package cmd

import (
	"github.com/ymgyt/cli"
)

func New() *cli.Command {
	root := &cli.Command{
		Name:     "ops",
		LongDesc: "ops is utility tool for daily operations.",
	}

	return root.
		AddCommand(NewDiskUsageCommand()).
		AddCommand(NewVersionCommand()).
		AddCommand(NewRandCommand())
}
