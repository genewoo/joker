package holdem

import (
	"testing"

	"github.com/genewoo/joker/internal/dealer"
	"github.com/stretchr/testify/assert"
)

func TestNewGame(t *testing.T) {
	tests := []struct {
		name          string
		gameType      GameType
		numPlayers    int
		expectedCards int // number of cards in deck
	}{
		{
			name:          "Texas Hold'em",
			gameType:      Texas,
			numPlayers:    4,
			expectedCards: 52, // full deck
		},
		{
			name:          "Short Deck Hold'em",
			gameType:      Short,
			numPlayers:    4,
			expectedCards: 36, // 52 - 16 (cards 2-5)
		},
		{
			name:          "Omaha Hold'em",
			gameType:      Omaha,
			numPlayers:    4,
			expectedCards: 52, // full deck
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			game := NewGame(tt.gameType, tt.numPlayers)
			assert.NotNil(t, game)
			assert.Equal(t, tt.numPlayers, len(game.Players))
			assert.IsType(t, &dealer.StandardDealer{}, game.dealer)
			assert.NotNil(t, game.deck)
			assert.Equal(t, tt.expectedCards, len(game.deck.Cards))
			assert.Equal(t, tt.gameType, game.gameType)
		})
	}
}

func TestStartHand(t *testing.T) {
	tests := []struct {
		name              string
		gameType          GameType
		numPlayers        int
		expectedHoleCards int // number of hole cards per player
	}{
		{
			name:              "Texas Hold'em - 2 hole cards",
			gameType:          Texas,
			numPlayers:        2,
			expectedHoleCards: 2,
		},
		{
			name:              "Short Deck Hold'em - 2 hole cards",
			gameType:          Short,
			numPlayers:        3,
			expectedHoleCards: 2,
		},
		{
			name:              "Omaha Hold'em - 4 hole cards",
			gameType:          Omaha,
			numPlayers:        4,
			expectedHoleCards: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			game := NewGame(tt.gameType, tt.numPlayers)
			err := game.StartHand()
			assert.NoError(t, err)

			for _, player := range game.Players {
				assert.Equal(t, tt.expectedHoleCards, len(player.Cards))
			}
		})
	}
}

func TestDealFlop(t *testing.T) {
	game := NewGame(Texas, 2)
	_ = game.StartHand()

	err := game.DealFlop()
	assert.NoError(t, err)
	assert.Equal(t, 3, len(game.Community))
}

func TestDealTurnOrRiver(t *testing.T) {
	game := NewGame(Texas, 2)
	_ = game.StartHand()
	_ = game.DealFlop()

	err := game.DealTurnOrRiver()
	assert.NoError(t, err)
	assert.Equal(t, 4, len(game.Community))
	assert.Equal(t, 2, len(game.burnCards)) // 1 for flop, 1 for turn
}

func TestBurnCard(t *testing.T) {
	game := NewGame(Texas, 2)
	initialCount := game.deck.Count()

	err := game.burnCard()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(game.burnCards))
	assert.Equal(t, initialCount-1, game.deck.Count())
}

func TestBurnCardError(t *testing.T) {
	game := NewGame(Texas, 2)
	// Empty the deck
	for i := 0; i < 52; i++ {
		_ = game.burnCard()
	}

	err := game.burnCard()
	assert.Error(t, err)
	assert.Equal(t, "no cards left to burn", err.Error())
}
