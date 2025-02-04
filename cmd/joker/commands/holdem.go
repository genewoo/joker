package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/genewoo/joker/internal/deck"
	"github.com/genewoo/joker/internal/holdem"
	"github.com/spf13/cobra"
)

// NewHoldemCmd creates a new holdem game command
func NewHoldemCmd(options *HoldemOptions) *cobra.Command {
	// Store game type as string for flag
	var gameTypeStr string

	holdemCmd := &cobra.Command{
		Use:   "holdem",
		Short: "Texas Hold'em style game commands",
		Long:  `Texas Hold'em style poker game with various options.`,
	}

	dealCmd := createDealCmd(options, &gameTypeStr)
	eqCmd := createEquityCmd(options)

	holdemCmd.AddCommand(dealCmd, eqCmd)
	return holdemCmd
}

func createDealCmd(options *HoldemOptions, gameTypeStr *string) *cobra.Command {

	// Create help message for game types using AllGameTypes
	gameTypes := make([]string, len(holdem.AllGameTypes()))
	for i, gt := range holdem.AllGameTypes() {
		gameTypes[i] = gt.String()
	}
	gameTypeHelp := fmt.Sprintf("Game type (%s)", strings.Join(gameTypes, ", "))

	dealCmd := &cobra.Command{
		Use:   "deal",
		Short: "Deal cards in Hold'em style",
		Run: func(cmd *cobra.Command, args []string) {
			// Convert string game type to GameType enum
			gameType, err := holdem.ParseGameType(*gameTypeStr)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			options.GameType = gameType

			// Create new holdem game
			game := holdem.NewGame(options.GameType, options.NumPlayers)

			// Deal initial cards
			if err := game.StartHand(); err != nil {
				fmt.Printf("Error dealing cards: %v\n", err)
				os.Exit(1)
			}

			// Print player hands
			fmt.Printf("Dealing %d cards to %d players:\n", options.NumCardsPerPlayer, options.NumPlayers)
			for i, player := range game.Players {
				fmt.Printf("\nPlayer %d:\n", i+1)
				for _, card := range player.Cards {
					fmt.Printf("%s ", card.String())
				}
				fmt.Println()
			}

			// Deal and show flop
			if err := game.DealFlop(); err != nil {
				fmt.Printf("Error dealing flop: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("\nFlop:\n")
			for _, card := range game.Community[:3] {
				fmt.Printf("%s ", card.String())
			}
			fmt.Println()

			// Deal and show turn
			if err := game.DealTurnOrRiver(); err != nil {
				fmt.Printf("Error dealing turn: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("\nTurn:\n%s\n", game.Community[3].String())

			// Deal and show river
			if err := game.DealTurnOrRiver(); err != nil {
				fmt.Printf("Error dealing river: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("\nRiver:\n%s\n", game.Community[4].String())
		},
	}

	dealCmd.Flags().IntVarP(&options.NumPlayers, "players", "p", 2, "Number of players")
	dealCmd.Flags().IntVarP(&options.NumCardsPerPlayer, "numberofcards", "n", 2, "Number of cards per player")
	dealCmd.Flags().StringVarP(gameTypeStr, "type", "t", holdem.Texas.String(), gameTypeHelp)

	return dealCmd
}

func createEquityCmd(options *HoldemOptions) *cobra.Command {
	eqCmd := &cobra.Command{
		Use:   "eq",
		Short: "Calculate equity for players",
		Long: `Calculate equity (winning probability) for each player in a Texas Hold'em game.
Example card format: "As Kh" for Ace of spades and King of hearts.
Use "♠" for spades, "♥" for hearts, "♦" for diamonds, "♣" for clubs.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Parse player cards
			if len(options.PlayerCards) == 0 {
				fmt.Println("Error: At least one player's cards must be specified")
				os.Exit(1)
			}

			// Convert player cards strings to Card objects
			players := make([][]*deck.Card, len(options.PlayerCards))
			for i, cardStr := range options.PlayerCards {
				cards := strings.Fields(cardStr)
				if len(cards) != 2 {
					fmt.Printf("Error: Player %d must have exactly 2 cards, got: %s\n", i+1, cardStr)
					os.Exit(1)
				}

				players[i] = make([]*deck.Card, 2)
				for j, card := range cards {
					value := card[:len(card)-1]
					suit := card[len(card)-1:]
					players[i][j] = deck.NewCard(value, suit)
				}
			}

			// Parse community cards if provided
			var community []*deck.Card
			if options.CommunityCards != "" {
				cards := strings.Fields(options.CommunityCards)
				if len(cards) > 5 {
					fmt.Println("Error: Maximum 5 community cards allowed")
					os.Exit(1)
				}

				for _, card := range cards {
					value := card[:len(card)-1]
					suit := card[len(card)-1:]
					community = append(community, deck.NewCard(value, suit))
				}
			}

			// Create calculator and calculate probabilities
			calc := holdem.NewWinningCalculator(players, options.NumSimulations, holdem.NewDefaultHandRanker(), community...)
			probabilities := calc.CalculateWinProbabilities()

			// Display results
			fmt.Println("\nEquity calculation results:")
			for i := 0; i < len(players); i++ {
				fmt.Printf("Player %d (%s): %.2f%%\n", i+1, options.PlayerCards[i], probabilities[i]*100)
			}
			if len(community) > 0 {
				fmt.Printf("\nCommunity cards: %s\n", options.CommunityCards)
			}
			fmt.Printf("Tie probability: %.2f%%\n", probabilities[len(players)]*100)
		},
	}

	eqCmd.Flags().StringSliceVarP(&options.PlayerCards, "cards", "c", []string{}, "Player hole cards (e.g. \"As Kh\" \"Jd Tc\")")
	eqCmd.Flags().StringVarP(&options.CommunityCards, "board", "b", "", "Community cards (e.g. \"Ah Kd Qc\")")
	eqCmd.Flags().IntVarP(&options.NumSimulations, "simulations", "s", 10000, "Number of Monte Carlo simulations")
	eqCmd.MarkFlagRequired("cards")

	return eqCmd
}
