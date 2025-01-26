package holdem

import (
	"fmt"
	"testing"

	"github.com/genewoo/joker/internal/deck"
	"github.com/stretchr/testify/assert"
)

func TestWinningCalculator(t *testing.T) {
	tests := []struct {
		name        string
		players     [][]*deck.Card
		simulations int
		limitCount  int
		expected    []float64
	}{
		//{"♠", "♥", "♦", "♣"}
		{
			name: "Two players - clear winner",
			players: [][]*deck.Card{
				{deck.NewCard("A", "♠"), deck.NewCard("K", "♠")}, // Player 1
				{deck.NewCard("2", "♥"), deck.NewCard("3", "♥")}, // Player 2
			},
			simulations: 1000,
			expected:    []float64{0.9, 0.05}, // Player 1 should win most of the time
		},
		{
			name: "Two players - Tie",
			players: [][]*deck.Card{
				{deck.NewCard("A", "♠"), deck.NewCard("A", "♦")}, // Player 1
				{deck.NewCard("A", "♥"), deck.NewCard("A", "♣")}, // Player 2
			},
			simulations: 1000,
			expected:    []float64{0.5, 0.5}, // Player 1,2 should be tir most of the time
		},
		// {
		// 	name: "Three players - balanced",
		// 	players: [][]*deck.Card{
		// 		{deck.NewCard("A", "♠"), deck.NewCard("K", "♠")},
		// 		{deck.NewCard("Q", "♠"), deck.NewCard("J", "♠")},
		// 		{deck.NewCard("T", "♠"), deck.NewCard("9", "♠")},
		// 	},
		// 	simulations: 1000,
		// 	expected:    []float64{0.5, 0.3, 0.2}, // Rough expected distribution
		// },
		// {
		// 	name: "Tie scenario",
		// 	players: [][]*deck.Card{
		// 		{deck.NewCard("A", "♠"), deck.NewCard("K", "♠")},
		// 		{deck.NewCard("A", "♥"), deck.NewCard("K", "♥")},
		// 	},
		// 	simulations: 1000,
		// 	expected:    []float64{0.5, 0.5}, // Equal chance of winning
		// },
		// {
		// 	name:        "No players",
		// 	players:     [][]*deck.Card{},
		// 	simulations: 1000,
		// 	expected:    nil,
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calc := NewWinningCalculator(tt.players, tt.simulations)
			probs := calc.CalculateWinProbabilities()
			fmt.Println(probs)

			assert.NotNil(t, probs)

			assert.Equal(t, len(tt.expected), len(probs))
			for i := range tt.expected {
				assert.InDelta(t, tt.expected[i], probs[i], 0.1)
			}
		})
	}
}

func TestFindWinners(t *testing.T) {
	tests := []struct {
		name     string
		hands    []HandStrength
		expected []int
	}{
		{
			name: "Single winner",
			hands: []HandStrength{
				{Rank: HighCard, Values: []int{14, 13, 12, 11, 9}},
				{Rank: OnePair, Values: []int{10, 9, 8, 7}},
			},
			expected: []int{1},
		},
		{
			name: "Tie",
			hands: []HandStrength{
				{Rank: OnePair, Values: []int{10, 9, 8, 7}},
				{Rank: OnePair, Values: []int{10, 9, 8, 7}},
			},
			expected: []int{0, 1},
		},
		{
			name:     "No hands",
			hands:    []HandStrength{},
			expected: nil,
		},
		{
			name: "Three players - single winner",
			hands: []HandStrength{
				{Rank: FullHouse, Values: []int{10, 5, 5}},
				{Rank: Flush, Values: []int{9, 8, 7, 6, 5}},
				{Rank: Straight, Values: []int{8, 7, 6, 5, 4}},
			},
			expected: []int{0},
		},
		{
			name: "Three players - tie for first",
			hands: []HandStrength{
				{Rank: TwoPair, Values: []int{10, 5, 9}},
				{Rank: TwoPair, Values: []int{10, 5, 9}},
				{Rank: OnePair, Values: []int{9, 9, 8, 7, 6}},
			},
			expected: []int{0, 1},
		},
		{
			name: "Four players - single winner",
			hands: []HandStrength{
				{Rank: FourOfAKind, Values: []int{10, 5}},
				{Rank: FullHouse, Values: []int{9, 8}},
				{Rank: Flush, Values: []int{8, 7, 6, 5, 4}},
				{Rank: Straight, Values: []int{7, 6, 5, 4, 3}},
			},
			expected: []int{0},
		},
		{
			name: "Four players - tie for first two",
			hands: []HandStrength{
				{Rank: ThreeOfAKind, Values: []int{10, 5, 9}},
				{Rank: ThreeOfAKind, Values: []int{10, 5, 9}},
				{Rank: ThreeOfAKind, Values: []int{10, 5, 7}},
				{Rank: TwoPair, Values: []int{9, 8, 7}},
			},
			expected: []int{0, 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			winners := findWinners(tt.hands)
			assert.Equal(t, tt.expected, winners)
		})
	}
}
