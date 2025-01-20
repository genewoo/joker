package deck

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type DeckTestSuite struct {
	suite.Suite
}

func TestDeckSuite(t *testing.T) {
	suite.Run(t, new(DeckTestSuite))
}

func (s *DeckTestSuite) TestNewCard() {
	card := NewCard("A", "♠")
	assert.Equal(s.T(), "♠", card.Suit, "Suit should match")
	assert.Equal(s.T(), "A", card.Value, "Value should match")
}

func (s *DeckTestSuite) TestNewDeck() {
	deck := NewDeck()
	assert.Equal(s.T(), 52, deck.Count(), "Standard deck should have 52 cards")
}

func (s *DeckTestSuite) TestNewDeckWithJokers() {
	deck := NewDeckWithJokers()
	assert.Equal(s.T(), 54, deck.Count(), "Deck with jokers should have 54 cards")
}

func (s *DeckTestSuite) TestCount() {
	deck := NewDeck()
	assert.Equal(s.T(), 52, deck.Count(), "Standard deck should have 52 cards")
}

func (s *DeckTestSuite) TestShuffle() {
	deck1 := NewDeck()
	deck2 := NewDeck()

	initialOrder1 := make([]*Card, len(deck1.Cards))
	initialOrder2 := make([]*Card, len(deck2.Cards))
	copy(initialOrder1, deck1.Cards)
	copy(initialOrder2, deck2.Cards)

	deck1.Shuffle()
	deck2.Shuffle()

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

	assert.False(s.T(), sameOrder1 || sameOrder2, "Shuffle should change the deck order")
}

func (s *DeckTestSuite) TestNewDeckWithMasks() {
	deck := NewDeck("A♠", "K♥")
	assert.Equal(s.T(), 50, deck.Count(), "Deck with masks should have 50 cards")

	// Verify masked cards are not in the deck
	for _, card := range deck.Cards {
		assert.False(s.T(),
			card.Value == "A" && card.Suit == "♠",
			"A♠ should be masked")
		assert.False(s.T(),
			card.Value == "K" && card.Suit == "♥",
			"K♥ should be masked")
	}
}

func (s *DeckTestSuite) TestNewDeckWithJokersWithMasks() {
	deck := NewDeckWithJokers("A♠", "K♥")
	assert.Equal(s.T(), 52, deck.Count(), "Deck with jokers and masks should have 52 cards")

	// Verify both jokers are present
	jokerCount := 0
	for _, card := range deck.Cards {
		if card.Value == "Joker" {
			jokerCount++
		}
	}
	assert.Equal(s.T(), 2, jokerCount, "Deck should contain exactly 2 jokers")
}
