package holdem

import (
	"testing"

	"github.com/genewoo/joker/internal/deck"
	"github.com/stretchr/testify/assert"
)

var handTestCases = []struct {
	name           string
	playerCards    []*deck.Card
	communityCards []*deck.Card
	expectedRank   HandStrength
}{
	{
		name: "Royal Flush",
		playerCards: []*deck.Card{
			{Value: "A", Suit: "♠"},
			{Value: "K", Suit: "♠"},
		},
		communityCards: []*deck.Card{
			{Value: "Q", Suit: "♠"},
			{Value: "J", Suit: "♠"},
			{Value: "10", Suit: "♠"},
			{Value: "2", Suit: "♥"},
			{Value: "3", Suit: "♦"},
		},
		expectedRank: HandStrength{
			Rank:   RoyalFlush,
			Values: []int{14},
		},
	},
	{
		name: "Straight Flush",
		playerCards: []*deck.Card{
			{Value: "9", Suit: "♠"},
			{Value: "8", Suit: "♠"},
		},
		communityCards: []*deck.Card{
			{Value: "7", Suit: "♠"},
			{Value: "6", Suit: "♠"},
			{Value: "5", Suit: "♠"},
			{Value: "2", Suit: "♥"},
			{Value: "3", Suit: "♦"},
		},
		expectedRank: HandStrength{
			Rank:   StraightFlush,
			Values: []int{9},
		},
	},
	{
		name: "Four of a Kind",
		playerCards: []*deck.Card{
			{Value: "A", Suit: "♠"},
			{Value: "A", Suit: "♥"},
		},
		communityCards: []*deck.Card{
			{Value: "A", Suit: "♦"},
			{Value: "A", Suit: "♣"},
			{Value: "K", Suit: "♠"},
			{Value: "2", Suit: "♥"},
			{Value: "3", Suit: "♦"},
		},
		expectedRank: HandStrength{
			Rank:   FourOfAKind,
			Values: []int{14, 13},
		},
	},
	{
		name: "Full House",
		playerCards: []*deck.Card{
			{Value: "A", Suit: "♠"},
			{Value: "A", Suit: "♥"},
		},
		communityCards: []*deck.Card{
			{Value: "K", Suit: "♦"},
			{Value: "K", Suit: "♣"},
			{Value: "K", Suit: "♠"},
			{Value: "2", Suit: "♥"},
			{Value: "3", Suit: "♦"},
		},
		expectedRank: HandStrength{
			Rank:   FullHouse,
			Values: []int{13, 14},
		},
	},
	{
		name: "Flush",
		playerCards: []*deck.Card{
			{Value: "A", Suit: "♠"},
			{Value: "K", Suit: "♠"},
		},
		communityCards: []*deck.Card{
			{Value: "Q", Suit: "♠"},
			{Value: "J", Suit: "♠"},
			{Value: "9", Suit: "♠"},
			{Value: "2", Suit: "♥"},
			{Value: "3", Suit: "♦"},
		},
		expectedRank: HandStrength{
			Rank:   Flush,
			Values: []int{14, 13, 12, 11, 9},
		},
	},
	{
		name: "Straight 5-9",
		playerCards: []*deck.Card{
			{Value: "9", Suit: "♠"},
			{Value: "8", Suit: "♥"},
		},
		communityCards: []*deck.Card{
			{Value: "7", Suit: "♦"},
			{Value: "6", Suit: "♣"},
			{Value: "5", Suit: "♠"},
			{Value: "2", Suit: "♥"},
			{Value: "3", Suit: "♦"},
		},
		expectedRank: HandStrength{
			Rank:   Straight,
			Values: []int{9},
		},
	},
	{
		name: "Straight A-5",
		playerCards: []*deck.Card{
			{Value: "A", Suit: "♠"},
			{Value: "4", Suit: "♥"},
		},
		communityCards: []*deck.Card{
			{Value: "7", Suit: "♦"},
			{Value: "8", Suit: "♣"},
			{Value: "5", Suit: "♠"},
			{Value: "2", Suit: "♥"},
			{Value: "3", Suit: "♦"},
		},
		expectedRank: HandStrength{Rank: Straight},
	},
	{
		name: "Straight T-A",
		playerCards: []*deck.Card{
			{Value: "10", Suit: "♠"},
			{Value: "K", Suit: "♥"},
		},
		communityCards: []*deck.Card{
			{Value: "7", Suit: "♦"},
			{Value: "8", Suit: "♣"},
			{Value: "Q", Suit: "♠"},
			{Value: "J", Suit: "♥"},
			{Value: "A", Suit: "♦"},
		},
		expectedRank: HandStrength{Rank: Straight},
	},
	{
		name: "Three of a Kind",
		playerCards: []*deck.Card{
			{Value: "A", Suit: "♠"},
			{Value: "A", Suit: "♥"},
		},
		communityCards: []*deck.Card{
			{Value: "A", Suit: "♦"},
			{Value: "K", Suit: "♣"},
			{Value: "Q", Suit: "♠"},
			{Value: "2", Suit: "♥"},
			{Value: "3", Suit: "♦"},
		},
		expectedRank: HandStrength{
			Rank:   ThreeOfAKind,
			Values: []int{14, 13, 12},
		},
	},
	{
		name: "Two Pair",
		playerCards: []*deck.Card{
			{Value: "A", Suit: "♠"},
			{Value: "A", Suit: "♥"},
		},
		communityCards: []*deck.Card{
			{Value: "K", Suit: "♦"},
			{Value: "K", Suit: "♣"},
			{Value: "Q", Suit: "♠"},
			{Value: "2", Suit: "♥"},
			{Value: "3", Suit: "♦"},
		},
		expectedRank: HandStrength{
			Rank:   TwoPair,
			Values: []int{14, 13, 12},
		},
	},
	{
		name: "Two Pair Common",
		playerCards: []*deck.Card{
			{Value: "A", Suit: "♠"},
			{Value: "4", Suit: "♥"},
		},
		communityCards: []*deck.Card{
			{Value: "K", Suit: "♦"},
			{Value: "K", Suit: "♣"},
			{Value: "Q", Suit: "♠"},
			{Value: "2", Suit: "♥"},
			{Value: "Q", Suit: "♦"},
		},
		expectedRank: HandStrength{
			Rank:   TwoPair,
			Values: []int{13, 12, 14},
		},
	},
	{
		name: "Three Pair - Only Select Top 2",
		playerCards: []*deck.Card{
			{Value: "A", Suit: "♠"},
			{Value: "4", Suit: "♥"},
		},
		communityCards: []*deck.Card{
			{Value: "K", Suit: "♦"},
			{Value: "K", Suit: "♣"},
			{Value: "J", Suit: "♠"},
			{Value: "A", Suit: "♥"},
			{Value: "J", Suit: "♦"},
		},
		expectedRank: HandStrength{
			Rank:   TwoPair,
			Values: []int{14, 13, 11},
		},
	},
	{
		name: "One Pair",
		playerCards: []*deck.Card{
			{Value: "A", Suit: "♠"},
			{Value: "A", Suit: "♥"},
		},
		communityCards: []*deck.Card{
			{Value: "K", Suit: "♦"},
			{Value: "Q", Suit: "♣"},
			{Value: "J", Suit: "♠"},
			{Value: "2", Suit: "♥"},
			{Value: "3", Suit: "♦"},
		},
		expectedRank: HandStrength{
			Rank:   OnePair,
			Values: []int{14, 13, 12, 11},
		},
	},
	{
		name: "High Card",
		playerCards: []*deck.Card{
			{Value: "A", Suit: "♠"},
			{Value: "K", Suit: "♥"},
		},
		communityCards: []*deck.Card{
			{Value: "Q", Suit: "♦"},
			{Value: "J", Suit: "♣"},
			{Value: "9", Suit: "♠"},
			{Value: "2", Suit: "♥"},
			{Value: "3", Suit: "♦"},
		},
		expectedRank: HandStrength{
			Rank:   HighCard,
			Values: []int{14, 13, 12, 11, 9},
		},
	},
	{
		name: "Invalid Hand - too few player cards",
		playerCards: []*deck.Card{
			{Value: "A", Suit: "♠"},
		},
		communityCards: []*deck.Card{
			{Value: "K", Suit: "♦"},
			{Value: "Q", Suit: "♣"},
			{Value: "J", Suit: "♠"},
			{Value: "2", Suit: "♥"},
			{Value: "3", Suit: "♦"},
		},
		expectedRank: HandStrength{Rank: InvalidHand},
	},
	{
		name: "Invalid Hand - too many community cards",
		playerCards: []*deck.Card{
			{Value: "A", Suit: "♠"},
			{Value: "K", Suit: "♥"},
		},
		communityCards: []*deck.Card{
			{Value: "Q", Suit: "♦"},
			{Value: "J", Suit: "♣"},
			{Value: "9", Suit: "♠"},
			{Value: "2", Suit: "♥"},
			{Value: "3", Suit: "♦"},
			{Value: "4", Suit: "♠"},
		},
		expectedRank: HandStrength{Rank: InvalidHand},
	},
	{
		name:        "Edge Cases",
		playerCards: nil,
		communityCards: []*deck.Card{
			{Value: "A", Suit: "♠"},
			{Value: "K", Suit: "♠"},
			{Value: "Q", Suit: "♠"},
			{Value: "J", Suit: "♠"},
			{Value: "10", Suit: "♠"},
		},
		expectedRank: HandStrength{Rank: InvalidHand},
	},
	{
		name: "Duplicate cards",
		playerCards: []*deck.Card{
			{Value: "A", Suit: "♠"},
			{Value: "A", Suit: "♠"},
		},
		communityCards: []*deck.Card{
			{Value: "K", Suit: "♠"},
			{Value: "Q", Suit: "♠"},
			{Value: "J", Suit: "♠"},
			{Value: "10", Suit: "♠"},
		},
		expectedRank: HandStrength{Rank: InvalidHand},
	},
	{
		name: "Invalid card value",
		playerCards: []*deck.Card{
			{Value: "A", Suit: "♠"},
			{Value: "K", Suit: "♠"},
		},
		communityCards: []*deck.Card{
			{Value: "1", Suit: "♠"}, // Invalid value
			{Value: "Q", Suit: "♠"},
			{Value: "J", Suit: "♠"},
			{Value: "10", Suit: "♠"},
		},
		expectedRank: HandStrength{Rank: InvalidHand},
	},
	{
		name: "Invalid card suit",
		playerCards: []*deck.Card{
			{Value: "A", Suit: "♠"},
			{Value: "K", Suit: "♠"},
		},
		communityCards: []*deck.Card{
			{Value: "Q", Suit: "X"}, // Invalid suit
			{Value: "J", Suit: "♠"},
			{Value: "10", Suit: "♠"},
		},
		expectedRank: HandStrength{Rank: InvalidHand},
	},
	{
		name:        "Empty player cards",
		playerCards: []*deck.Card{},
		communityCards: []*deck.Card{
			{Value: "A", Suit: "♠"},
			{Value: "K", Suit: "♠"},
			{Value: "Q", Suit: "♠"},
			{Value: "J", Suit: "♠"},
			{Value: "10", Suit: "♠"},
		},
		expectedRank: HandStrength{Rank: InvalidHand},
	},
}

func TestRankHand(t *testing.T) {
	for _, tt := range handTestCases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			rank, cards := RankHand(tt.playerCards, tt.communityCards)
			if tt.expectedRank.Rank == InvalidHand {
				assert.Nil(t, cards)
			} else {
				assert.NotNil(t, cards)
			}
			assert.Equal(t, tt.expectedRank.Rank, rank.Rank)
		})
	}
}

func TestHandStrengthCompare(t *testing.T) {
	tests := []struct {
		name     string
		h1       HandStrength
		h2       HandStrength
		expected int
	}{
		{
			name: "Different ranks - h1 stronger",
			h1: HandStrength{
				Rank:   FullHouse,
				Values: []int{10, 10},
			},
			h2: HandStrength{
				Rank:   Flush,
				Values: []int{14, 13, 12, 11, 9},
			},
			expected: 1,
		},
		{
			name: "Different ranks - h2 stronger",
			h1: HandStrength{
				Rank:   OnePair,
				Values: []int{14, 13, 12, 11},
			},
			h2: HandStrength{
				Rank:   TwoPair,
				Values: []int{13, 12, 14},
			},
			expected: -1,
		},
		{
			name: "Same rank - h1 stronger values",
			h1: HandStrength{
				Rank:   Flush,
				Values: []int{14, 13, 12, 11, 9},
			},
			h2: HandStrength{
				Rank:   Flush,
				Values: []int{13, 12, 11, 10, 9},
			},
			expected: 1,
		},
		{
			name: "Same rank - h2 stronger values",
			h1: HandStrength{
				Rank:   TwoPair,
				Values: []int{13, 12, 10},
			},
			h2: HandStrength{
				Rank:   TwoPair,
				Values: []int{14, 13, 12},
			},
			expected: -1,
		},
		{
			name: "Same rank and values",
			h1: HandStrength{
				Rank:   HighCard,
				Values: []int{14, 13, 12, 11, 9},
			},
			h2: HandStrength{
				Rank:   HighCard,
				Values: []int{14, 13, 12, 11, 9},
			},
			expected: 0,
		},
		{
			name: "Same rank - High Rank",
			h1: HandStrength{
				Rank:   FullHouse,
				Values: []int{14, 12},
			},
			h2: HandStrength{
				Rank:   FullHouse,
				Values: []int{14, 13},
			},
			expected: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := tt.h1.Compare(tt.h2)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSmartRankHand(t *testing.T) {
	for _, tt := range handTestCases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			rank, cards := SmartRankHand(tt.playerCards, tt.communityCards)
			if tt.expectedRank.Rank == InvalidHand {
				assert.Nil(t, cards)
			} else {
				assert.NotNil(t, cards)
			}
			assert.Equal(t, tt.expectedRank.Rank, rank.Rank)
		})
	}
}
