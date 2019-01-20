package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ymgyt/cli"
	"github.com/ymgyt/ops/filesystem"
)

const (
	duDesc = "du print provided directory disk usage recursively."
)

func NewDiskUsageCommand() *cli.Command {
	du := &DiskUsage{}

	cmd := &cli.Command{
		Name:      "du",
		ShortDesc: duDesc,
		LongDesc:  duDesc + "\n\nUSAGE\n  ops du [OPTIONS] <root_directory>",
		Aliases:   []string{"diskusage"},
		Run:       du.Run,
	}

	err := cmd.Options().
		Add(&cli.BoolOpt{Var: &du.showHelp, Long: "help", Short: "h", Description: "print this."}).
		Add(&cli.BoolOpt{Var: &du.humanize, Long: "humanize", Default: true, Description: "humanize bytes. ex MB"}).
		Add(&cli.BoolOpt{Var: &du.ibytes, Long: "ibytes", Short: "i", Description: "humanize ibytes. ex MiB"}).
		Add(&cli.StringOpt{Var: &du.order, Long: "order", Default: "size", Description: "output disk usaage order."}).
		Add(&cli.IntOpt{Var: &du.level, Long: "level", Short: "l", Default: -1, Description: "print directory level."}).
		Add(&cli.IntOpt{Var: &du.maxRecursion, Long: "max-recursion", Aliases: []string{"max"}, Default: 20, Description: "max recursion."}).
		Add(&cli.BoolOpt{Var: &du.relative, Long: "relative", Short: "r", Aliases: []string{"rel"}, Description: "print relative file path from root."}).
		Add(&cli.BoolOpt{Var: &du.verbose, Long: "verbose", Short: "v", Description: "verbose"}).Err
	if err != nil {
		panic(err)
	}

	return cmd
}

type DiskUsage struct {
	showHelp     bool
	level        int
	maxRecursion int
	order        string
	humanize     bool
	ibytes       bool
	relative     bool
	verbose      bool
}

func (du *DiskUsage) Run(ctx context.Context, cmd *cli.Command, args []string) {

	if du.showHelp {
		cli.HelpFunc()(os.Stderr, cmd)
		return
	}

	root, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if len(args) == 1 {
		root = args[0]
	}
	if len(args) > 1 {
		fmt.Fprintln(os.Stderr, "multi root args not supported yet.")
		os.Exit(1)
	}
	if !strings.HasPrefix(root, "/") {
		if abs, err := filepath.Abs(root); err == nil {
			root = abs
		}
	}

	reporter := filesystem.New().DiskUsageReporter()

	res, err := reporter.Do(ctx, &filesystem.DiskUsageRequest{
		Root:         root,
		MaxRecursion: du.maxRecursion,
		Verbose:      du.verbose,
	})
	if err != nil {
		// きちんとerror定義してハンドリングしたい
		if err.Error() == "context canceled" { // signal ctr + c
			fmt.Fprintln(os.Stderr, "stop du")
			os.Exit(3)
		}
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	res.Root.Entries().
		Relative(du.relative).
		Level(du.level).
		Order(du.order).
		Humanize(du.humanize).
		IBytes(du.ibytes).
		Print(os.Stdout)
}
