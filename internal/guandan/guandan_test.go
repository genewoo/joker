package guandan

import (
	"testing"

	"github.com/genewoo/joker/internal/deck"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type GuandanTestSuite struct {
	suite.Suite
	lastRanking [4]int
	teamLevels  [2]string
	hands       [4]*deck.Hand
}

func (suite *GuandanTestSuite) SetupTest() {
	suite.lastRanking = [4]int{1, 2, 3, 4}
	suite.teamLevels = [2]string{"2", "3"}

	// Initialize empty hands
	suite.hands = [4]*deck.Hand{
		deck.NewHand(),
		deck.NewHand(),
		deck.NewHand(),
		deck.NewHand(),
	}
}

// Helper to create a hand with specific cards
func (suite *GuandanTestSuite) createHand(cards [][]string) *deck.Hand {
	h := deck.NewHand()
	for _, card := range cards {
		h.AddCard(&deck.Card{Value: card[0], Suit: string(card[1])})
	}
	return h
}

func (suite *GuandanTestSuite) TestSwapCardsRedJokerRule() {
	tests := []struct {
		name           string
		lastRanking    [4]int
		giverHands     [][][]string
		expectedResult bool
	}{
		{
			name:        "Single giver with 1 red joker",
			lastRanking: [4]int{1, 2, 3, 4},
			giverHands: [][][]string{
				{[]string{"2", "♠"}, []string{"3", "♠"}},
				{},
				{},
				{[]string{"Joker", "Red"}, []string{"4", "♠"}, []string{"2", "♠"}, []string{"3", "♠"}},
			},
			expectedResult: true,
		},
		{
			name:        "Single giver with 2 red jokers",
			lastRanking: [4]int{1, 2, 3, 4},
			giverHands: [][][]string{
				{[]string{"2", "♠"}, []string{"3", "♠"}},
				{},
				{},
				{[]string{"Joker", "Red"}, []string{"Joker", "Red"}, []string{"2", "♠"}, []string{"3", "♠"}},
			},
			expectedResult: false,
		},
		{
			name:        "Two givers with total 2 red jokers",
			lastRanking: [4]int{1, 3, 2, 4},
			giverHands: [][][]string{
				{[]string{"2", "♠"}, []string{"5", "♠"}},
				{[]string{"4", "♠"}},
				{[]string{"Joker", "Red"}, []string{"3", "♠"}},
				{[]string{"Joker", "Red"}, []string{"2", "♠"}, []string{"3", "♠"}},
			},
			expectedResult: false,
		},
		{
			name:        "Two givers with total 1 black and 1 red jokers",
			lastRanking: [4]int{1, 3, 2, 4},
			giverHands: [][][]string{
				{[]string{"2", "♠"}, []string{"5", "♠"}},
				{[]string{"4", "♠"}},
				{[]string{"Joker", "BW"}, []string{"3", "♠"}},
				{[]string{"Joker", "Red"}, []string{"2", "♠"}, []string{"3", "♠"}},
			},
			expectedResult: true,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			// Reset hands
			suite.SetupTest()

			// Assign hands based on last ranking
			// Last ranking is [1,2,3,4], so:
			// - Player 1 is lastRanking[0]
			// - Player 2 is lastRanking[1]
			// - Player 3 is lastRanking[2]
			// - Player 4 is lastRanking[3]

			// Givers are always the last players in ranking
			for i, giverHand := range tt.giverHands {
				// Assign to last players (3 and 4)
				playerIndex := len(suite.lastRanking) - len(tt.giverHands) + i
				suite.hands[suite.lastRanking[playerIndex]-1] = suite.createHand(giverHand)
			}

			game := newGameWithHands(tt.lastRanking, suite.teamLevels, suite.hands)

			// Test swapCards
			result := game.SwapCards()
			suite.Equal(tt.expectedResult, result)
		})
	}

	// Create test hands
	for i := 0; i < 4; i++ {
		suite.hands[i] = deck.NewHand()
		for j := 0; j < 27; j++ {
			suite.hands[i].AddCard(&deck.Card{Value: "2", Suit: "♠"})
		}
	}
}

func TestGuandanSuite(t *testing.T) {
	suite.Run(t, new(GuandanTestSuite))
}

func (suite *GuandanTestSuite) TestNewGame() {
	tests := []struct {
		name        string
		lastRanking [4]int
		teamLevels  [2]string
		shouldPanic bool
	}{
		{
			name:        "Valid input",
			lastRanking: [4]int{1, 2, 3, 4},
			teamLevels:  [2]string{"2", "3"},
			shouldPanic: false,
		},
		{
			name:        "Invalid lastRanking - out of range",
			lastRanking: [4]int{5, 2, 3, 4},
			teamLevels:  [2]string{"2", "3"},
			shouldPanic: true,
		},
		{
			name:        "Invalid lastRanking - duplicate values",
			lastRanking: [4]int{1, 2, 3, 1},
			teamLevels:  [2]string{"2", "3"},
			shouldPanic: true,
		},
		{
			name:        "Invalid teamLevels - invalid card value",
			lastRanking: [4]int{1, 2, 3, 4},
			teamLevels:  [2]string{"1", "3"},
			shouldPanic: true,
		},
		{
			name:        "Invalid teamLevels - empty value",
			lastRanking: [4]int{1, 2, 3, 4},
			teamLevels:  [2]string{"", "3"},
			shouldPanic: true,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			if tt.shouldPanic {
				assert.Panics(suite.T(), func() {
					NewGame(tt.lastRanking, tt.teamLevels)
				})
			} else {
				game := NewGame(tt.lastRanking, tt.teamLevels)
				assert.NotNil(suite.T(), game)
				assert.Equal(suite.T(), tt.lastRanking, game.lastRanking)
				assert.Equal(suite.T(), tt.teamLevels[0], game.teams[0].level)
				assert.Equal(suite.T(), tt.teamLevels[1], game.teams[1].level)
			}
		})
	}
}

func (suite *GuandanTestSuite) TestDealCards() {
	game := NewGame(suite.lastRanking, suite.teamLevels)
	game.DealCards()

	suite.Run("Check deck initialization", func() {
		assert.NotNil(suite.T(), game.deck)
		assert.Equal(suite.T(), 0, len(game.deck.Cards))
	})

	suite.Run("Check dealer assignment", func() {
		assert.Equal(suite.T(), suite.lastRanking[0], game.dealer)
	})

	suite.Run("Check card distribution", func() {
		for _, player := range game.players {
			assert.Equal(suite.T(), 27, len(player.hand.Cards))
		}
	})
}

func (suite *GuandanTestSuite) TestSwapCards() {
	hands := [4]*deck.Hand{
		deck.NewHand(&deck.Card{Value: "10", Suit: "♠"}),
		deck.NewHand(&deck.Card{Value: "K", Suit: "♠"}),
		deck.NewHand(&deck.Card{Value: "9", Suit: "♠"}),
		deck.NewHand(&deck.Card{Value: "J", Suit: "♠"}),
	}
	suite.Run("Standard swap", func() {

		game := newGameWithHands([4]int{1, 2, 3, 4}, suite.teamLevels, hands)
		game.SwapCards()

		assert.Equal(suite.T(), "J", game.players[0].hand.Cards[0].Value)
		assert.Equal(suite.T(), "10", game.players[3].hand.Cards[0].Value)
	})

	suite.Run("Team swap", func() {
		game := newGameWithHands([4]int{1, 3, 2, 4}, suite.teamLevels, hands)
		game.SwapCards()

		// assert.Equal(suite.T(), "K", game.players[0].hand.Cards[0].Value)
		// assert.Equal(suite.T(), "10", game.players[1].hand.Cards[0].Value)
		// assert.Equal(suite.T(), "J", game.players[2].hand.Cards[0].Value)
		// assert.Equal(suite.T(), "9", game.players[3].hand.Cards[0].Value)
	})
}

func (suite *GuandanTestSuite) TestUpdateLevel() {
	game := newGameWithHands(suite.lastRanking, suite.teamLevels, suite.hands)

	suite.Run("Update level for winning team", func() {
		winningTeam := game.teams[0]
		expectedLevel := winningTeam.level
		game.UpdateLevel(winningTeam)
		assert.Equal(suite.T(), expectedLevel, game.currentLevel)
	})
}

func (suite *GuandanTestSuite) TestNextLevel() {
	game := newGameWithHands(suite.lastRanking, suite.teamLevels, suite.hands)

	suite.Run("Normal level progression", func() {
		winningTeam := game.teams[0]
		expectedLevel := "3"
		game.NextLevel(winningTeam)
		assert.Equal(suite.T(), expectedLevel, winningTeam.level)
		assert.Equal(suite.T(), expectedLevel, game.currentLevel)
	})

	suite.Run("Max level progression", func() {
		winningTeam := game.teams[0]
		winningTeam.level = "A"
		game.NextLevel(winningTeam)
		assert.Equal(suite.T(), "A", winningTeam.level)
	})
}
