package guandan

import (
	"testing"

	"github.com/genewoo/joker/internal/deck"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CombinerTestSuite struct {
	suite.Suite
	combiner *Combiner
}

func (suite *CombinerTestSuite) SetupTest() {
	suite.combiner = NewCombiner("5") // Using 5 as current level
}

func TestCombinerSuite(t *testing.T) {
	suite.Run(t, new(CombinerTestSuite))
}

func (suite *CombinerTestSuite) TestEvaluateCombination() {
	tests := []struct {
		name     string
		cards    []*deck.Card
		expected CombinationType
	}{
		{
			name:     "Single",
			cards:    []*deck.Card{{Value: "9", Suit: "♠"}},
			expected: Single,
		},
		{
			name:     "Pair",
			cards:    []*deck.Card{{Value: "5", Suit: "♦"}, {Value: "5", Suit: "♣"}},
			expected: Pair,
		},
		{
			name: "Triple",
			cards: []*deck.Card{
				{Value: "6", Suit: "♠"},
				{Value: "6", Suit: "♥"},
				{Value: "6", Suit: "♦"},
			},
			expected: Triple,
		},
		{
			name: "Plate",
			cards: []*deck.Card{
				{Value: "3", Suit: "♥"}, {Value: "3", Suit: "♠"}, {Value: "3", Suit: "♣"},
				{Value: "4", Suit: "♣"}, {Value: "4", Suit: "♦"}, {Value: "4", Suit: "♥"},
			},
			expected: Plate,
		},
		{
			name: "Tube",
			cards: []*deck.Card{
				{Value: "2", Suit: "♠"}, {Value: "2", Suit: "♥"},
				{Value: "3", Suit: "♠"}, {Value: "3", Suit: "♣"},
				{Value: "4", Suit: "♣"}, {Value: "4", Suit: "♦"},
			},
			expected: Tube,
		},
		{
			name: "Full House",
			cards: []*deck.Card{
				{Value: "9", Suit: "♠"}, {Value: "9", Suit: "♥"}, {Value: "9", Suit: "♦"},
				{Value: "J", Suit: "♥"}, {Value: "J", Suit: "♣"},
			},
			expected: FullHouse,
		},
		{
			name: "Straight",
			cards: []*deck.Card{
				{Value: "8", Suit: "♥"}, {Value: "9", Suit: "♠"},
				{Value: "10", Suit: "♣"}, {Value: "J", Suit: "♥"},
				{Value: "Q", Suit: "♦"},
			},
			expected: Straight,
		},
		{
			name: "Bomb",
			cards: []*deck.Card{
				{Value: "4", Suit: "♠"}, {Value: "4", Suit: "♥"},
				{Value: "4", Suit: "♣"}, {Value: "4", Suit: "♦"},
			},
			expected: Bomb,
		},
		{
			name: "Straight Flush",
			cards: []*deck.Card{
				{Value: "10", Suit: "♠"}, {Value: "J", Suit: "♠"},
				{Value: "Q", Suit: "♠"}, {Value: "K", Suit: "♠"},
				{Value: "A", Suit: "♠"},
			},
			expected: StraightFlush,
		},
		{
			name: "Joker Bomb",
			cards: []*deck.Card{
				{Value: "Joker", Suit: "Red"},
				{Value: "Joker", Suit: "Black"},
				{Value: "Joker", Suit: "Red"},
				{Value: "Joker", Suit: "Black"},
			},
			expected: JokerBomb,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			// result := suite.combiner.EvaluateCombination(tt.cards)
			// assert.Equal(suite.T(), tt.expected, result.Type)
		})
	}
}

func (suite *CombinerTestSuite) TestInvalidCombination() {
	tests := []struct {
		name  string
		cards []*deck.Card
	}{
		{
			name:  "Empty hand",
			cards: []*deck.Card{},
		},
		{
			name: "Too many cards for single",
			cards: []*deck.Card{
				{Value: "9", Suit: "♠"},
				{Value: "10", Suit: "♠"},
			},
		},
		{
			name: "Invalid straight",
			cards: []*deck.Card{
				{Value: "8", Suit: "♥"}, {Value: "9", Suit: "♠"},
				{Value: "10", Suit: "♣"}, {Value: "J", Suit: "♥"},
				{Value: "K", Suit: "♦"},
			},
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			result := suite.combiner.EvaluateCombination(tt.cards)
			assert.Equal(suite.T(), InvalidCombination, result.Type)
		})
	}
}

func (suite *CombinerTestSuite) TestCurrentLevelHandling() {
	tests := []struct {
		name     string
		cards    []*deck.Card
		expected bool
	}{
		{
			name: "Valid straight with current level wrap",
			cards: []*deck.Card{
				{Value: "A", Suit: "♠"},
				{Value: "2", Suit: "♠"},
				{Value: "3", Suit: "♠"},
				{Value: "4", Suit: "♠"},
				{Value: "5", Suit: "♠"},
			},
			expected: true,
		},
		{
			name: "Invalid straight without current level wrap",
			cards: []*deck.Card{
				{Value: "A", Suit: "♠"},
				{Value: "2", Suit: "♠"},
				{Value: "3", Suit: "♠"},
				{Value: "4", Suit: "♠"},
				{Value: "6", Suit: "♠"},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			result := suite.combiner.EvaluateCombination(tt.cards)
			if tt.expected {
				assert.Equal(suite.T(), Straight, result.Type)
			} else {
				assert.NotEqual(suite.T(), Straight, result.Type)
			}
		})
	}
}
