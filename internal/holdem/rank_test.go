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
		expectedRank   HandRank
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
			expectedRank: RoyalFlush,
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
			expectedRank: StraightFlush,
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
			expectedRank: FourOfAKind,
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
			expectedRank: FullHouse,
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
			expectedRank: Flush,
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
			expectedRank: Straight,
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
			expectedRank: Straight,
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
			expectedRank: Straight,
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
			expectedRank: ThreeOfAKind,
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
			expectedRank: TwoPair,
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
			expectedRank: TwoPair,
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
			expectedRank: OnePair,
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
			expectedRank: HighCard,
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
			expectedRank: InvalidHand,
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
			expectedRank: InvalidHand,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rank, _ := RankHand(tt.playerCards, tt.communityCards)
			assert.Equal(t, tt.expectedRank.String(), rank.String())
		})
	}
}
