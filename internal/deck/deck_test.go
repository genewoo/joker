package deck

import (
	"testing"
)

func TestNewCard(t *testing.T) {
	card := NewCard("A", "♠")
	if card.Suit != "♠" || card.Value != "A" {
		t.Errorf("Expected A♠, got %s%s", card.Value, card.Suit)
	}
}

func TestNewDeck(t *testing.T) {
	deck := NewDeck()
	if deck.Count() != 52 {
		t.Errorf("Expected 52 cards, got %d", deck.Count())
	}
}

func TestNewDeckWithJokers(t *testing.T) {
	deck := NewDeckWithJokers()
	if deck.Count() != 54 {
		t.Errorf("Expected 54 cards, got %d", deck.Count())
	}
}

func TestCount(t *testing.T) {
	deck := NewDeck()
	if deck.Count() != 52 {
		t.Errorf("Expected 52 cards, got %d", deck.Count())
	}
}

func TestShuffle(t *testing.T) {
	deck1 := NewDeck()
	deck2 := NewDeck()

	// Get initial order of both decks
	initialOrder1 := make([]*Card, len(deck1.Cards))
	initialOrder2 := make([]*Card, len(deck2.Cards))
	copy(initialOrder1, deck1.Cards)
	copy(initialOrder2, deck2.Cards)

	// Shuffle both decks
	deck1.Shuffle()
	deck2.Shuffle()

	// Check if shuffled decks are different from initial order
	sameOrder1 := true
	sameOrder2 := true
	for i := range deck1.Cards {
		if deck1.Cards[i] != initialOrder1[i] {
			sameOrder1 = false
			break
		}
	}
	for i := range deck2.Cards {
		if deck2.Cards[i] != initialOrder2[i] {
			sameOrder2 = false
			break
		}
	}

	if sameOrder1 || sameOrder2 {
		t.Error("Shuffle did not change the deck order")
	}
}

func TestNewDeckWithMasks(t *testing.T) {
	deck := NewDeck("A♠", "K♥")
	if deck.Count() != 50 {
		t.Errorf("Expected 50 cards after masking, got %d", deck.Count())
	}
}

func TestNewDeckWithJokersWithMasks(t *testing.T) {
	deck := NewDeckWithJokers("A♠", "K♥")
	if deck.Count() != 52 {
		t.Errorf("Expected 52 cards after masking, got %d", deck.Count())
	}
}
