package guandan

import (
	"sort"

	"github.com/genewoo/joker/internal/deck"
)

// Combiner handles card combination logic
type Combiner struct {
	currentLevel string
	valueToRank  map[string]int
}

// CombinationType represents the type of card combination
type CombinationType int

const (
	InvalidCombination CombinationType = iota
	Single
	Pair
	Triple
	Plate
	Tube
	FullHouse
	Straight
	Bomb
	StraightFlush
	JokerBomb
)

// CombinationStrength contains detailed information about a combination
type CombinationStrength struct {
	Type   CombinationType
	Values []int // Card values in descending order of importance
}

// NewCombiner creates a new Combiner instance
func NewCombiner(currentLevel string) *Combiner {
	return &Combiner{
		currentLevel: currentLevel,
		valueToRank: map[string]int{
			"2": 2, "3": 3, "4": 4, "5": 5,
			"6": 6, "7": 7, "8": 8, "9": 9,
			"10": 10, "J": 11, "Q": 12,
			"K": 13, "A": 14, "Joker": 15,
		},
	}
}

// EvaluateCombination determines the type and strength of a card combination
func (c *Combiner) EvaluateCombination(cards []*deck.Card) CombinationStrength {
	strength := CombinationStrength{
		Type:   InvalidCombination,
		Values: make([]int, 0, len(cards)),
	}

	// Convert card values to ranks
	for _, card := range cards {
		strength.Values = append(strength.Values, c.valueToRank[card.Value])
	}
	sort.Sort(sort.Reverse(sort.IntSlice(strength.Values)))

	// Check combination types in descending order of strength
	switch {
	case c.isJokerBomb(cards):
		strength.Type = JokerBomb
	case c.isStraightFlush(cards):
		strength.Type = StraightFlush
	case c.isBomb(cards):
		strength.Type = Bomb
	case c.isStraight(cards):
		strength.Type = Straight
	case c.isFullHouse(cards):
		strength.Type = FullHouse
	case c.isTube(cards):
		strength.Type = Tube
	case c.isPlate(cards):
		strength.Type = Plate
	case c.isTriple(cards):
		strength.Type = Triple
	case c.isPair(cards):
		strength.Type = Pair
	case c.isSingle(cards):
		strength.Type = Single
	}

	return strength
}

// isSingle checks if cards form a single card combination
func (c *Combiner) isSingle(cards []*deck.Card) bool {
	return len(cards) == 1
}

// isPair checks if cards form a pair combination
func (c *Combiner) isPair(cards []*deck.Card) bool {
	if len(cards) != 2 {
		return false
	}
	return cards[0].Value == cards[1].Value
}

// isTriple checks if cards form a triple combination
func (c *Combiner) isTriple(cards []*deck.Card) bool {
	if len(cards) != 3 {
		return false
	}
	return cards[0].Value == cards[1].Value &&
		cards[1].Value == cards[2].Value
}

// isPlate checks if cards form a plate combination (2 consecutive triples)
func (c *Combiner) isPlate(cards []*deck.Card) bool {
	if len(cards) != 6 {
		return false
	}

	// Split into two triples
	first := cards[:3]
	second := cards[3:]

	if !c.isTriple(first) || !c.isTriple(second) {
		return false
	}

	// Check if consecutive
	values := []string{first[0].Value, second[0].Value}
	sort.Strings(values)
	return c.areConsecutive(values[0], values[1])
}

// isTube checks if cards form a tube combination (3 consecutive pairs)
func (c *Combiner) isTube(cards []*deck.Card) bool {
	if len(cards) != 6 {
		return false
	}

	// Split into three pairs
	pairs := [][]*deck.Card{
		cards[:2],
		cards[2:4],
		cards[4:],
	}

	// Check each pair and collect values
	values := make([]string, 0, 3)
	for _, pair := range pairs {
		if !c.isPair(pair) {
			return false
		}
		values = append(values, pair[0].Value)
	}

	// Check if consecutive
	sort.Strings(values)
	return c.areConsecutive(values[0], values[1]) &&
		c.areConsecutive(values[1], values[2])
}

// isFullHouse checks if cards form a full house combination
func (c *Combiner) isFullHouse(cards []*deck.Card) bool {
	if len(cards) != 5 {
		return false
	}

	// Count card values
	valueCount := make(map[string]int)
	for _, card := range cards {
		valueCount[card.Value]++
	}

	// Should have one triple and one pair
	hasTriple := false
	hasPair := false
	for _, count := range valueCount {
		if count == 3 {
			hasTriple = true
		} else if count == 2 {
			hasPair = true
		}
	}

	return hasTriple && hasPair
}

// isStraight checks if cards form a straight combination
func (c *Combiner) isStraight(cards []*deck.Card) bool {
	if len(cards) != 5 {
		return false
	}

	// Get unique values
	uniqueValues := make([]string, 0, 5)
	seen := make(map[string]bool)
	for _, card := range cards {
		if !seen[card.Value] {
			seen[card.Value] = true
			uniqueValues = append(uniqueValues, card.Value)
		}
	}

	if len(uniqueValues) != 5 {
		return false
	}

	// Sort values
	sort.Strings(uniqueValues)

	// Check if consecutive
	for i := 0; i < 4; i++ {
		if !c.areConsecutive(uniqueValues[i], uniqueValues[i+1]) {
			return false
		}
	}

	return true
}

// isBomb checks if cards form a bomb combination
func (c *Combiner) isBomb(cards []*deck.Card) bool {
	if len(cards) < 4 {
		return false
	}

	// All cards must have same value
	firstValue := cards[0].Value
	for _, card := range cards {
		if card.Value != firstValue {
			return false
		}
	}

	return true
}

// isStraightFlush checks if cards form a straight flush combination
func (c *Combiner) isStraightFlush(cards []*deck.Card) bool {
	if len(cards) != 5 {
		return false
	}

	// First check if straight
	if !c.isStraight(cards) {
		return false
	}

	// Then check if all same suit
	firstSuit := cards[0].Suit
	for _, card := range cards {
		if card.Suit != firstSuit {
			return false
		}
	}

	return true
}

// isJokerBomb checks if cards form a joker bomb combination
func (c *Combiner) isJokerBomb(cards []*deck.Card) bool {
	if len(cards) != 4 {
		return false
	}

	for _, card := range cards {
		if card.Value != "Joker" {
			return false
		}
	}

	return true
}

// areConsecutive checks if two card values are consecutive
// considering the current level
func (c *Combiner) areConsecutive(value1, value2 string) bool {
	order := []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"}

	// Find indices
	index1 := -1
	index2 := -1
	for i, val := range order {
		if val == value1 {
			index1 = i
		}
		if val == value2 {
			index2 = i
		}
	}

	if index1 == -1 || index2 == -1 {
		return false
	}

	// Handle wrap-around for current level
	if value1 == c.currentLevel && value2 == "2" {
		return true
	}
	if value2 == c.currentLevel && value1 == "2" {
		return true
	}

	return abs(index1-index2) == 1
}

// abs returns absolute value of an integer
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
