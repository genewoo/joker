package holdem

import (
	"math/rand"
	"runtime"
	"sync"
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

// Hardcode flag to disable goroutines for debugging
var disableGoroutines = false

// CalculateWinProbabilities calculates winning probabilities for each player
func (wc *WinningCalculator) CalculateWinProbabilities() []float64 {
	if len(wc.players) == 0 {
		return nil
	}

	// Initialize results with a mutex for concurrent access
	var mu sync.Mutex
	results := make([]float64, len(wc.players))
	var tieCount float64

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

	if disableGoroutines {
		// Run simulations sequentially
		localResults := make([]float64, len(wc.players))
		for j := 0; j < len(communityCardHands); j++ {
			communityCards := communityCardHands[j]
			bestHands := make([]HandStrength, len(wc.players))
			for k, hand := range wc.players {
				// var x []*deck.Card
				bestHands[k], _ = RankHand(hand, communityCards.Cards)
				// fmt.Println(x)
			}

			winners := findWinners(bestHands)
			if len(winners) == 1 {
				localResults[winners[0]] += 1.0
			} else if len(winners) > 1 && len(winners) < len(wc.players) {
				winnerPercentage := 1.0 / float64(len(winners))
				for _, winner := range winners {
					localResults[winner] += winnerPercentage
				}
			} else if len(winners) == len(wc.players) {
				tieCount += 1.0
			}
		}
		for k := range results {
			results[k] += localResults[k]
		}
	} else {
		// Original goroutine-based implementation
		var wg sync.WaitGroup
		chunkSize := wc.simulations / runtime.NumCPU()

		for i := 0; i < wc.simulations; i += chunkSize {
			wg.Add(1)
			go func(start, end int) {
				defer wg.Done()
				localResults := make([]float64, len(wc.players))

				for j := start; j < end && j < len(communityCardHands); j++ {
					communityCards := communityCardHands[j]
					bestHands := make([]HandStrength, len(wc.players))
					for k, hand := range wc.players {
						bestHands[k], _ = RankHand(hand, communityCards.Cards)
					}

					winners := findWinners(bestHands)
					if len(winners) == 1 {
						localResults[winners[0]] += 1.0
					} else if len(winners) > 1 && len(winners) < len(wc.players) {
						winnerPercentage := 1.0 / float64(len(winners))
						for _, winner := range winners {
							localResults[winner] += winnerPercentage
						}
					} else if len(winners) == len(wc.players) {
						tieCount += 1.0
					}
				}

				mu.Lock()
				for k := range results {
					results[k] += localResults[k]
				}
				mu.Unlock()
			}(i, i+chunkSize)
		}
		wg.Wait()
	}

	// Calculate probabilities
	probabilities := make([]float64, len(wc.players)+1) // Add extra slot for tie percentage
	for i, wins := range results {
		probabilities[i] = float64(wins) / float64(wc.simulations)
	}
	probabilities[len(wc.players)] = tieCount / float64(wc.simulations) // Add tie probability

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
