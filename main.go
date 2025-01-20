package main

import (
	"fmt"

	"github.com/genewoo/joker/internal/deck"
)

// DealStrategy defines the interface for dealing cards
type DealStrategy interface {
	Deal(deck *deck.Deck, numCards int) (*deck.Hand, error)
}

// StandardDealer implements the standard dealing strategy
type StandardDealer struct{}

func (sd *StandardDealer) Deal(d *deck.Deck, numCards int) (*deck.Hand, error) {
	if d.Count() == 0 {
		return nil, fmt.Errorf("cannot deal from empty deck")
	}
	if numCards > d.Count() {
		numCards = d.Count()
	}
	cards := d.Cards[:numCards]
	d.Cards = d.Cards[numCards:]

	hand := deck.NewHand()
	for _, card := range cards {
		hand.AddCard(card)
	}
	return hand, nil
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
	for i := 0; i < hand.Count(); i++ {
		card := hand.Cards[i]
		fmt.Printf("%d: %s%s\n", i+1, card.Value, card.Suit)
	}

	fmt.Printf("\nRemaining cards in deck: %d\n", deck.Count())
}
