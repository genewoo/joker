package main

import (
	"fmt"
	"os"

	"github.com/genewoo/joker/cmd/joker/commands"
	"github.com/genewoo/joker/internal/holdem"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "joker",
		Short: "Joker is a poker calculation tool",
		Long: `A git-like command line tool for poker players to perform various calculations
and card dealings for different types of poker games.`,
	}

	// Create options for commands
	standardOpts := &commands.StandardOptions{
		CommonOptions: commands.CommonOptions{
			NumPlayers:        2,
			NumCardsPerPlayer: 0,
		},
		NumDecks:      1,
		IncludeJokers: true,
		KeepCards:     0,
	}

	holdemOpts := &commands.HoldemOptions{
		CommonOptions: commands.CommonOptions{
			NumPlayers:        2,
			NumCardsPerPlayer: 2,
		},
		GameType:       holdem.Texas,
		NumSimulations: 10000,
	}

	// Add commands
	rootCmd.AddCommand(
		commands.NewStandardCmd(standardOpts),
		commands.NewHoldemCmd(holdemOpts),
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
