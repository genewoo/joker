package main

import (
	"fmt"

	"github.com/genewoo/joker/internal/deck"
)

// DealStrategy defines the interface for dealing cards
type DealStrategy interface {
	Deal(deck *deck.Deck, numCards int) ([]*deck.Card, error)
}

// StandardDealer implements the standard dealing strategy
type StandardDealer struct{}

func (sd *StandardDealer) Deal(deck *deck.Deck, numCards int) ([]*deck.Card, error) {
	if deck.Count() == 0 {
		return nil, fmt.Errorf("cannot deal from empty deck")
	}
	if numCards > deck.Count() {
		numCards = deck.Count()
	}
	cards := deck.Cards[:numCards]
	deck.Cards = deck.Cards[numCards:]
	return cards, nil
}

func main() {
	// Create a deck excluding specific cards
	deck := deck.NewDeck("A♠", "K♥", "7♦")
	deck.Shuffle()

	fmt.Println("Created deck excluding A♠, K♥, and 7♦")

	dealer := &StandardDealer{}

	// Deal 5 cards
	hand, err := dealer.Deal(deck, 5)
	if err != nil {
		fmt.Println("Error dealing cards:", err)
		return
	}

	fmt.Println("\nDealt cards:")
	for i, card := range hand {
		fmt.Printf("%d: %s%s\n", i+1, card.Value, card.Suit)
	}

	fmt.Printf("\nRemaining cards in deck: %d\n", deck.Count())
}
