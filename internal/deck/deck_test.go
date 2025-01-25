package deck

import (
	"strings"
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
	tests := []struct {
		value    string
		suit     string
		expected *Card
	}{
		{"A", "♠", &Card{Suit: "♠", Value: "A"}},
		{"10", "♥", &Card{Suit: "♥", Value: "10"}},
		{"Joker", "Red", &Card{Suit: "Red", Value: "Joker"}},
	}

	for _, tt := range tests {
		s.Run(tt.value+tt.suit, func() {
			card := NewCard(tt.value, tt.suit)
			assert.Equal(s.T(), tt.expected, card)
		})
	}
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
	tests := []struct {
		masks    []string
		expected int
		excluded []*Card
	}{
		{
			masks:    []string{"A♠", "K♥"},
			expected: 50,
			excluded: []*Card{{Suit: "♠", Value: "A"}, {Suit: "♥", Value: "K"}},
		},
		{
			masks:    []string{"10♦", "J♣"},
			expected: 50,
			excluded: []*Card{{Suit: "♦", Value: "10"}, {Suit: "♣", Value: "J"}},
		},
		{
			masks:    []string{},
			expected: 52,
			excluded: []*Card{},
		},
	}

	for _, tt := range tests {
		s.Run(strings.Join(tt.masks, ","), func() {
			deck := NewDeck(tt.masks...)
			assert.Equal(s.T(), tt.expected, deck.Count())

			// Verify masked cards are not in the deck
			for _, card := range deck.Cards {
				for _, excluded := range tt.excluded {
					assert.NotEqual(s.T(), excluded, card)
				}
			}
		})
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

func (s *DeckTestSuite) TestTimes() {
	deck := NewDeck()

	// Test single copy
	singleCopy := deck.Times(1)
	assert.Equal(s.T(), 52, singleCopy.Count(), "Single copy should have 52 cards")
	assert.Equal(s.T(), deck.Cards, singleCopy.Cards, "Single copy should match original deck")

	// Test multiple copies
	doubleDeck := deck.Times(2)
	assert.Equal(s.T(), 104, doubleDeck.Count(), "Double deck should have 104 cards")

	// Verify first and second halves match original deck
	assert.Equal(s.T(), deck.Cards, doubleDeck.Cards[:52], "First half should match original deck")
	assert.Equal(s.T(), deck.Cards, doubleDeck.Cards[52:], "Second half should match original deck")

	// Test edge cases
	emptyDeck := deck.Times(0)
	assert.Equal(s.T(), 0, emptyDeck.Count(), "Zero copies should create empty deck")

	negativeDeck := deck.Times(-1)
	assert.Equal(s.T(), 0, negativeDeck.Count(), "Negative copies should create empty deck")
}

func (s *DeckTestSuite) TestTimesAndShuffle() {
	deck := NewDeck()
	doubleDeck := deck.Times(2)

	// Make copies of original deck halves
	firstHalf := make([]*Card, 52)
	secondHalf := make([]*Card, 52)
	copy(firstHalf, doubleDeck.Cards[:52])
	copy(secondHalf, doubleDeck.Cards[52:])

	// Shuffle the double deck
	doubleDeck.Shuffle()

	// Verify all cards are still present
	cardCount := make(map[*Card]int)
	for _, card := range doubleDeck.Cards {
		cardCount[card]++
	}

	// Verify each card appears exactly twice
	for _, card := range deck.Cards {
		assert.Equal(s.T(), 2, cardCount[card], "Each card should appear exactly twice")
	}

	// Verify the deck is shuffled (both halves are mixed)
	firstHalfShuffled := false
	secondHalfShuffled := false

	for i := 0; i < 52; i++ {
		if doubleDeck.Cards[i] != firstHalf[i] {
			firstHalfShuffled = true
		}
		if doubleDeck.Cards[i+52] != secondHalf[i] {
			secondHalfShuffled = true
		}
	}

	assert.True(s.T(), firstHalfShuffled || secondHalfShuffled,
		"Shuffle should mix both halves of the double deck")
}

func (s *DeckTestSuite) TestComboCount() {
	tests := []struct {
		name      string
		deck      *Deck
		drawCount int
		expected  int
	}{
		{
			name:      "Standard deck draw 5",
			deck:      NewDeck(),
			drawCount: 5,
			expected:  2598960, // C(52,5)
		},
		{
			name:      "Standard deck draw 7",
			deck:      NewDeck(),
			drawCount: 7,
			expected:  133784560, // C(52,7)
		},
		{
			name:      "Deck with jokers draw 5",
			deck:      NewDeckWithJokers(),
			drawCount: 5,
			expected:  3162510, // C(54,5)
		},
		{
			name:      "Draw all cards",
			deck:      NewDeck(),
			drawCount: 52,
			expected:  1, // C(52,52)
		},
		{
			name:      "Draw 1 card",
			deck:      NewDeck(),
			drawCount: 1,
			expected:  52, // C(52,1)
		},
		{
			name:      "Draw more than deck size",
			deck:      NewDeck(),
			drawCount: 53,
			expected:  0,
		},
		{
			name:      "Draw zero cards",
			deck:      NewDeck(),
			drawCount: 0,
			expected:  0,
		},
		{
			name:      "Negative draw count",
			deck:      NewDeck(),
			drawCount: -1,
			expected:  0,
		},
		{
			name:      "Empty deck",
			deck:      &Deck{Cards: []*Card{}},
			drawCount: 5,
			expected:  0,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			result := tt.deck.ComboCount(tt.drawCount)
			assert.Equal(s.T(), tt.expected, result)
		})
	}
}

func (s *DeckTestSuite) TestDrawWithLimitHands() {
	deck := NewDeck()

	tests := []struct {
		name      string
		deck      *Deck
		drawCount int
		limit     int
		expected  int // expected number of hands
		handSize  int // expected cards per hand
	}{
		{
			name:      "Standard 5-card hands",
			deck:      deck,
			drawCount: 5,
			limit:     10,
			expected:  10,
			handSize:  5,
		},
		{
			name:      "Single hand",
			deck:      deck,
			drawCount: 7,
			limit:     1,
			expected:  1,
			handSize:  7,
		},
		{
			name:      "Limit exceeds available cards",
			deck:      deck,
			drawCount: 5,
			limit:     20,
			expected:  10, // 52/5 = 10 hands
			handSize:  5,
		},
		{
			name:      "Limit exceeds combo count",
			deck:      deck,
			drawCount: 5,
			limit:     3000000, // More than C(52,5) = 2,598,960
			expected:  2598960, // Should be capped at ComboCount
			handSize:  5,
		},
		{
			name:      "Draw all cards as single hand",
			deck:      deck,
			drawCount: 52,
			limit:     1,
			expected:  1,
			handSize:  52,
		},
		{
			name:      "Invalid draw count",
			deck:      deck,
			drawCount: 0,
			limit:     5,
			expected:  0,
		},
		{
			name:      "Invalid limit",
			deck:      deck,
			drawCount: 5,
			limit:     0,
			expected:  0,
		},
		{
			name:      "Draw count exceeds deck size",
			deck:      deck,
			drawCount: 53,
			limit:     1,
			expected:  0,
		},
		{
			name:      "Empty deck",
			deck:      &Deck{Cards: []*Card{}},
			drawCount: 5,
			limit:     1,
			expected:  0,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			hands := tt.deck.DrawWithLimitHands(tt.drawCount, tt.limit)
			assert.Equal(s.T(), tt.expected, len(hands))

			if tt.expected > 0 {
				// Verify all hands have correct size
				for _, hand := range hands {
					assert.Equal(s.T(), tt.handSize, len(hand.Cards))
				}

				// Verify all cards are unique across hands
				seenCards := make(map[*Card]bool)
				for _, hand := range hands {
					for _, card := range hand.Cards {
						assert.False(s.T(), seenCards[card], "Duplicate card found")
						seenCards[card] = true
					}
				}
			}
		})
	}
}
