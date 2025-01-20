package main

import (
	"fmt"

	"github.com/genewoo/joker/internal/dealer"
	"github.com/genewoo/joker/internal/deck"
)

func main() {
	// Create a deck excluding specific cards
	deck := deck.NewDeck("A♠", "K♥", "7♦")
	deck.Shuffle()

	fmt.Println("Created deck excluding A♠, K♥, and 7♦")

	dealer := &dealer.StandardDealer{}

	// Deal 5 cards
	hand, err := dealer.Deal(deck, 5)
	if err != nil {
		fmt.Println("Error dealing cards:", err)
		return
	}

	fmt.Println("\nDealt cards:")
	for i := 0; i < hand.Count(); i++ {
		card := hand.Cards[i]
		fmt.Printf("%d: %s%s\n", i+1, card.Value, card.Suit)
	}

	fmt.Printf("\nRemaining cards in deck: %d\n", deck.Count())
}
