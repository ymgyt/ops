package cmd

import (
	"context"
	"fmt"

	"github.com/ymgyt/cli"
)

var version = "v0.0.2"

func NewVersionCommand() *cli.Command {
	cmd := &cli.Command{
		Name:      "version",
		ShortDesc: "print version",
		Run: func(_ context.Context, _ *cli.Command, _ []string) {
			fmt.Println(version)
		},
	}
	return cmd
}
