package cmd

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/ymgyt/cli"
)

const (
	randDesc            = "rand generate random characters."
	defaultRandomLength = 40
	randCharacters      = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"
)

func NewRandCommand() *cli.Command {
	rd := &RandomGenerator{}

	cmd := &cli.Command{
		Name:      "rand",
		Aliases:   []string{"random"},
		ShortDesc: randDesc,
		LongDesc: randDesc + "\n\nUSAGE\n ops rand [OPTIONS]\n\n" +
			"default character pool: " + randCharacters + "\n",
		Run: rd.Run,
	}

	err := cmd.Options().
		Add(&cli.BoolOpt{Var: &rd.showHelp, Long: "help", Short: "h", Description: "print this."}).
		Add(&cli.BoolOpt{Var: &rd.noTrailingNewLine, Long: "without-trailing-newline", Description: "output without adding trailing newline code"}).
		Add(&cli.IntOpt{
			Var: &rd.length, Long: "len", Aliases: []string{"max", "length"}, Default: defaultRandomLength,
			Description: "random strings length."}).
		Add(&cli.StringOpt{
			Var: &rd.specialChars, Long: "special-chars", Aliases: []string{"special-char"},
			Description: "characters used in addition to the standard string pool for random string generation",
		}).Err
	if err != nil {
		panic(err)
	}

	return cmd
}

type RandomGenerator struct {
	showHelp          bool
	length            int
	specialChars      string
	noTrailingNewLine bool
}

func (rd *RandomGenerator) Run(ctx context.Context, cmd *cli.Command, args []string) {
	if rd.showHelp {
		cli.HelpFunc(os.Stdout, cmd)
		return
	}

	pool := randCharacters
	pool += rd.specialChars

	generated := rd.generate(pool, rd.length, time.Now().UnixNano())

	if rd.noTrailingNewLine {
		fmt.Print(generated)
	} else {
		fmt.Println(generated)
	}
}

func (rd *RandomGenerator) generate(pool string, length int, seed int64) string {
	r := rand.New(rand.NewSource(seed))

	var sb strings.Builder
	for i := 0; i < length; i++ {
		sb.WriteByte(pool[r.Intn(len(pool))])
	}
	return sb.String()
}
