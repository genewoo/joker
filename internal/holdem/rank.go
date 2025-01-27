package holdem

import (
	"sort"

	"github.com/genewoo/joker/internal/deck"
)

var organizer = deck.DefaultOrganizer{}

// RankHand evaluates the best 5-card hand from a player's 2 cards and 5 community cards
func RankHand(playerCards []*deck.Card, communityCards []*deck.Card) (HandStrength, []*deck.Card) {
	if len(playerCards) != 2 || len(communityCards) != 5 {
		strength := NewHandStrength()
		strength.Rank = InvalidHand
		return strength, nil
	}

	// Combine and sort all cards
	allCards := append(playerCards, communityCards...)
	organizer.Sort(allCards)

	// Generate and evaluate all possible 5-card combinations
	return evaluateAllCombinations(allCards)
}

// compareHands returns 1 if hand1 is stronger than hand2
// return -1 if hand2 is stronger than hand1
// returns 0 if equal
func compareHands(hand1, hand2 HandStrength) int {
	// First compare ranks
	if hand1.Rank < hand2.Rank {
		return -1
	}
	if hand1.Rank > hand2.Rank {
		return 1
	}

	// If ranks are equal, compare values element by element
	for i := 0; i < len(hand1.Values) && i < len(hand2.Values); i++ {
		if hand1.Values[i] < hand2.Values[i] {
			return -1
		}
		if hand1.Values[i] > hand2.Values[i] {
			return 1
		}
	}

	// Hands are equal
	return 0
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
						// fmt.Println("CURRENT STRENGTH:", currentStrength)
						// fmt.Println("BEST STRENGTH:", bestStrength)
						if bestHand == nil || compareHands(currentStrength, bestStrength) == 1 {

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
	strength := NewHandStrength()

	// Populate values with card ranks
	for _, card := range hand {
		if rank, ok := valueToRank[card.Value]; ok {
			strength.Values = append(strength.Values, rank)
		}
	}
	// Sort values in descending order
	sort.Sort(sort.Reverse(sort.IntSlice(strength.Values)))

	// Make a copy of the original sorted values for reference
	originalValues := make([]int, len(strength.Values))
	copy(originalValues, strength.Values)

	// Check for royal flush
	if flush && straight && lowestRank == 10 {
		strength.Rank = RoyalFlush
		return strength
	}

	// Check for straight flush
	if flush && straight {
		strength.Rank = StraightFlush
		strength.Values = []int{lowestRank + 4} // Highest card in straight
		return strength
	}

	// Check for four of a kind
	for value, count := range valueCount {
		if count == 4 {
			strength.Rank = FourOfAKind
			// Find the kicker from the original sorted values
			var kicker int
			for _, v := range originalValues {
				if v != valueToRank[value] {
					kicker = v
					break
				}
			}
			// Set values to four of a kind value followed by kicker
			strength.Values = []int{valueToRank[value], kicker}
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
		strength.Values = []int{valueToRank[threeValue], valueToRank[twoValue]}
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
		strength.Values = []int{lowestRank + 4} // Highest card in straight
		return strength
	}

	// Check for three of a kind
	if hasThree {
		strength.Rank = ThreeOfAKind
		// Get all card values
		allValues := make([]int, 0, 5)
		for _, card := range hand {
			allValues = append(allValues, valueToRank[card.Value])
		}
		// Get kickers by filtering out the three of a kind value
		kickers := make([]int, 0, 2)
		for _, v := range allValues {
			if v != valueToRank[threeValue] {
				kickers = append(kickers, v)
			}
		}
		// Set values to three of a kind value followed by kickers
		strength.Values = append([]int{valueToRank[threeValue]}, kickers...)
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

	// Comment out due to personal thought.

	// sort before pairCount
	if pairCount >= 1 {
		// Sort pairs by value
		sort.Slice(pairValues, func(i, j int) bool {
			return valueToRank[pairValues[i]] > valueToRank[pairValues[j]]
		})
	}
	if pairCount >= 2 {
		strength.Rank = TwoPair
		// find the kicker by filter out the two pair values
		kicker := 0
		for _, v := range originalValues {
			if v != valueToRank[pairValues[0]] && v != valueToRank[pairValues[1]] {
				kicker = v
			}
		}
		// Set values to higher pair, lower pair, then kicker
		strength.Values = []int{
			valueToRank[pairValues[0]],
			valueToRank[pairValues[1]],
			kicker,
		}
		return strength
	}

	// Check for one pair
	if pairCount == 1 {
		strength.Rank = OnePair
		// Set values to pair value followed by kickers
		pairValue := valueToRank[pairValues[0]]
		kickers := make([]int, 0, 3)
		for _, v := range originalValues {
			if v != pairValue {
				kickers = append(kickers, v)
			}
		}
		strength.Values = append([]int{pairValue}, kickers...)
		return strength
	}

	// Default to high card
	strength.Rank = HighCard
	strength.Values = originalValues
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
	Rank   HandRank
	Values []int // Card values in descending order of importance
}

// NewHandStrength creates and initializes a new HandStrength
func NewHandStrength() HandStrength {
	return HandStrength{
		Rank:   HighCard,
		Values: make([]int, 0, 5),
	}
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

// Compare compares two HandStrength values
// Returns -1 if h is weaker than other, 0 if equal, 1 if h is stronger
func (h HandStrength) Compare(other HandStrength) int {
	// First compare ranks
	if h.Rank < other.Rank {
		return -1
	}
	if h.Rank > other.Rank {
		return 1
	}

	// If ranks are equal, compare values element by element
	for i := 0; i < len(h.Values) && i < len(other.Values); i++ {
		if h.Values[i] < other.Values[i] {
			return -1
		}
		if h.Values[i] > other.Values[i] {
			return 1
		}
	}

	// Hands are equal
	return 0
}

// SmartRankHand evaluates the best 5-card hand from a player's 2 cards and 5 community cards
// without generating all possible 5-card combinations
func SmartRankHand(playerCards []*deck.Card, communityCards []*deck.Card) (HandStrength, []*deck.Card) {
	if len(playerCards) != 2 || len(communityCards) != 5 {
		strength := NewHandStrength()
		strength.Rank = InvalidHand
		return strength, nil
	}

	// Combine all cards
	allCards := append(playerCards, communityCards...)
	organizer.Sort(allCards)

	// Create frequency maps for values and suits
	valueCount := make(map[string]int)
	suitCount := make(map[string]int)
	rankBits := 0
	valueToRank := map[string]int{
		"2": 2, "3": 3, "4": 4, "5": 5, "6": 6, "7": 7, "8": 8, "9": 9,
		"10": 10, "J": 11, "Q": 12, "K": 13, "A": 14,
	}

	// Track cards by suit for flush detection
	suitedCards := make(map[string][]*deck.Card)

	// Populate maps
	for _, card := range allCards {
		valueCount[card.Value]++
		suitCount[card.Suit]++
		suitedCards[card.Suit] = append(suitedCards[card.Suit], card)
		if rank, ok := valueToRank[card.Value]; ok {
			rankBits |= 1 << uint(rank)
			if rank == 14 { // Ace can be low
				rankBits |= 1 << 1
			}
		}
	}

	strength := NewHandStrength()
	bestHand := make([]*deck.Card, 0, 5)

	// Check for flush first
	var flushCards []*deck.Card
	var flushSuit string
	for suit, cards := range suitedCards {
		if len(cards) >= 5 {
			flushCards = cards[:5] // Take top 5 cards of the same suit
			flushSuit = suit
			break
		}
	}

	// Check for straight in the entire hand
	straight := false
	lowestRank := 0
	for i := 0; i <= 15-5; i++ {
		if rankBits&(0b11111<<uint16(i)) == 0b11111<<uint16(i) {
			straight = true
			lowestRank = i
			break
		}
	}

	// Check for straight flush or royal flush
	if flushCards != nil && straight {
		// Verify straight flush by checking only flush cards
		flushRankBits := 0
		for _, card := range allCards {
			if card.Suit == flushSuit {
				rank := valueToRank[card.Value]
				flushRankBits |= 1 << uint(rank)
				if rank == 14 {
					flushRankBits |= 1 << 1
				}
			}
		}

		for i := 0; i <= 15-5; i++ {
			if flushRankBits&(0b11111<<uint16(i)) == 0b11111<<uint16(i) {
				if i == 10 {
					strength.Rank = RoyalFlush
				} else {
					strength.Rank = StraightFlush
				}
				strength.Values = []int{i + 4}
				// Collect the straight flush cards
				bestHand = collectStraightCards(allCards, i, flushSuit)
				return strength, bestHand
			}
		}
	}

	// Check for four of a kind
	for value, count := range valueCount {
		if count == 4 {
			strength.Rank = FourOfAKind
			rank := valueToRank[value]
			// Find highest kicker
			var kicker int
			for _, card := range allCards {
				if card.Value != value {
					kicker = valueToRank[card.Value]
					break
				}
			}
			strength.Values = []int{rank, kicker}
			bestHand = collectFourOfAKind(allCards, value)
			return strength, bestHand
		}
	}

	// Check for full house
	var threeValue string
	var pairValue string
	for value, count := range valueCount {
		if count == 3 {
			if threeValue == "" || valueToRank[value] > valueToRank[threeValue] {
				threeValue = value
			}
		}
	}
	if threeValue != "" {
		for value, count := range valueCount {
			if count >= 2 && value != threeValue {
				if pairValue == "" || valueToRank[value] > valueToRank[pairValue] {
					pairValue = value
				}
			}
		}
		if pairValue != "" {
			strength.Rank = FullHouse
			strength.Values = []int{valueToRank[threeValue], valueToRank[pairValue]}
			bestHand = collectFullHouse(allCards, threeValue, pairValue)
			return strength, bestHand
		}
	}

	// Handle flush
	if flushCards != nil {
		strength.Rank = Flush
		strength.Values = make([]int, 5)
		for i, card := range flushCards {
			strength.Values[i] = valueToRank[card.Value]
		}
		bestHand = flushCards
		return strength, bestHand
	}

	// Handle straight
	if straight {
		strength.Rank = Straight
		strength.Values = []int{lowestRank + 4}
		bestHand = collectStraightCards(allCards, lowestRank, "")
		return strength, bestHand
	}

	// Rest of the evaluations (three of a kind, two pair, one pair, high card)
	return evaluateRemainingCombinations(allCards, valueCount, valueToRank)
}

// Helper functions for collecting cards for each hand type
func collectStraightCards(cards []*deck.Card, lowestRank int, requiredSuit string) []*deck.Card {
	result := make([]*deck.Card, 0, 5)
	ranks := make(map[int]bool)
	for i := lowestRank; i < lowestRank+5; i++ {
		ranks[i] = true
	}

	for _, card := range cards {
		rank := getCardRank(card.Value)
		if ranks[rank] && (requiredSuit == "" || card.Suit == requiredSuit) {
			result = append(result, card)
			delete(ranks, rank)
			if len(result) == 5 {
				break
			}
		}
	}
	return result
}

func collectFourOfAKind(cards []*deck.Card, value string) []*deck.Card {
	result := make([]*deck.Card, 0, 5)
	var kicker *deck.Card

	// Collect four of a kind cards first
	for _, card := range cards {
		if card.Value == value {
			result = append(result, card)
		} else if kicker == nil {
			kicker = card
		}
	}

	// Add the kicker
	result = append(result, kicker)
	return result
}

func collectFullHouse(cards []*deck.Card, threeValue, pairValue string) []*deck.Card {
	result := make([]*deck.Card, 0, 5)
	threeCount := 0
	pairCount := 0

	for _, card := range cards {
		if card.Value == threeValue && threeCount < 3 {
			result = append(result, card)
			threeCount++
		} else if card.Value == pairValue && pairCount < 2 {
			result = append(result, card)
			pairCount++
		}
	}
	return result
}

// evaluateRemainingCombinations evaluates the remaining combinations for a hand
func evaluateRemainingCombinations(cards []*deck.Card, valueCount map[string]int, valueToRank map[string]int) (HandStrength, []*deck.Card) {
	strength := NewHandStrength()
	bestHand := make([]*deck.Card, 0, 5)

	// Check for three of a kind
	for value, count := range valueCount {
		if count == 3 {
			strength.Rank = ThreeOfAKind
			// Get all card values
			allValues := make([]int, 0, 5)
			for _, card := range cards {
				allValues = append(allValues, valueToRank[card.Value])
			}
			// Get kickers by filtering out the three of a kind value
			kickers := make([]int, 0, 2)
			for _, v := range allValues {
				if v != valueToRank[value] {
					kickers = append(kickers, v)
				}
			}
			// Set values to three of a kind value followed by kickers
			strength.Values = append([]int{valueToRank[value]}, kickers...)
			bestHand = cards
			return strength, bestHand
		}
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

	// Comment out due to personal thought.

	// sort before pairCount
	if pairCount >= 1 {
		// Sort pairs by value
		sort.Slice(pairValues, func(i, j int) bool {
			return valueToRank[pairValues[i]] > valueToRank[pairValues[j]]
		})
	}
	if pairCount >= 2 {
		strength.Rank = TwoPair
		// find the kicker by filter out the two pair values
		kicker := 0
		for _, v := range valueToRank {
			if v != valueToRank[pairValues[0]] && v != valueToRank[pairValues[1]] {
				kicker = v
			}
		}
		// Set values to higher pair, lower pair, then kicker
		strength.Values = []int{
			valueToRank[pairValues[0]],
			valueToRank[pairValues[1]],
			kicker,
		}
		bestHand = cards
		return strength, bestHand
	}

	// Check for one pair
	if pairCount == 1 {
		strength.Rank = OnePair
		// Set values to pair value followed by kickers
		pairValue := valueToRank[pairValues[0]]
		kickers := make([]int, 0, 3)
		for _, v := range valueToRank {
			if v != pairValue {
				kickers = append(kickers, v)
			}
		}
		strength.Values = append([]int{pairValue}, kickers...)
		bestHand = cards
		return strength, bestHand
	}

	// Default to high card
	strength.Rank = HighCard
	// Convert map values to sorted slice for high card values
	values := make([]int, 0, len(valueToRank))
	for _, v := range valueToRank {
		values = append(values, v)
	}
	// Sort in descending order
	sort.Sort(sort.Reverse(sort.IntSlice(values)))
	strength.Values = values
	bestHand = cards
	return strength, bestHand
}

// getCardRank returns the rank of a card
func getCardRank(value string) int {
	valueToRank := map[string]int{
		"2": 2, "3": 3, "4": 4, "5": 5, "6": 6, "7": 7, "8": 8, "9": 9,
		"10": 10, "J": 11, "Q": 12, "K": 13, "A": 14,
	}
	return valueToRank[value]
}
