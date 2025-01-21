package deck

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type HandsTestSuite struct {
	suite.Suite
}

func TestHandsSuite(t *testing.T) {
	suite.Run(t, new(HandsTestSuite))
}

func (s *HandsTestSuite) TestNewHand() {
	hand := NewHand()
	assert.NotNil(s.T(), hand)
	assert.Equal(s.T(), 0, hand.Count())
	assert.Empty(s.T(), hand.Cards)
}

func (s *HandsTestSuite) TestAddCard() {
	tests := []struct {
		name     string
		card     *Card
		expected *Card
	}{
		{"Add regular card", &Card{Suit: "♠", Value: "A"}, &Card{Suit: "♠", Value: "A"}},
		{"Add joker", &Card{Suit: "Red", Value: "Joker"}, &Card{Suit: "Red", Value: "Joker"}},
		{"Add nil card", nil, nil},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			hand := NewHand()
			hand.AddCard(tt.card)

			if tt.card != nil {
				assert.Equal(s.T(), 1, hand.Count())
				assert.Equal(s.T(), tt.expected, hand.Cards[0])
			} else {
				assert.Equal(s.T(), 0, hand.Count())
			}
		})
	}
}

func (s *HandsTestSuite) TestRemoveCard() {
	tests := []struct {
		name        string
		setup       func() *Hand
		index       int
		expected    *Card
		expectedLen int
	}{
		{
			"Remove first card",
			func() *Hand {
				h := NewHand()
				h.AddCard(&Card{Suit: "♠", Value: "A"})
				h.AddCard(&Card{Suit: "♥", Value: "K"})
				return h
			},
			0,
			&Card{Suit: "♠", Value: "A"},
			1,
		},
		{
			"Remove last card",
			func() *Hand {
				h := NewHand()
				h.AddCard(&Card{Suit: "♠", Value: "A"})
				h.AddCard(&Card{Suit: "♥", Value: "K"})
				return h
			},
			1,
			&Card{Suit: "♥", Value: "K"},
			1,
		},
		{
			"Remove from empty hand",
			func() *Hand { return NewHand() },
			0,
			nil,
			0,
		},
		{
			"Remove with negative index",
			func() *Hand {
				h := NewHand()
				h.AddCard(&Card{Suit: "♠", Value: "A"})
				return h
			},
			-1,
			nil,
			1,
		},
		{
			"Remove with out of bounds index",
			func() *Hand {
				h := NewHand()
				h.AddCard(&Card{Suit: "♠", Value: "A"})
				return h
			},
			1,
			nil,
			1,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			hand := tt.setup()
			card := hand.RemoveCard(tt.index)
			assert.Equal(s.T(), tt.expected, card)
			assert.Equal(s.T(), tt.expectedLen, hand.Count())
		})
	}
}

func (s *HandsTestSuite) TestCount() {
	tests := []struct {
		name     string
		setup    func() *Hand
		expected int
	}{
		{"Empty hand", func() *Hand { return NewHand() }, 0},
		{"Single card", func() *Hand {
			h := NewHand()
			h.AddCard(&Card{Suit: "♠", Value: "A"})
			return h
		}, 1},
		{"Multiple cards", func() *Hand {
			h := NewHand()
			h.AddCard(&Card{Suit: "♠", Value: "A"})
			h.AddCard(&Card{Suit: "♥", Value: "K"})
			h.AddCard(&Card{Suit: "♦", Value: "Q"})
			return h
		}, 3},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			hand := tt.setup()
			assert.Equal(s.T(), tt.expected, hand.Count())
		})
	}
}

func (s *HandsTestSuite) TestClear() {
	tests := []struct {
		name  string
		setup func() *Hand
	}{
		{"Clear empty hand", func() *Hand { return NewHand() }},
		{"Clear single card", func() *Hand {
			h := NewHand()
			h.AddCard(&Card{Suit: "♠", Value: "A"})
			return h
		}},
		{"Clear multiple cards", func() *Hand {
			h := NewHand()
			h.AddCard(&Card{Suit: "♠", Value: "A"})
			h.AddCard(&Card{Suit: "♥", Value: "K"})
			h.AddCard(&Card{Suit: "♦", Value: "Q"})
			return h
		}},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			hand := tt.setup()
			hand.Clear()
			assert.Equal(s.T(), 0, hand.Count())
			assert.Empty(s.T(), hand.Cards)
		})
	}
}
