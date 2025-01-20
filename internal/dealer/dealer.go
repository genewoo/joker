package dealer

import (
	"fmt"

	"github.com/genewoo/joker/internal/deck"
)

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
