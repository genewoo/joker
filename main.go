package main

import (
	"fmt"
)

// DealStrategy defines the interface for dealing cards
type DealStrategy interface {
	Deal(deck *Deck, numCards int) []Card
}

// StandardDealer implements the standard dealing strategy
type StandardDealer struct{}

func (sd *StandardDealer) Deal(deck *Deck, numCards int) []Card {
	if numCards > len(deck.cards) {
		numCards = len(deck.cards)
	}
	cards := deck.cards[:numCards]
	deck.cards = deck.cards[numCards:]
	return cards
}

func main() {
	// Create a deck excluding specific cards
	deck := NewDeck("A♠", "K♥", "7♦")
	deck.Shuffle()

	fmt.Println("Created deck excluding A♠, K♥, and 7♦")

	dealer := &StandardDealer{}

	// Deal 5 cards
	hand := dealer.Deal(deck, 5)
	fmt.Println("Dealt cards:")
	for _, card := range hand {
		fmt.Printf("%s%s\n", card.Value, card.Suit)
	}

	fmt.Printf("\nRemaining cards in deck: %d\n", len(deck.cards))
}
