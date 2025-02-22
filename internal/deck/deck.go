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
		cards = append(cards, &Card{Suit: "BW", Value: "Joker"})
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

// Times creates a new deck with multiple copies of the current deck
// count: number of copies to create (must be positive)
// Returns a new Deck containing count copies of the current deck's cards
func (d *Deck) Times(count int) *Deck {
	if count <= 0 {
		return &Deck{Cards: []*Card{}}
	}

	var cards []*Card
	for i := 0; i < count; i++ {
		cards = append(cards, d.Cards...)
	}
	return &Deck{Cards: cards}
}

// ComboCount calculates the number of possible combinations when drawing a specified number of cards
// drawCount: number of cards to draw (must be positive and not exceed deck size)
// Returns the number of possible combinations as an integer
func (d *Deck) ComboCount(drawCount int) int {
	n := len(d.Cards)
	if drawCount <= 0 || drawCount > n {
		return 0
	}

	// Use the multiplicative formula to avoid large factorial calculations
	result := 1
	for i := 1; i <= drawCount; i++ {
		result = result * (n - drawCount + i) / i
	}
	return result
}

// DrawWithLimitHands generates randomized hands of cards
// drawCount: number of cards per hand (must be positive and not exceed deck size)
// limit: maximum number of hands to generate (must be positive)
// Returns a slice of *Hand
func (d *Deck) DrawWithLimitHands(drawCount, limit int) []*Hand {
	if drawCount <= 0 || limit <= 0 || drawCount > len(d.Cards) {
		return nil
	}

	// Calculate maximum possible combinations
	maxCombinations := d.ComboCount(drawCount)
	if limit > maxCombinations {
		limit = maxCombinations
	}

	if maxCombinations == 0 {
		return nil
	}

	// Shuffle the deck if needed
	if len(d.Cards) < 1 {
		return nil
	}

	// build a slice of hands
	hands := make([]*Hand, 0, limit)
	drawnHands := make(map[string]bool)

	for i := 0; i < limit; i++ {
		if len(drawnHands) >= maxCombinations {
			break // Stop if all combinations are drawn
		}
		d.Shuffle()
		drawnCards := make([]Card, drawCount)
		for i, cardPtr := range d.Cards[:drawCount] {
			drawnCards[i] = *cardPtr
		}
		currentHand := NewHandByCards(drawnCards...)
		handKey := currentHand.String()
		if !drawnHands[handKey] {
			hands = append(hands, currentHand)
			drawnHands[handKey] = true
		} else {
			i-- // Decrement counter to retry drawing a unique hand
		}
	}

	return hands
}

// String returns a string representation of a card.
func (c *Card) String() string {
	return c.Value + c.Suit
}
