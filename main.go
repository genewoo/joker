package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/genewoo/joker/internal/deck"
	"github.com/genewoo/joker/internal/holdem"
)

func main() {
	// Parse command line flags
	numPlayers := flag.Int("n", 2, "Number of players")
	flag.Parse()

	if *numPlayers < 2 {
		log.Fatal("Number of players must be at least 2")
	}

	// Create new game
	game := holdem.NewGame(*numPlayers)

	// Start the hand - deal cards to players
	if err := game.StartHand(); err != nil {
		log.Fatalf("Failed to start hand: %v", err)
	}

	// Print each player's cards
	fmt.Println("Initial hands:")
	playerCards := make([][]*deck.Card, *numPlayers)
	for i, player := range game.Players {
		fmt.Printf("Player %d: %v %v\n", i+1, player.Cards[0], player.Cards[1])
		playerCards[i] = player.Cards
	}

	ranker := holdem.NewSmartHandRanker()
	// Calculate initial winning probabilities
	calculator := holdem.NewWinningCalculator(playerCards, 10000, ranker)
	probabilities := calculator.CalculateWinProbabilities()

	fmt.Println("\nInitial winning probabilities:")
	for i := 0; i < len(probabilities)-1; i++ {
		fmt.Printf("Player %d: %.2f%%\n", i+1, probabilities[i]*100)
	}
	fmt.Printf("Tie probability: %.2f%%\n", probabilities[len(probabilities)-1]*100)

	// Deal the flop
	fmt.Println("\nDealing the flop...")
	if err := game.DealFlop(); err != nil {
		log.Fatalf("Failed to deal flop: %v", err)
	}
	fmt.Printf("Flop: %v %v %v\n", game.Community[0], game.Community[1], game.Community[2])

	// Deal the turn
	fmt.Println("\nDealing the turn...")
	if err := game.DealTurnOrRiver(); err != nil {
		log.Fatalf("Failed to deal turn: %v", err)
	}
	fmt.Printf("Turn: %v\n", game.Community[3])

	// Deal the river
	fmt.Println("\nDealing the river...")
	if err := game.DealTurnOrRiver(); err != nil {
		log.Fatalf("Failed to deal river: %v", err)
	}
	fmt.Printf("River: %v\n", game.Community[4])

	// Calculate final hands and determine winner
	fmt.Println("\nFinal hands:")

	bestHands := make([]holdem.HandStrength, *numPlayers)

	for i, player := range game.Players {
		bestHand, _ := ranker.RankHand(player.Cards, game.Community)
		bestHands[i] = bestHand
		fmt.Printf("Player %d: %s %s - %v\n", i+1, player.Cards[0], player.Cards[1], bestHand)
	}

	// Find winners
	winners := holdem.FindWinners(bestHands)

	fmt.Println("\nWinners:")
	if len(winners) == 1 {
		fmt.Printf("Player %d wins!\n", winners[0]+1)
	} else {
		fmt.Printf("Tie between players: ")
		for i, winner := range winners {
			if i > 0 {
				fmt.Print(", ")
			}
			fmt.Printf("%d", winner+1)
		}
		fmt.Println()
	}
}
