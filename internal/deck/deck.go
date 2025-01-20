// Package deck provides functionality for creating and manipulating decks of cards
// Changes made:
// - Made original NewDeck private as newDeck
// - Added public NewDeck that creates decks without jokers
// - Added NewDeckWithJokers for decks with jokers
// - Added detailed documentation for all public methods
package deck

import (
	"math/rand"
	"time"
)

// Card represents a playing card
type Card struct {
	Suit  string
	Value string
}

// NewCard creates a new Card instance with the specified value and suit
// value: The card's value (e.g., "A", "2", "J", "Q", "K")
// suit: The card's suit (e.g., "♠", "♥", "♦", "♣")
// Returns a pointer to the newly created Card
func NewCard(value, suit string) *Card {
	return &Card{
		Suit:  suit,
		Value: value,
	}
}

// Deck represents a collection of cards
type Deck struct {
	Cards []*Card
}

// Count returns the number of remaining cards in the deck
// Returns the current number of cards in the deck as an integer
func (d *Deck) Count() int {
	return len(d.Cards)
}

// newDeck creates a new deck of cards, optionally excluding cards that match the provided masks
// Each mask should be in the format "ValueSuit" (e.g., "A♠", "10♥")
// includeJokers determines whether to include two joker cards (Red and White)
func newDeck(includeJokers bool, masks ...string) *Deck {
	suits := []string{"♠", "♥", "♦", "♣"}
	values := []string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}

	// Convert masks to a map for O(1) lookups
	maskMap := make(map[string]bool)
	for _, mask := range masks {
		maskMap[mask] = true
	}

	var cards []*Card
	for _, suit := range suits {
		for _, value := range values {
			cardStr := value + suit
			if !maskMap[cardStr] {
				cards = append(cards, &Card{Suit: suit, Value: value})
			}
		}
	}

	if includeJokers {
		cards = append(cards, &Card{Suit: "Red", Value: "Joker"})
		cards = append(cards, &Card{Suit: "White", Value: "Joker"})
	}
	return &Deck{Cards: cards}
}

// NewDeck creates a new standard deck of 52 cards without jokers
// masks: Optional list of cards to exclude from the deck in "ValueSuit" format (e.g., "A♠", "10♥")
// Returns a pointer to the newly created Deck
func NewDeck(masks ...string) *Deck {
	return newDeck(false, masks...)
}

// NewDeckWithJokers creates a new deck of 54 cards including two jokers (Red and White)
// masks: Optional list of cards to exclude from the deck in "ValueSuit" format (e.g., "A♠", "10♥")
// Returns a pointer to the newly created Deck
func NewDeckWithJokers(masks ...string) *Deck {
	return newDeck(true, masks...)
}

// Shuffle randomizes the order of cards in the deck using the Fisher-Yates algorithm
// The shuffle is seeded with the current time to ensure different results each time
func (d *Deck) Shuffle() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(d.Cards), func(i, j int) {
		d.Cards[i], d.Cards[j] = d.Cards[j], d.Cards[i]
	})
}
