package dealer

import (
	"fmt"

	"github.com/genewoo/joker/internal/deck"
)

type DealStrategy interface {
	Deal(deck *deck.Deck, numCards, hands int) ([]*deck.Hand, error)
}

// StandardDealer implements the standard dealing strategy
type StandardDealer struct{}

func (sd *StandardDealer) Deal(d *deck.Deck, numCards, hands int) ([]*deck.Hand, error) {
	if d.Count() == 0 {
		return nil, fmt.Errorf("cannot deal from empty deck")
	}
	if numCards <= 0 {
		return nil, fmt.Errorf("numCards must be positive")
	}
	totalCards := numCards * hands
	if totalCards > d.Count() {
		return nil, fmt.Errorf("not enough cards in deck")
	}

	result := make([]*deck.Hand, hands)
	for i := 0; i < hands; i++ {
		result[i] = &deck.Hand{Cards: []*deck.Card{}}
	}

	for i := 0; i < numCards; i++ {
		for h := 0; h < hands; h++ {
			card := d.Cards[0]
			d.Cards = d.Cards[1:]
			result[h].AddCard(card)
		}
	}
	return result, nil
}
