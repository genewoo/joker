package holdem

import (
	"sort"

	"github.com/genewoo/joker/internal/deck"
)

var organizer = deck.DefaultOrganizer{}

// RankHand evaluates the best 5-card hand from a player's 2 cards and 5 community cards
func RankHand(playerCards []*deck.Card, communityCards []*deck.Card) (HandStrength, []*deck.Card) {
	if len(playerCards) != 2 || len(communityCards) != 5 {
		return HandStrength{Rank: InvalidHand}, nil
	}

	// Combine and sort all cards
	allCards := append(playerCards, communityCards...)
	organizer.Sort(allCards)

	// Generate and evaluate all possible 5-card combinations
	return evaluateAllCombinations(allCards)
}

// compareHands returns true if hand1 is stronger than hand2
func compareHands(hand1, hand2 HandStrength) bool {
	if hand1.Rank != hand2.Rank {
		return hand1.Rank > hand2.Rank
	}

	if hand1.Primary != hand2.Primary {
		return hand1.Primary > hand2.Primary
	}

	if hand1.Secondary != hand2.Secondary {
		return hand1.Secondary > hand2.Secondary
	}

	// Compare kickers
	for i := 0; i < len(hand1.Kickers) && i < len(hand2.Kickers); i++ {
		if hand1.Kickers[i] != hand2.Kickers[i] {
			return hand1.Kickers[i] > hand2.Kickers[i]
		}
	}

	// If we get here, the hands are equal
	return false
}

// evaluateAllCombinations generates all possible 5-card combinations from the given cards
// and returns the highest ranking hand found
func evaluateAllCombinations(cards []*deck.Card) (HandStrength, []*deck.Card) {
	var bestStrength HandStrength
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
						currentStrength := evaluateHand(currentHand)
						if bestHand == nil || compareHands(currentStrength, bestStrength) {
							bestStrength = currentStrength
							bestHand = currentHand
						}
					}
				}
			}
		}
	}
	return bestStrength, bestHand

}

// evaluateHand determines the rank and strength of a given 5-card hand
func evaluateHand(hand []*deck.Card) HandStrength {
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

	// fmt.Printf("Rank Bits (binary): %b\n", rankBits)

	// Check for flush (all cards same suit)
	flush := false
	for _, count := range suitCount {
		if count == 5 {
			flush = true
			break
		}
	}
	// Print out cards of hand
	// for _, card := range hand {
	// 	fmt.Printf("%s%s ", card.Value, card.Suit)
	// }
	// fmt.Println()
	// Check for straight using bitwise operations
	straight := false
	lowestRank := 0
	// find the rank in the straight, from the lowest to the highest.
	for i := 0; i <= 15-5; i++ {
		mask := 0b11111 << uint16(i)
		// fmt.Printf("mask Bits (binary): %d %b\n", i, mask)
		if rankBits&mask == mask {
			straight = true
			lowestRank = i
			break
		}
	}

	// Create HandStrength with basic rank
	strength := HandStrength{
		Rank:    HighCard,
		Kickers: make([]int, 0, 5),
	}

	// Populate kickers with card values in descending order
	for _, card := range hand {
		if rank, ok := valueToRank[card.Value]; ok {
			strength.Kickers = append(strength.Kickers, rank)
		}
	}
	sort.Sort(sort.Reverse(sort.IntSlice(strength.Kickers)))

	// Check for royal flush
	if flush && straight && lowestRank == 10 {
		strength.Rank = RoyalFlush
		return strength
	}

	// Check for straight flush
	if flush && straight {
		strength.Rank = StraightFlush
		strength.Primary = lowestRank + 4 // Highest card in straight
		return strength
	}

	// Check for four of a kind
	for value, count := range valueCount {
		if count == 4 {
			strength.Rank = FourOfAKind
			strength.Primary = valueToRank[value]
			// Remove the four of a kind from kickers
			strength.Kickers = []int{strength.Kickers[4]}
			return strength
		}
	}

	// Check for full house
	hasThree := false
	hasTwo := false
	var threeValue, twoValue string
	for value, count := range valueCount {
		if count == 3 {
			hasThree = true
			threeValue = value
		} else if count == 2 {
			hasTwo = true
			twoValue = value
		}
	}
	if hasThree && hasTwo {
		strength.Rank = FullHouse
		strength.Primary = valueToRank[threeValue]
		strength.Secondary = valueToRank[twoValue]
		return strength
	}

	// Check for flush
	if flush {
		strength.Rank = Flush
		return strength
	}

	// Check for straight
	if straight {
		strength.Rank = Straight
		strength.Primary = lowestRank + 4 // Highest card in straight
		return strength
	}

	// Check for three of a kind
	if hasThree {
		strength.Rank = ThreeOfAKind
		strength.Primary = valueToRank[threeValue]
		// Remove the three of a kind from kickers
		strength.Kickers = strength.Kickers[3:]
		return strength
	}

	// Check for two pair
	pairCount := 0
	var pairValues []string
	for value, count := range valueCount {
		if count == 2 {
			pairCount++
			pairValues = append(pairValues, value)
		}
	}
	if pairCount >= 2 {
		strength.Rank = TwoPair
		// Sort pairs by value
		sort.Slice(pairValues, func(i, j int) bool {
			return valueToRank[pairValues[i]] > valueToRank[pairValues[j]]
		})
		strength.Primary = valueToRank[pairValues[0]]
		strength.Secondary = valueToRank[pairValues[1]]
		// Remove pairs from kickers
		strength.Kickers = []int{strength.Kickers[4]}
		return strength
	}

	// Check for one pair
	if pairCount == 1 {
		strength.Rank = OnePair
		strength.Primary = valueToRank[pairValues[0]]
		// Remove pair from kickers
		strength.Kickers = strength.Kickers[2:]
		return strength
	}

	// Default to high card
	strength.Rank = HighCard
	return strength
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

// HandStrength contains detailed information about a hand's strength
type HandStrength struct {
	Rank      HandRank
	Primary   int   // Highest card for high card/straight, value of set for pairs/trips
	Secondary int   // Second highest card or value of second pair
	Kickers   []int // Remaining card values in descending order
	IsWheel   bool  // Special case for A-2-3-4-5 straight
}

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
