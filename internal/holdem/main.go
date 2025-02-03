package holdem

import (
	"flag"
	"fmt"
	"log"

	"github.com/genewoo/joker/internal/deck"
)

// printProbabilities prints the winning probabilities for each player
func printProbabilities(probabilities []float64) {
	total := 0.0
	for i := 0; i < len(probabilities); i++ {
		total += probabilities[i]
		if i < len(probabilities)-1 {
			fmt.Printf("Player %d: %.2f%%\n", i+1, probabilities[i]*100)
		} else {
			fmt.Printf("Tie probability: %.2f%%\n", probabilities[i]*100)
		}
	}
	fmt.Printf("Total probability: %.2f%%\n", total*100)
}

func main() {
	// Parse command line flags
	numPlayers := flag.Int("n", 2, "Number of players")
	flag.Parse()

	if *numPlayers < 2 {
		log.Fatal("Number of players must be at least 2")
	}

	// Create new game
	game := NewGame(Texas, *numPlayers)

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

	ranker := NewSmartHandRanker()
	// Calculate initial winning probabilities
	calculator := NewWinningCalculator(playerCards, 10000, ranker)
	fmt.Println("\nInitial winning probabilities:")
	printProbabilities(calculator.CalculateWinProbabilities())

	// Deal the flop
	fmt.Println("\nDealing the flop...")
	if err := game.DealFlop(); err != nil {
		log.Fatalf("Failed to deal flop: %v", err)
	}
	fmt.Printf("Flop: %v %v %v\n", game.Community[0], game.Community[1], game.Community[2])

	// Update calculator with flop cards and show new probabilities
	if err := calculator.AppendCommunityCards(game.Community[0], game.Community[1], game.Community[2]); err != nil {
		log.Fatalf("Failed to append flop cards: %v", err)
	}
	fmt.Println("\nProbabilities after flop:")
	printProbabilities(calculator.CalculateWinProbabilities())

	// Deal the turn
	fmt.Println("\nDealing the turn...")
	if err := game.DealTurnOrRiver(); err != nil {
		log.Fatalf("Failed to deal turn: %v", err)
	}
	fmt.Printf("Turn: %v\n", game.Community[3])

	// Update calculator with turn card and show new probabilities
	if err := calculator.AppendCommunityCards(game.Community[3]); err != nil {
		log.Fatalf("Failed to append turn card: %v", err)
	}
	fmt.Println("\nProbabilities after turn:")
	printProbabilities(calculator.CalculateWinProbabilities())

	// Deal the river
	fmt.Println("\nDealing the river...")
	if err := game.DealTurnOrRiver(); err != nil {
		log.Fatalf("Failed to deal river: %v", err)
	}
	fmt.Printf("River: %v\n", game.Community[4])

	// Update calculator with river card and evaluate final hands
	if err := calculator.AppendCommunityCards(game.Community[4]); err != nil {
		log.Fatalf("Failed to append river card: %v", err)
	}

	// Calculate final hands and determine winner using showdown
	fmt.Println("\nFinal hands:")
	result, err := calculator.EvaluateShowdown()
	if err != nil {
		log.Fatalf("Failed to evaluate showdown: %v", err)
	}

	// Print final hands with their best 5-card combinations
	for i := range result.HandStrengths {
		fmt.Printf("Player %d: %s %s - %v (Best hand: %v)\n",
			i+1,
			game.Players[i].Cards[0],
			game.Players[i].Cards[1],
			result.HandStrengths[i],
			result.BestHands[i])
	}

	fmt.Println("\nWinners:")
	if len(result.Winners) == 1 {
		fmt.Printf("Player %d wins!\n", result.Winners[0]+1)
	} else {
		fmt.Printf("Tie between players: ")
		for i, winner := range result.Winners {
			if i > 0 {
				fmt.Print(", ")
			}
			fmt.Printf("%d", winner+1)
		}
		fmt.Println()
	}
}
