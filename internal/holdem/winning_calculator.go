package holdem

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"

	"github.com/genewoo/joker/internal/deck"
)

// WinningCalculator calculates winning probabilities for Texas Hold'em hands
// by simulating multiple games with random community cards and evaluating
// the best possible hand for each player.
type WinningCalculator struct {
	simulations    int            // Number of simulations to run
	players        [][]*deck.Card // Each player's hole cards
	communityCards []*deck.Card   // Pre-existing community cards
	rng            *rand.Rand     // Random number generator for simulations
	ranker         HandRanker     // Hand ranking implementation to use
}

// NewWinningCalculator creates a new WinningCalculator with specified players and simulations.
// Parameters:
//   - players: A slice of player hole cards, where each element is a slice of 2 cards
//   - simulations: Number of Monte Carlo simulations to run
//   - ranker: The hand ranking implementation to use for evaluating hands
//   - communityCards: Optional pre-existing community cards (0-5 cards)
func NewWinningCalculator(players [][]*deck.Card, simulations int, ranker HandRanker, communityCards ...*deck.Card) *WinningCalculator {
	if len(communityCards) > 5 {
		communityCards = communityCards[:5] // Allow up to 5 community cards for showdown
	}
	return &WinningCalculator{
		simulations:    simulations,
		players:        players,
		communityCards: communityCards,
		rng:            rand.New(rand.NewSource(time.Now().UnixNano())),
		ranker:         ranker,
	}
}

// Hardcode flag to disable goroutines for debugging
var disableGoroutines = true

// calculateRequiredSimulations determines the number of simulations needed
// based on the number of remaining community cards.
func (wc *WinningCalculator) calculateRequiredSimulations() int {
	remainingCards := 5 - len(wc.communityCards)
	switch remainingCards {
	case 5: // Pre-flop
		return wc.simulations
	case 2: // After flop
		return min(wc.simulations, 990) // C(47,2) = 1,081
	case 1: // After turn
		return min(wc.simulations, 45) // C(46,1) = 46
	case 0: // After river
		return 1
	default:
		return wc.simulations
	}
}

// min returns the smaller of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// CalculateWinProbabilities calculates winning probabilities for each player
// by running Monte Carlo simulations with random community cards.
// Returns a slice of probabilities where:
//   - indices 0 to n-1 contain each player's probability of winning
//   - index n contains the probability of a complete tie between all players
//
// The probabilities sum to 1.0.
func (wc *WinningCalculator) CalculateWinProbabilities() []float64 {
	if len(wc.players) == 0 {
		return nil
	}

	// Initialize results with a mutex for concurrent access
	var mu sync.Mutex
	results := make([]float64, len(wc.players))
	var tieCount float64

	// Create a new deck with player's cards and community cards marked
	var markedCardsMasks []string
	for _, hand := range wc.players {
		for _, card := range hand {
			markedCardsMasks = append(markedCardsMasks, card.Value+card.Suit)
		}
	}
	for _, card := range wc.communityCards {
		markedCardsMasks = append(markedCardsMasks, card.Value+card.Suit)
	}
	d := deck.NewDeck(markedCardsMasks...)

	// Calculate remaining community cards needed and required simulations
	remainingCards := 5 - len(wc.communityCards)
	requiredSimulations := wc.calculateRequiredSimulations()

	// Draw remaining community cards for each simulation
	communityCardHands := d.DrawWithLimitHands(remainingCards, requiredSimulations)

	if disableGoroutines {
		// Run simulations sequentially
		localResults := make([]float64, len(wc.players))
		for j := 0; j < len(communityCardHands); j++ {
			drawnCards := communityCardHands[j]
			// Combine pre-existing community cards with drawn cards
			allCommunityCards := append([]*deck.Card{}, wc.communityCards...)
			allCommunityCards = append(allCommunityCards, drawnCards.Cards...)

			bestHands := make([]HandStrength, len(wc.players))
			for k, hand := range wc.players {
				bestHands[k], _ = wc.ranker.RankHand(Texas, hand, allCommunityCards)
			}

			winners := FindWinners(bestHands)
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
		chunkSize := requiredSimulations / runtime.NumCPU()
		if chunkSize == 0 {
			chunkSize = 1
		}

		for i := 0; i < requiredSimulations; i += chunkSize {
			wg.Add(1)
			go func(start, end int) {
				defer wg.Done()
				localResults := make([]float64, len(wc.players))
				localTieCount := 0.0

				for j := start; j < end && j < len(communityCardHands); j++ {
					drawnCards := communityCardHands[j]
					// Combine pre-existing community cards with drawn cards
					allCommunityCards := append([]*deck.Card{}, wc.communityCards...)
					allCommunityCards = append(allCommunityCards, drawnCards.Cards...)

					bestHands := make([]HandStrength, len(wc.players))
					for k, hand := range wc.players {
						bestHands[k], _ = wc.ranker.RankHand(Texas, hand, allCommunityCards)
					}

					winners := FindWinners(bestHands)
					if len(winners) == 1 {
						localResults[winners[0]] += 1.0
					} else if len(winners) > 1 && len(winners) < len(wc.players) {
						winnerPercentage := 1.0 / float64(len(winners))
						for _, winner := range winners {
							localResults[winner] += winnerPercentage
						}
					} else if len(winners) == len(wc.players) {
						localTieCount += 1.0
					}
				}

				mu.Lock()
				for k := range results {
					results[k] += localResults[k]
				}
				tieCount += localTieCount
				mu.Unlock()
			}(i, i+chunkSize)
		}
		wg.Wait()
	}

	// Calculate probabilities
	probabilities := make([]float64, len(wc.players)+1) // Add extra slot for tie percentage
	totalSimulations := float64(requiredSimulations)    // Use the calculated required simulations
	for i, wins := range results {
		probabilities[i] = float64(wins) / totalSimulations
	}
	probabilities[len(wc.players)] = tieCount / totalSimulations // Add tie probability

	return probabilities
}

// AppendCommunityCards adds additional community cards to the calculator.
// Returns an error if adding the new cards would exceed 5 total community cards.
// This is useful for updating probabilities as more community cards are revealed (e.g., turn and river).
func (wc *WinningCalculator) AppendCommunityCards(cards ...*deck.Card) error {
	if len(wc.communityCards)+len(cards) > 5 {
		return fmt.Errorf("cannot add %d cards: would exceed maximum of 5 community cards (current: %d)",
			len(cards), len(wc.communityCards))
	}

	wc.communityCards = append(wc.communityCards, cards...)
	return nil
}

// ShowdownResult contains the evaluation results for the showdown
type ShowdownResult struct {
	HandStrengths []HandStrength // Hand strength for each player
	BestHands     [][]*deck.Card // The actual 5 cards making up each player's best hand
	Winners       []int          // Indices of winning players
}

// EvaluateShowdown evaluates the final hands when all 5 community cards are available.
// Returns the best hand for each player and the indices of winners.
// Returns an error if there are not exactly 5 community cards.
func (wc *WinningCalculator) EvaluateShowdown() (*ShowdownResult, error) {
	if len(wc.communityCards) != 5 {
		return nil, fmt.Errorf("showdown requires exactly 5 community cards, got %d", len(wc.communityCards))
	}

	handStrengths := make([]HandStrength, len(wc.players))
	bestHands := make([][]*deck.Card, len(wc.players))
	for i, hand := range wc.players {
		var bestHand []*deck.Card
		handStrengths[i], bestHand = wc.ranker.RankHand(Texas, hand, wc.communityCards)
		bestHands[i] = bestHand
	}

	winners := FindWinners(handStrengths)
	return &ShowdownResult{
		HandStrengths: handStrengths,
		BestHands:     bestHands,
		Winners:       winners,
	}, nil
}

// FindWinners returns indices of players with the best hand(s).
// If multiple players have equally strong hands, all their indices are returned.
func FindWinners(hands []HandStrength) []int {
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
