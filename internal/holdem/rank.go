package holdem

import (
	"sort"

	"github.com/genewoo/joker/internal/deck"
)

// Types and Constants
// ==================

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

// HandRanker Interface
// ===================

// HandRanker defines the interface for ranking poker hands
type HandRanker interface {
	// RankHand evaluates the best 5-card hand from a player's 2 cards and 5 community cards
	RankHand(playerCards []*deck.Card, communityCards []*deck.Card) (HandStrength, []*deck.Card)
}

// HandRanker Implementations
// =========================

// DefaultHandRanker implements HandRanker using the traditional all-combinations approach
type DefaultHandRanker struct {
	organizer deck.Organizer
}

// SmartHandRanker implements HandRanker using a more efficient algorithm
type SmartHandRanker struct {
	organizer deck.Organizer
}

// Factory Functions
// ================

// NewDefaultHandRanker creates a new DefaultHandRanker
func NewDefaultHandRanker() *DefaultHandRanker {
	return &DefaultHandRanker{
		organizer: &deck.DefaultOrganizer{},
	}
}

// NewSmartHandRanker creates a new SmartHandRanker
func NewSmartHandRanker() *SmartHandRanker {
	return &SmartHandRanker{
		organizer: &deck.DefaultOrganizer{},
	}
}

// HandStrength Methods
// ===================

// NewHandStrength creates and initializes a new HandStrength
func NewHandStrength() HandStrength {
	return HandStrength{
		Rank:   HighCard,
		Values: make([]int, 0, 5),
	}
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

// HandRank Methods
// ===============

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

// DefaultHandRanker Methods
// ========================

// RankHand evaluates the best 5-card hand from a player's 2 cards and 5 community cards
func (r *DefaultHandRanker) RankHand(playerCards []*deck.Card, communityCards []*deck.Card) (HandStrength, []*deck.Card) {
	if len(playerCards) != 2 || len(communityCards) != 5 {
		strength := NewHandStrength()
		strength.Rank = InvalidHand
		return strength, nil
	}

	// Combine and sort all cards
	allCards := append(playerCards, communityCards...)
	r.organizer.Sort(allCards)

	return evaluateAllCombinations(allCards)
}

// SmartHandRanker Methods
// ======================

// RankHand evaluates the best 5-card hand from a player's 2 cards and 5 community cards
func (r *SmartHandRanker) RankHand(playerCards []*deck.Card, communityCards []*deck.Card) (HandStrength, []*deck.Card) {
	if len(playerCards) != 2 || len(communityCards) != 5 {
		strength := NewHandStrength()
		strength.Rank = InvalidHand
		return strength, nil
	}

	// Combine and sort all cards
	allCards := append(playerCards, communityCards...)
	r.organizer.Sort(allCards)

	// Build analysis maps
	valueCount, _, rankBits, suitedCards := buildHandAnalysis(allCards)

	strength := NewHandStrength()
	var bestHand []*deck.Card

	// Check for flush
	flushCards, flushSuit := findFlush(suitedCards)

	// Check for straight
	straight, lowestRank := findStraight(rankBits)

	// Check for straight flush or royal flush
	if flushCards != nil && straight {
		// Verify straight flush by checking only flush cards
		flushRankBits := 0
		for _, card := range allCards {
			if card.Suit == flushSuit {
				rank := valueToRank[card.Value]
				flushRankBits |= 1 << uint(rank)
				if rank == 14 { // Ace can be low
					flushRankBits |= 1 << 1
				}
			}
		}

		straightFound, straightLowestRank := findStraight(flushRankBits)
		if straightFound {
			bestHand = collectStraightCards(allCards, straightLowestRank, flushSuit)
			if straightLowestRank == 10 {
				strength.Rank = RoyalFlush
				strength.Values = []int{14}
			} else {
				strength.Rank = StraightFlush
				// For A-5 straight flush, the highest card is 5
				if straightLowestRank == 1 {
					strength.Values = []int{5}
				} else {
					strength.Values = []int{straightLowestRank + 4} // Highest card in straight
				}
			}
			return strength, bestHand
		}
	}

	// Check for four of a kind
	if value, kicker, found := findFourOfAKind(valueCount, allCards); found {
		strength.Rank = FourOfAKind
		strength.Values = []int{valueToRank[value], kicker}
		bestHand = collectFourOfAKind(allCards, value)
		return strength, bestHand
	}

	// Check for full house
	if threeValue, pairValue, found := findFullHouse(valueCount); found {
		strength.Rank = FullHouse
		strength.Values = []int{valueToRank[threeValue], valueToRank[pairValue]}
		bestHand = collectFullHouse(allCards, threeValue, pairValue)
		return strength, bestHand
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

// Helper Functions
// ===============

// Common value to rank mapping used across functions
var valueToRank = map[string]int{
	"2": 2, "3": 3, "4": 4, "5": 5, "6": 6, "7": 7, "8": 8, "9": 9,
	"10": 10, "J": 11, "Q": 12, "K": 13, "A": 14,
}

// Analysis Helper Functions
// =======================

func buildHandAnalysis(cards []*deck.Card) (map[string]int, map[string]int, int, map[string][]*deck.Card) {
	valueCount := make(map[string]int)
	suitCount := make(map[string]int)
	suitedCards := make(map[string][]*deck.Card)
	rankBits := 0

	for _, card := range cards {
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

	return valueCount, suitCount, rankBits, suitedCards
}

func findStraight(rankBits int) (bool, int) {
	for i := 0; i <= 15-5; i++ {
		if rankBits&(0b11111<<uint16(i)) == 0b11111<<uint16(i) {
			return true, i
		}
	}
	return false, 0
}

func findFlush(suitedCards map[string][]*deck.Card) ([]*deck.Card, string) {
	for suit, cards := range suitedCards {
		if len(cards) >= 5 {
			return cards[:5], suit
		}
	}
	return nil, ""
}

func findFourOfAKind(valueCount map[string]int, cards []*deck.Card) (string, int, bool) {
	for value, count := range valueCount {
		if count == 4 {
			// Find highest kicker
			var kicker int
			for _, card := range cards {
				if card.Value != value {
					kicker = valueToRank[card.Value]
					break
				}
			}
			return value, kicker, true
		}
	}
	return "", 0, false
}

func findFullHouse(valueCount map[string]int) (string, string, bool) {
	var threeValue string
	var pairValue string

	// Find highest three of a kind
	for value, count := range valueCount {
		if count == 3 {
			if threeValue == "" || valueToRank[value] > valueToRank[threeValue] {
				threeValue = value
			}
		}
	}

	if threeValue != "" {
		// Find highest pair that's not the three of a kind
		for value, count := range valueCount {
			if count >= 2 && value != threeValue {
				if pairValue == "" || valueToRank[value] > valueToRank[pairValue] {
					pairValue = value
				}
			}
		}
	}

	return threeValue, pairValue, threeValue != "" && pairValue != ""
}

// Card Collection Helper Functions
// ==============================

func collectStraightCards(cards []*deck.Card, lowestRank int, requiredSuit string) []*deck.Card {
	result := make([]*deck.Card, 0, 5)

	// First pass: collect cards in ascending order
	for i := lowestRank; i < lowestRank+5; i++ {
		found := false
		for _, card := range cards {
			rank := valueToRank[card.Value]
			if rank == i && (requiredSuit == "" || card.Suit == requiredSuit) {
				result = append(result, card)
				found = true
				break
			}
		}
		if !found {
			// Special case for Ace when it's used as 1 in A-5 straight
			if i == 1 {
				for _, card := range cards {
					if card.Value == "A" && (requiredSuit == "" || card.Suit == requiredSuit) {
						result = append(result, card)
						found = true
						break
					}
				}
			}
			if !found {
				// If we can't find a card for this rank, something is wrong
				// This shouldn't happen since we already verified the straight exists
				return result
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

func collectThreeOfAKind(cards []*deck.Card, value string) []*deck.Card {
	result := make([]*deck.Card, 0, 5)
	kickers := make([]*deck.Card, 0, 2)

	// Collect three of a kind first
	for _, card := range cards {
		if card.Value == value && len(result) < 3 {
			result = append(result, card)
		} else if len(kickers) < 2 {
			kickers = append(kickers, card)
		}
	}

	return append(result, kickers...)
}

func collectTwoPair(cards []*deck.Card, firstPair, secondPair string) []*deck.Card {
	result := make([]*deck.Card, 0, 5)
	firstCount, secondCount := 0, 0
	var kicker *deck.Card

	for _, card := range cards {
		if card.Value == firstPair && firstCount < 2 {
			result = append(result, card)
			firstCount++
		} else if card.Value == secondPair && secondCount < 2 {
			result = append(result, card)
			secondCount++
		} else if kicker == nil {
			kicker = card
		}
	}

	return append(result, kicker)
}

func collectOnePair(cards []*deck.Card, value string) []*deck.Card {
	result := make([]*deck.Card, 0, 5)
	pairCount := 0
	kickers := make([]*deck.Card, 0, 3)

	for _, card := range cards {
		if card.Value == value && pairCount < 2 {
			result = append(result, card)
			pairCount++
		} else if len(kickers) < 3 {
			kickers = append(kickers, card)
		}
	}

	return append(result, kickers...)
}

func collectHighCard(cards []*deck.Card) []*deck.Card {
	result := make([]*deck.Card, 0, 5)
	for i := 0; i < 5 && i < len(cards); i++ {
		result = append(result, cards[i])
	}
	return result
}

// Evaluation Helper Functions
// =========================

func compareHands(hand1, hand2 HandStrength) int {
	return hand1.Compare(hand2)
}

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

func evaluateHand(hand []*deck.Card) HandStrength {
	// Build analysis maps
	valueCount, suitCount, rankBits, _ := buildHandAnalysis(hand)

	// Check for flush (all cards same suit)
	flush := false
	for _, count := range suitCount {
		if count == 5 {
			flush = true
			break
		}
	}

	// Check for straight
	straight, lowestRank := findStraight(rankBits)

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
	if value, kicker, found := findFourOfAKind(valueCount, hand); found {
		strength.Rank = FourOfAKind
		strength.Values = []int{valueToRank[value], kicker}
		return strength
	}

	// Check for full house
	if threeValue, pairValue, found := findFullHouse(valueCount); found {
		strength.Rank = FullHouse
		strength.Values = []int{valueToRank[threeValue], valueToRank[pairValue]}
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
	for value, count := range valueCount {
		if count == 3 {
			strength.Rank = ThreeOfAKind
			// Get kickers by filtering out the three of a kind value
			kickers := make([]int, 0, 2)
			for _, v := range originalValues {
				if v != valueToRank[value] {
					kickers = append(kickers, v)
				}
			}
			// Set values to three of a kind value followed by kickers
			strength.Values = append([]int{valueToRank[value]}, kickers...)
			return strength
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
				break
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

func evaluateRemainingCombinations(cards []*deck.Card, valueCount map[string]int, valueToRank map[string]int) (HandStrength, []*deck.Card) {
	strength := NewHandStrength()
	var bestHand []*deck.Card

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
			bestHand = collectThreeOfAKind(cards, value)
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
		for _, card := range cards {
			if card.Value != pairValues[0] && card.Value != pairValues[1] {
				kicker = valueToRank[card.Value]
				break
			}
		}
		// Set values to higher pair, lower pair, then kicker
		strength.Values = []int{
			valueToRank[pairValues[0]],
			valueToRank[pairValues[1]],
			kicker,
		}
		bestHand = collectTwoPair(cards, pairValues[0], pairValues[1])
		return strength, bestHand
	}

	// Check for one pair
	if pairCount == 1 {
		strength.Rank = OnePair
		// Set values to pair value followed by kickers
		pairValue := valueToRank[pairValues[0]]
		kickers := make([]int, 0, 3)
		for _, card := range cards {
			if card.Value != pairValues[0] {
				kickers = append(kickers, valueToRank[card.Value])
			}
		}
		strength.Values = append([]int{pairValue}, kickers...)
		bestHand = collectOnePair(cards, pairValues[0])
		return strength, bestHand
	}

	// Default to high card
	strength.Rank = HighCard
	// Convert map values to sorted slice for high card values
	values := make([]int, 0, 5)
	for _, card := range cards {
		values = append(values, valueToRank[card.Value])
	}
	// Sort in descending order
	sort.Sort(sort.Reverse(sort.IntSlice(values)))
	strength.Values = values[:5]
	bestHand = collectHighCard(cards)
	return strength, bestHand
}

func getCardRank(value string) int {
	return valueToRank[value]
}
