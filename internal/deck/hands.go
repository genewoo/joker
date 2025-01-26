package deck

import (
	"sort"
	"strings"
)

// Organizer defines the interface for sorting hands
type Organizer interface {
	Sort(cards []*Card)
}

// DefaultOrganizer sorts cards by value and suit, with jokers at the end
type DefaultOrganizer struct{}

var valueOrder = map[string]int{
	"A":  14,
	"K":  13,
	"Q":  12,
	"J":  11,
	"10": 10,
	"9":  9,
	"8":  8,
	"7":  7,
	"6":  6,
	"5":  5,
	"4":  4,
	"3":  3,
	"2":  2,
}

var suitOrder = map[string]int{
	"♠":   4,
	"♥":   3,
	"♦":   2,
	"♣":   1,
	"Red": 0,
	"BW":  -1,
}

// Sort implements the Organizer interface
func (o *DefaultOrganizer) Sort(cards []*Card) {
	sort.Slice(cards, func(i, j int) bool {
		// Jokers go to the beginning
		if cards[i].Value == "Joker" && cards[j].Value != "Joker" {
			return true
		}
		if cards[j].Value == "Joker" && cards[i].Value != "Joker" {
			return false
		}

		// Compare values
		if valueOrder[cards[i].Value] != valueOrder[cards[j].Value] {
			return valueOrder[cards[i].Value] > valueOrder[cards[j].Value]
		}

		// Compare suits
		return suitOrder[cards[i].Suit] > suitOrder[cards[j].Suit]
	})
}

// Hand represents a player's hand of cards
type Hand struct {
	Cards     []*Card
	organizer Organizer
}

// NewHand creates a new hand with default organizer and optional initial cards
func NewHand(cards ...*Card) *Hand {
	return &Hand{
		Cards:     cards,
		organizer: &DefaultOrganizer{},
	}
}

// NewHandByCards creates a new hand with default organizer and initial cards
func NewHandByCards(cards ...Card) *Hand {
	cardsPointers := make([]*Card, len(cards))
	for i := range cards {
		cardsPointers[i] = &cards[i]
	}
	return &Hand{
		Cards:     cardsPointers,
		organizer: &DefaultOrganizer{},
	}
}

// SetOrganizer sets a custom organizer for the hand
func (h *Hand) SetOrganizer(organizer Organizer) {
	h.organizer = organizer
}

// Sort sorts the hand's cards using the current organizer
func (h *Hand) Sort() {
	h.organizer.Sort(h.Cards)
}

// AddCard adds a card to the hand
func (h *Hand) AddCard(card *Card) {
	if card != nil {
		h.Cards = append(h.Cards, card)
	}
}

// RemoveCard removes a card from the hand at the specified index
// Returns the removed card or nil if index is out of bounds
func (h *Hand) RemoveCard(index int) *Card {
	if index < 0 || index >= len(h.Cards) {
		return nil
	}

	card := h.Cards[index]
	h.Cards = append(h.Cards[:index], h.Cards[index+1:]...)
	return card
}

// Count returns the number of cards in the hand
func (h *Hand) Count() int {
	return len(h.Cards)
}

// Clear removes all cards from the hand
func (h *Hand) Clear() {
	h.Cards = []*Card{}
}

// IndexOf returns the index of a card in the hand
func (h *Hand) IndexOf(card *Card) int {
	for i, c := range h.Cards {
		if c == card {
			return i
		}
	}
	return -1
}

// String returns a string representation of the hand, with cards sorted.
func (h *Hand) String() string {
	h.Sort() // Sort the hand before stringifying
	var cardStrings []string
	for _, card := range h.Cards {
		cardStrings = append(cardStrings, card.String()) // Assuming Card has a String() method
	}
	return strings.Join(cardStrings, ",")
}
