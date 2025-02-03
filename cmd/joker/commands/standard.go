package commands

import (
	"fmt"
	"os"

	"github.com/genewoo/joker/internal/deck"
	"github.com/spf13/cobra"
)

// NewStandardCmd creates a new standard game command
func NewStandardCmd(options *StandardOptions) *cobra.Command {
	standardCmd := &cobra.Command{
		Use:   "standard",
		Short: "Standard card game commands",
		Long:  `Standard card game with various options for dealing cards.`,
	}

	dealCmd := &cobra.Command{
		Use:   "deal",
		Short: "Deal cards to players",
		Run: func(cmd *cobra.Command, args []string) {
			// Create a new deck based on options
			var d *deck.Deck
			if options.IncludeJokers {
				d = deck.NewDeckWithJokers()
			} else {
				d = deck.NewDeck()
			}

			// Multiply deck if needed
			if options.NumDecks > 1 {
				d = d.Times(options.NumDecks)
			}

			// Calculate cards per player if not specified
			if options.NumCardsPerPlayer == 0 {
				totalCards := d.Count() - options.KeepCards
				options.NumCardsPerPlayer = totalCards / options.NumPlayers
			}

			// Validate parameters
			if options.NumCardsPerPlayer*options.NumPlayers+options.KeepCards > d.Count() {
				fmt.Printf("Error: Not enough cards in deck. Have %d cards, need %d cards (%d players × %d cards + %d kept cards)\n",
					d.Count(), options.NumCardsPerPlayer*options.NumPlayers+options.KeepCards,
					options.NumPlayers, options.NumCardsPerPlayer, options.KeepCards)
				os.Exit(1)
			}

			// Shuffle the deck
			d.Shuffle()

			// Deal cards to each player
			fmt.Printf("Dealing %d cards to %d players (keeping %d cards):\n",
				options.NumCardsPerPlayer, options.NumPlayers, options.KeepCards)
			cards := d.Cards
			currentCard := 0

			// First deal the cards that will be kept aside
			if options.KeepCards > 0 {
				fmt.Println("\nKept cards:")
				for i := 0; i < options.KeepCards; i++ {
					fmt.Printf("%s ", cards[currentCard].String())
					currentCard++
				}
				fmt.Println()
			}

			// Then deal to players
			for player := 1; player <= options.NumPlayers; player++ {
				fmt.Printf("\nPlayer %d:\n", player)
				for i := 0; i < options.NumCardsPerPlayer; i++ {
					fmt.Printf("%s ", cards[currentCard].String())
					currentCard++
				}
				fmt.Println()
			}
		},
	}

	// Add flags
	dealCmd.Flags().IntVarP(&options.NumPlayers, "players", "p", 2, "Number of players")
	dealCmd.Flags().IntVarP(&options.NumDecks, "decks", "d", 1, "Number of decks")
	dealCmd.Flags().BoolVarP(&options.IncludeJokers, "joker", "j", true, "Include jokers")
	dealCmd.Flags().IntVarP(&options.KeepCards, "keep", "k", 0, "Number of cards to keep")
	dealCmd.Flags().IntVarP(&options.NumCardsPerPlayer, "numberofcards", "n", 0, "Number of cards per player")

	standardCmd.AddCommand(dealCmd)
	return standardCmd
}
