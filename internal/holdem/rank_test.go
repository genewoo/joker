package holdem

import (
	"testing"

	"github.com/genewoo/joker/internal/deck"
	"github.com/stretchr/testify/assert"
)

func TestRankHand(t *testing.T) {
	tests := []struct {
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
			expectedRank: HandStrength{Rank: RoyalFlush},
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
			expectedRank: HandStrength{Rank: StraightFlush},
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
			expectedRank: HandStrength{Rank: FourOfAKind},
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
			expectedRank: HandStrength{Rank: FullHouse},
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
			expectedRank: HandStrength{Rank: Flush},
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
			expectedRank: HandStrength{Rank: Straight},
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
			expectedRank: HandStrength{Rank: ThreeOfAKind},
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
			expectedRank: HandStrength{Rank: TwoPair},
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
			expectedRank: HandStrength{Rank: TwoPair},
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
			expectedRank: HandStrength{Rank: OnePair},
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
			expectedRank: HandStrength{Rank: HighCard},
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rank, _ := RankHand(tt.playerCards, tt.communityCards)
			assert.Equal(t, tt.expectedRank.Rank, rank.Rank)
		})
	}
}
