package main

import (
	"github.com/johnietre/term-games/common"
	minesweeper "github.com/johnietre/term-games/minesweeper/cli"
	wordle "github.com/johnietre/term-games/wordle/cli"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:                   "term-games",
		Short:                 "A suite of terminal games",
		DisableFlagsInUseLine: true,
	}
	rootCmd.AddCommand(minesweeper.MakeCmd())
	rootCmd.AddCommand(wordle.MakeCmd())

	if err := rootCmd.Execute(); err != nil {
		common.Fatal("error running: ", err)
	}
}
