package holdem

import (
	"fmt"

	"github.com/genewoo/joker/internal/deck"
)

var organizer = deck.DefaultOrganizer{}

// RankHand evaluates the best 5-card hand from a player's 2 cards and 5 community cards
func RankHand(playerCards []*deck.Card, communityCards []*deck.Card) (HandRank, []*deck.Card) {
	if len(playerCards) != 2 || len(communityCards) != 5 {
		return InvalidHand, nil
	}

	// Combine and sort all cards
	allCards := append(playerCards, communityCards...)
	organizer.Sort(allCards)

	// Generate and evaluate all possible 5-card combinations
	return evaluateAllCombinations(allCards)
}

// evaluateAllCombinations generates all possible 5-card combinations from the given cards
// and returns the highest ranking hand found
func evaluateAllCombinations(cards []*deck.Card) (HandRank, []*deck.Card) {
	var bestRank HandRank = HighCard
	var bestHand []*deck.Card

	// We need to choose 5 cards from len(cards) (typically 7)
	// Using nested loops to generate all combinations without repetition
	for first := 0; first < len(cards); first++ {
		for second := first + 1; second < len(cards); second++ {
			for third := second + 1; third < len(cards); third++ {
				for fourth := third + 1; fourth < len(cards); fourth++ {
					for fifth := fourth + 1; fifth < len(cards); fifth++ {
						currentHand := []*deck.Card{
							cards[first],
							cards[second],
							cards[third],
							cards[fourth],
							cards[fifth],
						}
						currentRank := evaluateHand(currentHand)
						if currentRank > bestRank {
							bestRank = currentRank
							bestHand = currentHand
						}
					}
				}
			}
		}
	}
	return bestRank, bestHand

}

// evaluateHand determines the rank of a given 5-card hand
func evaluateHand(hand []*deck.Card) HandRank {
	// Create frequency maps for values and suits
	valueCount := make(map[string]int)
	suitCount := make(map[string]int)
	rankBits := 0

	// Map card values to numeric ranks
	valueToRank := map[string]int{
		"2": 2, "3": 3, "4": 4, "5": 5,
		"6": 6, "7": 7, "8": 8, "9": 9,
		"10": 10, "J": 11, "Q": 12,
		"K": 13, "A": 14,
	}

	// Populate value and suite maps and build a rank bits
	for _, card := range hand {
		valueCount[card.Value]++
		suitCount[card.Suit]++
		if rank, ok := valueToRank[card.Value]; ok {
			rankBits |= 1 << uint(rank)
			// // speical handling for A which could make a strage both A-5 10-A
			if rank == 14 {
				rankBits |= 1 << 1
			}
		}
	}

	fmt.Printf("Rank Bits (binary): %b\n", rankBits)

	// Check for flush (all cards same suit)
	flush := false
	for _, count := range suitCount {
		if count == 5 {
			flush = true
			break
		}
	}
	// Print out cards of hand
	for _, card := range hand {
		fmt.Printf("%s%s ", card.Value, card.Suit)
	}
	fmt.Println()
	// Check for straight using bitwise operations
	straight := false
	lowestRank := 0
	// find the rank in the straight, from the lowest to the highest.
	for i := 0; i <= 15-5; i++ {
		mask := 0b11111 << uint16(i)
		fmt.Printf("mask Bits (binary): %d %b\n", i, mask)
		if rankBits&mask == mask {
			straight = true
			lowestRank = i
			break
		}
	}

	// Check for royal flush
	if flush && straight && lowestRank == 10 {
		return RoyalFlush
	}

	// Check for straight flush
	if flush && straight {
		return StraightFlush
	}

	// Check for four of a kind
	for _, count := range valueCount {
		if count == 4 {
			return FourOfAKind
		}
	}

	// Check for full house
	hasThree := false
	hasTwo := false
	for _, count := range valueCount {
		if count == 3 {
			hasThree = true
		} else if count == 2 {
			hasTwo = true
		}
	}
	if hasThree && hasTwo {
		return FullHouse
	}

	// Check for flush
	if flush {
		return Flush
	}

	// Check for straight
	if straight {
		return Straight
	}

	// Check for three of a kind
	if hasThree {
		return ThreeOfAKind
	}

	// Check for two pair
	pairCount := 0
	for _, count := range valueCount {
		if count == 2 {
			pairCount++
		}
	}
	if pairCount >= 2 {
		return TwoPair
	}

	// Check for one pair
	if pairCount == 1 {
		return OnePair
	}

	// Default to high card
	return HighCard
}

// HandRank represents the strength of a poker hand
type HandRank int

const (
	InvalidHand HandRank = iota
	HighCard
	OnePair
	TwoPair
	ThreeOfAKind
	Straight
	Flush
	FullHouse
	FourOfAKind
	StraightFlush
	RoyalFlush
)

// String returns a human-readable representation of the hand rank
func (hr HandRank) String() string {
	return [...]string{
		"Invalid Hand",
		"High Card",
		"One Pair",
		"Two Pair",
		"Three of a Kind",
		"Straight",
		"Flush",
		"Full House",
		"Four of a Kind",
		"Straight Flush",
		"Royal Flush",
	}[hr]
}
