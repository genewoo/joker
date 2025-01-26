package holdem

import (
	"math/rand"
	"time"

	"github.com/genewoo/joker/internal/deck"
)

// WinningCalculator calculates winning probabilities for Texas Hold'em hands
type WinningCalculator struct {
	simulations int
	players     [][]*deck.Card
	rng         *rand.Rand
}

// NewWinningCalculator creates a new WinningCalculator with specified players and simulations
func NewWinningCalculator(players [][]*deck.Card, simulations int) *WinningCalculator {
	return &WinningCalculator{
		simulations: simulations,
		players:     players,
		rng:         rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// CalculateWinProbabilities calculates winning probabilities for each player
func (wc *WinningCalculator) CalculateWinProbabilities() []float64 {
	if len(wc.players) == 0 {
		return nil
	}

	// Initialize results
	results := make([]float64, len(wc.players))

	// Create a new deck with player's cards marked
	var markedCardsMasks []string
	for _, hand := range wc.players {
		for _, card := range hand {
			markedCardsMasks = append(markedCardsMasks, card.Value+card.Suit)
		}
	}
	d := deck.NewDeck(markedCardsMasks...)
	// Draw community cards
	communityCardHands := d.DrawWithLimitHands(5, wc.simulations)

	for _, communityCards := range communityCardHands {
		// Evaluate each player's hand
		// fmt.Println(communityCards)
		bestHands := make([]HandStrength, len(wc.players))
		for i, hand := range wc.players {
			bestHands[i], _ = RankHand(hand, communityCards.Cards)
		}

		// fmt.Println(bestHands)
		// Determine winner(s)
		winners := findWinners(bestHands)
		// fmt.Println(winners)
		if len(winners) == 1 {
			results[winners[0]] += 1.0
		} else {
			// in case of tie or all tie

			/// TODO: comment tie situation for now
			winnerPercentage := 1 / float64(len(winners))
			for _, winner := range winners {
				results[winner] += winnerPercentage
			}
		}
	}
	// Calculate probabilities
	probabilities := make([]float64, len(wc.players))
	for i, wins := range results {
		probabilities[i] = float64(wins) / float64(wc.simulations)
	}

	return probabilities
}

// findWinners returns indices of players with the best hand(s)
func findWinners(hands []HandStrength) []int {
	if len(hands) == 0 {
		return nil
	}

	winners := []int{0}
	bestHand := hands[0]

	for i := 1; i < len(hands); i++ {
		comparison := hands[i].Compare(bestHand)
		switch {
		case comparison > 0:
			// New best hand
			winners = []int{i}
			bestHand = hands[i]
		case comparison == 0:
			// Tie
			winners = append(winners, i)
		}
	}

	return winners
}
