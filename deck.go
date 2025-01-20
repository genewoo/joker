package main

import (
	"math/rand"
	"time"
)

// Card represents a playing card
type Card struct {
	Suit  string
	Value string
}

// Deck represents a collection of cards
type Deck struct {
	cards []Card
}

// NewDeck creates a new deck of cards, optionally excluding cards that match the provided masks
// Each mask should be in the format "ValueSuit" (e.g., "A♠", "10♥")
func NewDeck(masks ...string) *Deck {
	suits := []string{"♠", "♥", "♦", "♣"}
	values := []string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}

	// Convert masks to a map for O(1) lookups
	maskMap := make(map[string]bool)
	for _, mask := range masks {
		maskMap[mask] = true
	}

	var cards []Card
	for _, suit := range suits {
		for _, value := range values {
			cardStr := value + suit
			if !maskMap[cardStr] {
				cards = append(cards, Card{Suit: suit, Value: value})
			}
		}
	}
	return &Deck{cards: cards}
}

// Shuffle randomizes the order of cards in the deck
func (d *Deck) Shuffle() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(d.cards), func(i, j int) {
		d.cards[i], d.cards[j] = d.cards[j], d.cards[i]
	})
}
