package holdem

import (
	"fmt"
	"testing"

	"github.com/genewoo/joker/internal/deck"
	"github.com/stretchr/testify/assert"
)

func TestWinningCalculator(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		players     [][]*deck.Card
		simulations int
		limitCount  int
		expected    []float64
	}{
		//{"♠", "♥", "♦", "♣"}
		{
			name: "Two players - 64 vs 36 winner Player 1",
			players: [][]*deck.Card{
				{deck.NewCard("A", "♠"), deck.NewCard("K", "♠")}, // Player 1
				{deck.NewCard("2", "♥"), deck.NewCard("3", "♥")}, // Player 2
			},
			simulations: 10000,
			expected:    []float64{0.64, 0.36, 0.005}, // Player 1 should win most of the time
		},
		{
			name: "Two players - Tie",
			players: [][]*deck.Card{
				{deck.NewCard("A", "♠"), deck.NewCard("A", "♦")}, // Player 1
				{deck.NewCard("A", "♥"), deck.NewCard("A", "♣")}, // Player 2
			},
			simulations: 10000,
			expected:    []float64{0.02, 0.02, 0.95}, // Player 1,2 should be tir most of the time
		},
		{
			name: "Two players AKs vs 23o - 36 vs 64 winner Player 2",
			players: [][]*deck.Card{
				{deck.NewCard("2", "♥"), deck.NewCard("3", "♦")}, // Player 1
				{deck.NewCard("A", "♠"), deck.NewCard("K", "♠")}, // Player 2
			},
			simulations: 10000,
			expected:    []float64{0.32, 0.68, 0.006}, // Player 1 should lost most of the time
		},
		{
			name: "Two players AKs vs 23s - 36 vs 64 winner Player 2",
			players: [][]*deck.Card{
				{deck.NewCard("2", "♥"), deck.NewCard("3", "♥")}, // Player 1
				{deck.NewCard("A", "♠"), deck.NewCard("K", "♠")}, // Player 2
			},
			simulations: 10000,
			expected:    []float64{0.36, 0.64, 0.005}, // Player 1 should lost most of the time
		},
		{
			name: "Two players AKs vs Qs - 54 vs 46 winner Player 2",
			players: [][]*deck.Card{
				{deck.NewCard("Q", "♥"), deck.NewCard("Q", "♦")}, // Player 1
				{deck.NewCard("A", "♠"), deck.NewCard("K", "♠")}, // Player 2
			},
			simulations: 10000,
			expected:    []float64{0.536, 0.46, 0.004}, // Player 1 should lost most of the time
		},
		{
			name: "Three players AKs Even",
			players: [][]*deck.Card{
				{deck.NewCard("A", "♠"), deck.NewCard("K", "♠")}, // Player 1
				{deck.NewCard("A", "♦"), deck.NewCard("K", "♦")}, // Player 2
				{deck.NewCard("A", "♥"), deck.NewCard("K", "♥")}, // Player 3
			},
			simulations: 10000,
			expected:    []float64{0.08, 0.08, 0.08, 0.76}, // Player 1 should lost most of the time
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calc := NewWinningCalculator(tt.players, tt.simulations, NewDefaultHandRanker())
			calcSmart := NewWinningCalculator(tt.players, tt.simulations, NewSmartHandRanker())
			probs := calc.CalculateWinProbabilities()
			probsSmart := calcSmart.CalculateWinProbabilities()
			fmt.Println(probs)
			fmt.Println(probsSmart)
			assert.NotNil(t, probs)
			assert.NotNil(t, probsSmart)
			assert.Equal(t, len(tt.expected), len(probs))
			for i := range tt.expected {
				assert.InEpsilon(t, tt.expected[i], probs[i], 0.5)
				assert.InEpsilon(t, tt.expected[i], probsSmart[i], 0.5)
				assert.InEpsilon(t, probsSmart[i], probs[i], 1)
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
			winners := FindWinners(tt.hands)
			assert.Equal(t, tt.expected, winners)
		})
	}
}

func BenchmarkWinningCalculator(b *testing.B) {
	players := [][]*deck.Card{
		{deck.NewCard("A", "♠"), deck.NewCard("K", "♠")},
		{deck.NewCard("Q", "♥"), deck.NewCard("Q", "♦")},
		{deck.NewCard("A", "♦"), deck.NewCard("K", "♦")},
	}

	b.Run("DefaultHandRanker", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			calc := NewWinningCalculator(players, 1000, NewDefaultHandRanker())
			_ = calc.CalculateWinProbabilities()
		}
	})

	b.Run("SmartHandRanker", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			calc := NewWinningCalculator(players, 1000, NewSmartHandRanker())
			_ = calc.CalculateWinProbabilities()
		}
	})
}

func TestWinningCalculatorWithCommunityCards(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		players        [][]*deck.Card
		communityCards []*deck.Card
		simulations    int
		expected       []float64
	}{
		{
			name: "Two players with flop - AK vs QQ",
			players: [][]*deck.Card{
				{deck.NewCard("A", "♠"), deck.NewCard("K", "♠")}, // Player 1
				{deck.NewCard("Q", "♥"), deck.NewCard("Q", "♦")}, // Player 2
			},
			communityCards: []*deck.Card{
				deck.NewCard("Q", "♣"), deck.NewCard("K", "♥"), deck.NewCard("2", "♦"),
			},
			simulations: 10000,
			expected:    []float64{0.032, 0.968, 0.0}, // Player 2 dominates with three queens
		},
		{
			name: "Two players with turn - Flush draw vs pair",
			players: [][]*deck.Card{
				{deck.NewCard("A", "♥"), deck.NewCard("K", "♥")}, // Player 1 (flush draw)
				{deck.NewCard("J", "♠"), deck.NewCard("J", "♦")}, // Player 2 (pair)
			},
			communityCards: []*deck.Card{
				deck.NewCard("2", "♥"), deck.NewCard("5", "♥"), deck.NewCard("8", "♣"), deck.NewCard("T", "♠"),
			},
			simulations: 10000,
			expected:    []float64{0.333, 0.667, 0.0}, // Player 1 has 9 outs for flush
		},
		{
			name: "Three players with flop - all drawing",
			players: [][]*deck.Card{
				{deck.NewCard("A", "♥"), deck.NewCard("K", "♥")}, // Player 1 (flush draw)
				{deck.NewCard("J", "♠"), deck.NewCard("T", "♠")}, // Player 2 (straight draw)
				{deck.NewCard("9", "♦"), deck.NewCard("9", "♣")}, // Player 3 (pair)
			},
			communityCards: []*deck.Card{
				deck.NewCard("Q", "♥"), deck.NewCard("8", "♥"), deck.NewCard("7", "♠"),
			},
			simulations: 10000,
			expected:    []float64{0.51515, 0.1030303, 0.337374, 0.00}, // Flush draw heavily favored
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calc := NewWinningCalculator(tt.players, tt.simulations, NewDefaultHandRanker(), tt.communityCards...)
			calcSmart := NewWinningCalculator(tt.players, tt.simulations, NewSmartHandRanker(), tt.communityCards...)

			probs := calc.CalculateWinProbabilities()
			probsSmart := calcSmart.CalculateWinProbabilities()

			assert.NotNil(t, probs)
			assert.NotNil(t, probsSmart)
			assert.Equal(t, len(tt.expected), len(probs))

			// Print actual probabilities for debugging
			fmt.Printf("\nTest: %s\n", tt.name)
			fmt.Printf("Default ranker probabilities: %v\n", probs)
			fmt.Printf("Smart ranker probabilities: %v\n", probsSmart)

			// Verify probabilities sum to approximately 1
			sum := 0.0
			for _, p := range probs {
				sum += p
			}
			assert.InDelta(t, 1.0, sum, 0.05, "Probabilities should sum to 1")

			for i := range tt.expected {
				assert.InDelta(t, tt.expected[i], probs[i], 0.05, "Default ranker probability mismatch at index %d", i)
				assert.InDelta(t, tt.expected[i], probsSmart[i], 0.05, "Smart ranker probability mismatch at index %d", i)
				assert.InDelta(t, probsSmart[i], probs[i], 0.02, "Rankers should produce similar results")
			}
		})
	}
}

func TestAppendCommunityCards(t *testing.T) {
	players := [][]*deck.Card{
		{deck.NewCard("A", "♠"), deck.NewCard("K", "♠")},
		{deck.NewCard("Q", "♥"), deck.NewCard("Q", "♦")},
	}

	tests := []struct {
		name          string
		initialCards  []*deck.Card
		cardsToAdd    []*deck.Card
		expectedError bool
		expectedCount int
	}{
		{
			name:         "Add flop to empty",
			initialCards: nil,
			cardsToAdd: []*deck.Card{
				deck.NewCard("2", "♥"), deck.NewCard("3", "♥"), deck.NewCard("4", "♥"),
			},
			expectedError: false,
			expectedCount: 3,
		},
		{
			name: "Add turn after flop",
			initialCards: []*deck.Card{
				deck.NewCard("2", "♥"), deck.NewCard("3", "♥"), deck.NewCard("4", "♥"),
			},
			cardsToAdd: []*deck.Card{
				deck.NewCard("5", "♥"),
			},
			expectedError: false,
			expectedCount: 4,
		},
		{
			name: "Add river after turn",
			initialCards: []*deck.Card{
				deck.NewCard("2", "♥"), deck.NewCard("3", "♥"), deck.NewCard("4", "♥"),
				deck.NewCard("5", "♥"),
			},
			cardsToAdd: []*deck.Card{
				deck.NewCard("6", "♥"),
			},
			expectedError: false,
			expectedCount: 5,
		},
		{
			name: "Try to add beyond river",
			initialCards: []*deck.Card{
				deck.NewCard("2", "♥"), deck.NewCard("3", "♥"), deck.NewCard("4", "♥"),
				deck.NewCard("5", "♥"), deck.NewCard("6", "♥"),
			},
			cardsToAdd: []*deck.Card{
				deck.NewCard("7", "♥"),
			},
			expectedError: true,
			expectedCount: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calc := NewWinningCalculator(players, 1000, NewDefaultHandRanker(), tt.initialCards...)
			err := calc.AppendCommunityCards(tt.cardsToAdd...)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expectedCount, len(calc.communityCards))
		})
	}
}

func TestEvaluateShowdown(t *testing.T) {
	tests := []struct {
		name            string
		players         [][]*deck.Card
		communityCards  []*deck.Card
		expectedWinners []int
		expectError     bool
	}{
		{
			name: "Complete board - one winner",
			players: [][]*deck.Card{
				{deck.NewCard("A", "♠"), deck.NewCard("K", "♠")},
				{deck.NewCard("Q", "♥"), deck.NewCard("Q", "♦")},
			},
			communityCards: []*deck.Card{
				deck.NewCard("A", "♥"), deck.NewCard("A", "♦"), deck.NewCard("K", "♥"),
				deck.NewCard("2", "♣"), deck.NewCard("3", "♦"),
			},
			expectedWinners: []int{0}, // Player 1 wins with two pair aces and kings
			expectError:     false,
		},
		{
			name: "Incomplete board - error",
			players: [][]*deck.Card{
				{deck.NewCard("A", "♠"), deck.NewCard("K", "♠")},
				{deck.NewCard("Q", "♥"), deck.NewCard("Q", "♦")},
			},
			communityCards: []*deck.Card{
				deck.NewCard("A", "♥"), deck.NewCard("A", "♦"), deck.NewCard("K", "♥"),
			},
			expectedWinners: nil,
			expectError:     true,
		},
		{
			name: "Complete board - tie",
			players: [][]*deck.Card{
				{deck.NewCard("A", "♠"), deck.NewCard("K", "♠")},
				{deck.NewCard("A", "♦"), deck.NewCard("K", "♦")},
			},
			communityCards: []*deck.Card{
				deck.NewCard("2", "♥"), deck.NewCard("3", "♥"), deck.NewCard("4", "♥"),
				deck.NewCard("5", "♣"), deck.NewCard("6", "♦"),
			},
			expectedWinners: []int{0, 1}, // Both players tie with ace-king high
			expectError:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calc := NewWinningCalculator(tt.players, 1000, NewDefaultHandRanker(), tt.communityCards...)
			result, err := calc.EvaluateShowdown()

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedWinners, result.Winners)
				assert.Equal(t, len(tt.players), len(result.HandStrengths))
				assert.Equal(t, len(tt.players), len(result.BestHands))

				// Verify each best hand has exactly 5 cards
				for _, hand := range result.BestHands {
					assert.Equal(t, 5, len(hand))
				}
			}
		})
	}
}
